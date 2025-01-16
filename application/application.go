package application

import (
	"fmt"
	"github.com/go-fuego/fuego"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"time"
)

type ResponseWriterWrapper struct {
	http.ResponseWriter
	StatusCode int
}

func (wrapper *ResponseWriterWrapper) WriteHeader(statusCode int) {
	wrapper.StatusCode = statusCode
	wrapper.ResponseWriter.WriteHeader(statusCode)
}

type Application struct {
	ContextCollection []*ApplicationContext
	AppConfig         *Config
	Server            *fuego.Server
}

func (app *Application) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the ResponseWriter to capture the status code
		wrapper := &ResponseWriterWrapper{ResponseWriter: w, StatusCode: http.StatusOK}

		next.ServeHTTP(wrapper, r)
		// Log the request details

		statusCode := wrapper.StatusCode
		logger := log.Info() // Default log level

		switch {
		case statusCode >= 500:
			logger = log.Error() // Server errors
		case statusCode >= 400:
			logger = log.Warn() // Client errors
		}

		logger.
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Str("remote_addr", r.RemoteAddr).
			Int("status", statusCode).
			Str("user_agent", r.UserAgent()).
			Dur("duration", time.Since(start)).
			Msg("Request processed")
	})
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
	fuego.Use(app.Server, app.LoggingMiddleware)

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
