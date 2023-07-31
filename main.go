package main

import (
	"andreishchedrin/gopherMQ/app"
	"andreishchedrin/gopherMQ/config"
	"context"
	"log"
	"net/http"
	_ "net/http/pprof"
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

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	app := app.NewApp(cfg)
	app.Start()
	<-ctx.Done()
	app.Shutdown()
}
