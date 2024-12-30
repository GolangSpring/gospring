package application

type ServerConfig struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

type Controller interface {
	RegisterRoutes()
}

type IService interface {
	PostConstruct()
}
