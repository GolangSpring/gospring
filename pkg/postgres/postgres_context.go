package postgres

import (
	"github.com/GolangSpring/gospring/application"
	"log"
)

var ContextName = "PostgresApplicationContext"

func MustNewPostgresApplicationContext(config *PostgresDataSourceConfig) *application.ApplicationContext {

	engineService, err := NewPostgresEngineService(config)
	if err != nil {
		log.Fatalf("Failed to create Postgres engine service: %v", err)
	}

	return &application.ApplicationContext{
		Name:     ContextName,
		Services: []application.IService{engineService},
	}
}
