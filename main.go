package main

import (
	"andreishchedrin/gopherMQ/app"
	"andreishchedrin/gopherMQ/config"
	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

		<-c
		cancel()
	}()

	cfg, err := config.NewConfig("config/.env")
	if err != nil {
		panic(err)
	}

	app := app.NewApp(cfg)
	app.Start()
	<-ctx.Done()
	app.Shutdown()
}
