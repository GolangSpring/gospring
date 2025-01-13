package application

import (
	"github.com/go-fuego/fuego"
)

type ApplicationServer struct {
	ServerConfig *ServerConfig
	Engine       *fuego.Server
}

func (server *ApplicationServer) Run() error {
	return server.Engine.Run()
}
