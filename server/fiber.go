package server

import "github.com/gofiber/fiber/v2"

type FiberServer struct {
	app  *fiber.App
	port string
}

func (s *FiberServer) Serve(h AbstractHandler) error {
	s.app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	s.app.Static("/info", "./static/info.html")

	s.app.Post("/set", func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		setter := new(Setter)
		if err := c.BodyParser(setter); err != nil {
			return err
		}
		h.SetterHandler(setter)
		return c.SendStatus(200)
	})

	s.app.Post("/push", func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		pusher := new(Pusher)
		if err := c.BodyParser(pusher); err != nil {
			return err
		}
		h.PusherHandler(pusher)
		return c.SendStatus(200)
	})

	s.app.Post("/pull", func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		puller := new(Puller)
		if err := c.BodyParser(puller); err != nil {
			return err
		}

		return c.SendStatus(200)
	})

	return s.app.Listen(":" + s.port)
}

func (s *FiberServer) Shutdown() error {
	return s.app.Shutdown()
}
