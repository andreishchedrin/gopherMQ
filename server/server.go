package server

import (
	"andreishchedrin/gopherMQ/logger"
	"github.com/gofiber/fiber/v2"
	"os"
	"sync"
)

type AbstractServer interface {
	Serve() error
	Shutdown() error
}

var srv AbstractServer

func init() {
	srv = &FiberServer{fiber.New(), os.Getenv("SERVER_PORT")}
}

func Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := srv.Serve()
		if err != nil {
			logger.Write(err)
		}
	}()
}

func Stop() error {
	return srv.Shutdown()
}
