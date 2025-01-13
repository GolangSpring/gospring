package application

import "github.com/go-fuego/fuego"

type ServerConfig struct {
	Address string `yaml:"address" validate:"required"`
	Port    int    `yaml:"port" validate:"required"`
}

type IController interface {
	Routes(server *fuego.Server)
}

type IService interface {
	PostConstruct()
}
