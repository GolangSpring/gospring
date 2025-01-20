package postgres

import (
	"fmt"
	"github.com/GolangSpring/gospring/application"
	"github.com/jmoiron/sqlx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type SSLMode string

const (
	SSLModeDisable    SSLMode = "disable"
	SSLModeRequire    SSLMode = "require"
	SSLModeVerifyCA   SSLMode = "verify-ca"
	SSLModeVerifyFull SSLMode = "verify-full"
)

type PostgresDataSourceConfig struct {
	Postgres struct {
		Host         string  `yaml:"host" validate:"required"`
		Port         int     `yaml:"port" validate:"required"`
		User         string  `yaml:"user" validate:"required"`
		Password     string  `yaml:"password" validate:"required"`
		DatabaseName string  `yaml:"db_name" validate:"required"`
		SSLMode      SSLMode `yaml:"ssl_mode" validate:"required,oneof=disable require verify-ca verify-full"`
	} `yaml:"postgres" validate:"required"`
}

func (postgresConfig *PostgresDataSourceConfig) AsDSN() string {
	config := postgresConfig.Postgres
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DatabaseName,
		config.SSLMode,
	)
}

var _ application.IService = (*PostgresEngineService)(nil)

type PostgresEngineService struct {
	Engine        *gorm.DB
	BuilderEngine *sqlx.DB
}

func (service *PostgresEngineService) PostConstruct() {}

func MustNewPostgresDataSourceConfig(configPath string) *PostgresDataSourceConfig {
	return application.MustNewConfigFromFile[PostgresDataSourceConfig](configPath)
}

func NewPostgresEngineService(config *PostgresDataSourceConfig) (*PostgresEngineService, error) {
	dialector := postgres.Open(config.AsDSN())
	sqlEngine, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	builderEngine, err := sqlx.Connect("postgres", config.AsDSN())
	if err != nil {
		return nil, err
	}

	return &PostgresEngineService{
		Engine:        sqlEngine,
		BuilderEngine: builderEngine,
	}, nil
}

func (service *PostgresEngineService) MigrateModels(models ...any) error {
	return service.Engine.AutoMigrate(models...)
}
