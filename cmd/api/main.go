package main

import "go-spring/application"

func main() {

	appConfig := application.MustNewAppConfig("./config.yaml")
	app := application.MustNewApplication(appConfig)
	app.InjectContextCollection()
	app.Run()
}
