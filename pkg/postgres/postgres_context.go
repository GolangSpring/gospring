package postgres

import "go-spring/application"

var ContextName = "PostgresApplicationContext"

func MustNewPostgresApplicationContext(config *PostgresDataSourceConfig) *application.ApplicationContext {

	engineService, err := NewPostgresEngineService(config)
	if err != nil {
		panic(err)
	}

	return &application.ApplicationContext{
		Name:     ContextName,
		Services: []application.IService{engineService},
	}
}
