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

	s.app.Post("/push", PusherHandler)

	s.app.Post("/pull", PullerHandler)

	s.app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		defer func() {
			unregister <- c
			c.Close()
		}()

		// Register the client
		register <- c

		for {
			messageType, message, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					logger.Write(fmt.Sprintf("Read error: %v", err))
				}

				return // Calls the deferred function, i.e. closes the connection on error
			}

			if messageType == websocket.TextMessage {
				// Broadcast the received message
				broadcast <- string(message)
			} else {
				logger.Write(fmt.Sprintf("Websocket message received of type: %v", messageType))
			}
		}
	}))

	return s.app.Listen(":" + s.port)
}

func (s *FiberServer) Shutdown() error {
	return s.app.Shutdown()
}
