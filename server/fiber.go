package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"sync"
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

func (s *FiberServer) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := s.Serve()
		if err != nil {
			s.Logger.Log(err)
		}
	}()

	s.StartScheduler(wg)
}

func (s *FiberServer) Stop() error {

	return s.Shutdown()
}
