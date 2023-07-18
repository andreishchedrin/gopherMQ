package main

import (
	"andreishchedrin/gopherMQ/app"
	"andreishchedrin/gopherMQ/config"
)

func main() {
	cfg, err := config.NewConfig("config/.env")
	if err != nil {
		panic(err)
	}

	app := app.NewApp(cfg)
	app.Start()
}
