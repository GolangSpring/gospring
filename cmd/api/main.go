package main

import (
	"go-spring/application"
)

func main() {
	configPath := "./config.yaml"
	appConfig := application.MustNewAppConfig(configPath)
	app := application.MustNewApplication(appConfig)
	app.Run()
}
