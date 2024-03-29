package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func (s *FiberServer) Serve() error {
	s.App.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	s.App.Static("/info", "./static/info.html")

	s.App.Post("/broadcast", s.BroadcastHandler)

	s.App.Post("/push", s.PushHandler)
	s.App.Post("/pull", s.PullHandler)

	s.App.Get("/ws", websocket.New(func(c *websocket.Conn) {
		defer func() {
			s.Ws.Unregister <- c
			c.Close()
		}()

		// Register the client
		s.Ws.Register <- c

		for {
			s.Ws.Ws <- c
			select {
			case messageError := <-s.Ws.MessageErrors:
				if messageError != nil {
					if websocket.IsUnexpectedCloseError(messageError, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						s.Logger.Log(fmt.Sprintf("Read error: %v", messageError))
					}

					return // Calls the deferred function, i.e. closes the connection on error
				}
			}
		}
	}))

	s.App.Post("/publish", s.PublishHandler)
	s.App.Post("/consume", s.ConsumeHandler)

	s.App.Post("/add-task", s.AddTaskHandler)
	s.App.Post("/remove-task", s.ConsumeHandler)

	return s.App.Listen(":" + s.Port)
}

func (s *FiberServer) Shutdown() error {
	return s.App.Shutdown()
}

func (s *FiberServer) Start() {
	go func() {
		err := s.Serve()
		if err != nil {
			s.Logger.Log(err)
		}
	}()
}

func (s *FiberServer) Stop() error {
	s.WsExit <- true
	return s.Shutdown()
}
