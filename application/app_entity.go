package application

import "github.com/go-fuego/fuego"

type ServerConfig struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

type IController interface {
	Routes(server *fuego.Server)
}

type IService interface {
	PostConstruct()
}
