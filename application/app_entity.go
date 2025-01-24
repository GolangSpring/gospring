package application

import (
	"github.com/go-fuego/fuego"
	"net/http"
)

type ServerMode string

const (
	Development ServerMode = "dev"
	Production  ServerMode = "prod"
)

type ServerConfig struct {
	Address string     `yaml:"address" validate:"required"`
	Port    int        `yaml:"port" validate:"required"`
	Mode    ServerMode `yaml:"mode" validate:"required,oneof=dev prod"`
}

type IController interface {
	Routes(server *fuego.Server)
	Middlewares() []func(next http.Handler) http.Handler
}

type IService interface {
	PostConstruct()
}
