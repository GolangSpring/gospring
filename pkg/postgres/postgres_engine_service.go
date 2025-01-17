package postgres

import (
	"fmt"
	"github.com/GolangSpring/gospring/application"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDataSourceConfig struct {
	Postgres struct {
		Host         string `yaml:"host"`
		Port         int    `yaml:"port"`
		User         string `yaml:"user"`
		Password     string `yaml:"password"`
		DatabaseName string `yaml:"db_name"`
	} `yaml:"postgres"`
}

func MustNewPostgresDataSourceConfig(configPath string) *PostgresDataSourceConfig {
	return application.MustNewConfigFromFile[PostgresDataSourceConfig](configPath)
}

func NewPostgresEngineService(config *PostgresDataSourceConfig) (*PostgresEngineService, error) {
	dialector := postgres.Open(config.AsDSN())
	sqlEngine, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &PostgresEngineService{Engine: sqlEngine}, nil
}

func (postgresConfig *PostgresDataSourceConfig) AsDSN() string {
	config := postgresConfig.Postgres
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=prefer",
		config.Host, config.Port, config.User, config.Password, config.DatabaseName,
	)
}

type PostgresEngineService struct {
	Engine *gorm.DB
}

func (service *PostgresEngineService) PostConstruct() {}

func (service *PostgresEngineService) MigrateModels(models ...any) error {
	return service.Engine.AutoMigrate(models...)
}
