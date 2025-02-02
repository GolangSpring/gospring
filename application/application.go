package application

import (
	"errors"
	"fmt"
	appMiddleware "github.com/GolangSpring/gospring/application/app_middleware"
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

func (app *Application) CheckInterfaceNilValues(interfaceType any) error {
	// Ensure the input is a struct or pointer to struct
	objectValue := reflect.ValueOf(interfaceType)
	if objectValue.Kind() == reflect.Ptr || objectValue.Kind() == reflect.Interface {
		objectValue = objectValue.Elem()
	}
	if objectValue.Kind() != reflect.Struct {
		return errors.New("input must be a struct or pointer to struct")
	}

	objectType := objectValue.Type()

	for idx := 0; idx < objectType.NumField(); idx++ {
		field := objectType.Field(idx)
		value := objectValue.Field(idx)

		// Skip unexported fields
		if !value.CanInterface() {
			continue
		}

		// Check if the field is nil for pointers, interfaces, or slices
		if (value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface || value.Kind() == reflect.Slice) && value.IsNil() {
			return fmt.Errorf("field '%s' is nil", field.Name)
		}
	}

	return nil
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
	fileLogger := &lumberjack.Logger{
		Filename:   logConfig.FileName,
		MaxSize:    logConfig.MaxSize,
		MaxBackups: logConfig.MaxBackups,
		MaxAge:     logConfig.MaxAge,
		Compress:   logConfig.Compress,
	}

	var consoleLogWriter io.Writer
	switch logConfig.LogMode {
	case LogModeJson:
		consoleLogWriter = fileLogger
	case LogModeText:
		consoleLogWriter = zerolog.ConsoleWriter{
			Out:          fileLogger,
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

	withOpenAIConfig := fuego.WithOpenAPIConfig(
		fuego.OpenAPIConfig{
			DisableLocalSave: true,
			DisableSwagger:   config.ServerConfig.Mode == Production,
		},
	)

	withAddr := fuego.WithAddr(fmt.Sprintf("%s:%d", config.ServerConfig.Address, config.ServerConfig.Port))
	_server := fuego.NewServer(withAddr, withOpenAIConfig)
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
func (app *Application) registerAppMiddlewares() {
	log.Info().Msg("Registering application middlewares")
	fuego.Use(app.Server, appMiddleware.LoggingMiddleware)
}

func (app *Application) checkContextNilValues() {
	for _, _context := range app.ContextCollection {
		services := _context.Services
		for _, service := range services {
			serviceName := reflect.TypeOf(service).String()
			if err := app.CheckInterfaceNilValues(service); err != nil {
				log.Fatal().Msgf("Service %s has nil values: %v", serviceName, err)
			}
		}

		controllers := _context.Controllers
		for _, controller := range controllers {
			controllerName := reflect.TypeOf(controller).String()
			if err := app.CheckInterfaceNilValues(controller); err != nil {
				log.Fatal().Msgf("Controller %s has nil values: %v", controllerName, err)
			}
		}
	}
}

func (app *Application) Run() {
	log.Info().Msgf("Checking nil values in context collection")
	app.checkContextNilValues()
	log.Info().Msg("Starting application...")
	app.registerAppMiddlewares()
	fmt.Printf("%s\n", app.AppConfig.AsJson())
	app.postConstructServices()
	app.registerControllerMiddlewares()
	app.registerControllerRoutes()

	err := app.Server.Run()
	if err != nil {
		log.Fatal().Msgf("Failed to start server: %v", err)
	}

}
