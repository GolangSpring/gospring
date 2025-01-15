package application

import (
	"fmt"
	"github.com/go-fuego/fuego"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"time"
)

type Application struct {
	ContextCollection []*ApplicationContext
	AppConfig         *Config
	Server            *fuego.Server
}

func (app *Application) setupLogger() {
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}

	consoleWriter := zerolog.ConsoleWriter{
		Out:          os.Stdout,
		TimeFormat:   time.DateTime,
		TimeLocation: time.UTC,
	}
	logConfig := app.AppConfig.LogConfig

	var consoleLogWriter io.Writer
	switch logConfig.LogMode {
	case LogModeJson:
		log.Info().Msg("Setting up JSON logger")
		consoleLogWriter = &lumberjack.Logger{
			Filename:   logConfig.FileName,
			MaxSize:    logConfig.MaxSize,
			MaxBackups: logConfig.MaxBackups,
			MaxAge:     logConfig.MaxAge,
			Compress:   logConfig.Compress,
			LocalTime:  false,
		}
	case LogModeText:
		log.Info().Msg("Setting up text logger")
		consoleLogWriter = zerolog.ConsoleWriter{
			Out: &lumberjack.Logger{
				Filename:   logConfig.FileName,
				MaxSize:    logConfig.MaxSize,
				MaxBackups: logConfig.MaxBackups,
				MaxAge:     logConfig.MaxAge,
				Compress:   logConfig.Compress,
			},
			NoColor:      true,
			TimeFormat:   time.DateTime,
			TimeLocation: time.UTC,
		}
	}
	multiWriter := zerolog.MultiLevelWriter(consoleWriter, consoleLogWriter)
	// Set the global logger to use the console writer
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		shortFileName := filepath.Base(file)
		return fmt.Sprintf("%s:%d", shortFileName, line)
	}

	log.Logger = zerolog.New(multiWriter).With().Timestamp().Caller().Logger()
}

func (app *Application) InjectContextCollection(appContextCollection ...*ApplicationContext) {
	app.ContextCollection = appContextCollection
}

func (app *Application) InjectContext(appContext *ApplicationContext) {
	app.ContextCollection = append(app.ContextCollection, appContext)
}

func (app *Application) GetContext(contextName string) *ApplicationContext {
	for _, ctx := range app.ContextCollection {
		if ctx.Name == contextName {
			return ctx
		}
	}
	return nil
}

func MustNewApplication(config *Config) *Application {
	withAddr := fuego.WithAddr(fmt.Sprintf("%s:%d", config.ServerConfig.Address, config.ServerConfig.Port))
	_server := fuego.NewServer(withAddr)
	app := &Application{
		AppConfig: config,
		Server:    _server,
	}
	app.setupLogger()
	return app
}

func (app *Application) registerControllerRoutes() {
	for _, _context := range app.ContextCollection {
		for _, _controller := range _context.Controllers {
			log.Info().Msgf("Registering routes for web: %s", reflect.TypeOf(_controller).String())
			_controller.Routes(app.Server)
		}
	}

}

func (app *Application) registerControllerMiddlewares() {
	for _, _context := range app.ContextCollection {
		for _, _controller := range _context.Controllers {
			log.Info().Msgf("Registering middleware for web: %s", reflect.TypeOf(_controller).String())
			for _, middleware := range _controller.Middlewares() {
				fuego.Use(app.Server, middleware)
			}

		}
	}

}

func (app *Application) postConstructServices() {
	for _, _context := range app.ContextCollection {
		for _, _service := range _context.Services {
			log.Info().Msgf("PostConstruct for service: %s", reflect.TypeOf(_service).String())
			_service.PostConstruct()
		}
	}
	log.Info().Msg("PostConstruct for all services completed")
}

func (app *Application) Run() {
	log.Info().Msg("Starting application...")
	fmt.Printf("%s\n", app.AppConfig.AsJson())
	app.postConstructServices()
	app.registerControllerMiddlewares()
	app.registerControllerRoutes()

	err := app.Server.Run()
	if err != nil {
		log.Fatal().Msgf("Failed to start server: %v", err)
	}

}
