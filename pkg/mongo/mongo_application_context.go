package mongo

import (
	"github.com/GolangSpring/gospring/application"
	"log"
)

var ContextName = "MongoApplicationContext"

func MustNewPostgresApplicationContext(config *MongoDataSourceConfig) *application.ApplicationContext {
	engineService, err := NewMongoEngineService(config)
	if err != nil {
		log.Fatalf("Failed to create Mongo engine service: %v", err)
	}

	return &application.ApplicationContext{
		Name:     ContextName,
		Services: []application.IService{engineService},
	}
}
