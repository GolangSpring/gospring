package mongo

import (
	"fmt"
	"github.com/GolangSpring/gospring/application"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDataSourceConfig struct {
	Mongo struct {
		Host         string `yaml:"host" validate:"required"`
		Port         int    `yaml:"port" validate:"required"`
		User         string `yaml:"user" validate:"required"`
		Password     string `yaml:"password" validate:"required"`
		DatabaseName string `yaml:"db_name" validate:"required"`
	} `yaml:"mongodb" validate:"required"`
}

func (config *MongoDataSourceConfig) AsDSN() string {
	return fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
		config.Mongo.User,
		config.Mongo.Password,
		config.Mongo.Host,
		config.Mongo.Port,
		config.Mongo.DatabaseName,
	)
}

type MongoEngineService struct {
	Engine *mongo.Client
}

func (service *MongoEngineService) PostConstruct() {}

func MustNewMongoDataSourceConfig(configPath string) *MongoDataSourceConfig {
	return application.MustNewConfigFromFile[MongoDataSourceConfig](configPath)
}

func NewMongoEngineService(config *MongoDataSourceConfig) (*MongoEngineService, error) {
	opts := options.Client().ApplyURI(config.AsDSN())
	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, err
	}
	return &MongoEngineService{Engine: client}, nil
}
