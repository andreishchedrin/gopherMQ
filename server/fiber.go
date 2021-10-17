package server

import (
	"andreishchedrin/gopherMQ/logger"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type FiberServer struct {
	app  *fiber.App
	port string
}

func (s *FiberServer) Serve() error {
	s.app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	s.app.Static("/info", "./static/info.html")

	s.app.Post("/broadcast", BroadcastHandler)

	s.app.Post("/push", PushHandler)

	s.app.Post("/pull", PullHandler)

	s.app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		defer func() {
			unregister <- c
			c.Close()
		}()

		// Register the client
		register <- c

		for {
			ws <- c
			select {
			case messageError := <-messageErrors:
				if messageError != nil {
					if websocket.IsUnexpectedCloseError(messageError, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						logger.Write(fmt.Sprintf("Read error: %v", messageError))
					}

					return // Calls the deferred function, i.e. closes the connection on error
				}
			}
		}
	}))

	s.app.Post("/publish", PublishHandler)

	//s.app.Post('/subscribe', SubscribeHandler)

	s.app.Post("/consume", ConsumeHandler)

	return s.app.Listen(":" + s.port)
}

func (s *FiberServer) Shutdown() error {
	return s.app.Shutdown()
}
