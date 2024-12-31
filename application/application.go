package application

import (
	"fmt"
	"github.com/go-fuego/fuego"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"reflect"
	"time"
)

type Application struct {
	ContextCollection []*ApplicationContext
	AppConfig         *Config
	Server            *fuego.Server
}

func (app *Application) setupLogger() {
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339, // Optional: specify the time format
	}
	// Set the global logger to use the console writer
	log.Logger = zerolog.New(consoleWriter).With().Timestamp().Logger()
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
	return &Application{
		AppConfig: config,
		Server:    _server,
	}
}

func (app *Application) registerControllerRoutes() {
	for _, _context := range app.ContextCollection {
		for _, _controller := range _context.Controllers {
			log.Info().Msgf("Registering routes for web: %s", reflect.TypeOf(_controller).String())
			_controller.Routes(app.Server)
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
	app.setupLogger()
	log.Info().Msg("Starting application...")
	fmt.Printf("%s\n", app.AppConfig.AsJson())
	app.postConstructServices()
	app.registerControllerRoutes()
	err := app.Server.Run()
	if err != nil {
		log.Fatal().Msgf("Failed to start server: %v", err)
	}

}
