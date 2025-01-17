package security

import (
	"github.com/GolangSpring/gospring/application"
	"github.com/GolangSpring/gospring/pkg/postgres"
	"github.com/GolangSpring/gospring/pkg/security/controller"
	securityRepository "github.com/GolangSpring/gospring/pkg/security/repository"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/rs/zerolog/log"

	securityService "github.com/GolangSpring/gospring/pkg/security/service"
)

var ContextName = "SecurityApplicationContext"

func MustNewSecurityContext(securityConfig *SecurityConfig, postgresContext *application.ApplicationContext) *application.ApplicationContext {

	service, err := application.GetServiceFromContext[*postgres.PostgresEngineService](postgresContext)
	if err != nil {
		log.Fatal().Msgf("Failed to get Postgres engine service from context: %v", err)
	}

	models := []any{securityRepository.User{}}

	if err = service.MigrateModels(models...); err != nil {
		log.Fatal().Msgf("Failed to migrate models: %v", err)
	}

	adapter, err := gormadapter.NewAdapterByDB(service.Engine)
	if err != nil {
		log.Fatal().Msgf("Failed to create Casbin adapter: %v", err)
	}

	casbinModel, err := model.NewModelFromString(securityService.ModelString)
	if err != nil {
		log.Fatal().Msgf("Failed to create Casbin model: %v", err)
	}

	enforcer, err := casbin.NewEnforcer(casbinModel, adapter)
	if err != nil {
		log.Fatal().Msgf("Failed to create Casbin enforcer: %v", err)
	}

	casbinService := securityService.NewCasbinService(enforcer)

	engine := service.Engine
	userRepo := securityRepository.NewUserRepository(engine)
	userService := securityService.NewUserService(userRepo)
	authService := securityService.NewAuthService(userService, securityConfig.Security.Secret)
	authController := controller.NewAuthController(authService, casbinService)
	casbinController := controller.NewCasbinController(casbinService, authService)
	systemController := controller.NewSystemController()

	return &application.ApplicationContext{
		Name: ContextName,
		Services: []application.IService{
			casbinService,
			userService,
			authService,
		},
		Controllers: []application.IController{
			authController,
			casbinController,
			systemController,
		},
	}
}
