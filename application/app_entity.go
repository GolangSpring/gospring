package application

import (
	"github.com/go-fuego/fuego"
	"net/http"
)

type ServerConfig struct {
	Address string `yaml:"address" validate:"required"`
	Port    int    `yaml:"port" validate:"required"`
}

type IController interface {
	Routes(server *fuego.Server)
	Middlewares() []func(next http.Handler) http.Handler
}

type IService interface {
	PostConstruct()
}
