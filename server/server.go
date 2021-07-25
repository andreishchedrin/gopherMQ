package server

import (
	"andreishchedrin/gopherMQ/logger"
	"github.com/gofiber/fiber/v2"
	"os"
	"sync"
)

type FiberServer struct {
	app  *fiber.App
	port string
}

type AbstractServer interface {
	Serve() error
	Shutdown() error
}

func (s *FiberServer) Serve() error {
	s.app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	return s.app.Listen(":" + s.port)
}

func (s *FiberServer) Shutdown() error {
	return s.app.Shutdown()
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
