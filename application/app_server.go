package application

import (
	"fmt"
	"github.com/go-fuego/fuego"
)

type ApplicationServer struct {
	ServerConfig *ServerConfig
	Engine       *fuego.Server
}

func NewApplicationServer(serverConfig *ServerConfig) *ApplicationServer {
	serverAddr := fmt.Sprintf("%s:%d", serverConfig.Address, serverConfig.Port)
	return &ApplicationServer{
		ServerConfig: serverConfig,
		Engine:       fuego.NewServer(fuego.WithAddr(serverAddr)),
	}
}

func (server *ApplicationServer) Run() error {
	return server.Engine.Run()
}
