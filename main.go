package main

import (
	_ "andreishchedrin/gopherMQ/config"
	"andreishchedrin/gopherMQ/container"
)

func main() {
	app := container.NewApp()
	app.Start()
}
