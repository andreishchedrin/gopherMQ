package server

import (
	"github.com/gofiber/fiber/v2"
	"os"
)

type FiberServer struct {
	app  *fiber.App
	port string
}

type Starter interface {
	Serve()
	Shutdown()
}

func (s *FiberServer) Serve() {
	s.app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	s.app.Listen(":" + s.port)
}

func (s *FiberServer) Shutdown() {
	s.app.Shutdown()
}

var srv Starter

func init() {
	srv = &FiberServer{fiber.New(), os.Getenv("SERVER_PORT")}
}

func Start() {
	srv.Serve()
}
