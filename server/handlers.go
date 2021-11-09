package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) BroadcastHandler(c *fiber.Ctx) error {
	c.Accepts("application/json")
	pusher := new(Pusher)

	if err := c.BodyParser(pusher); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	errors := ValidateStruct(*pusher)
	if errors != nil {
		return c.JSON(errors)
	}

	broadcastMessage <- pusher

	return c.SendStatus(200)
}

func (s *FiberServer) PushHandler(c *fiber.Ctx) error {
	c.Accepts("application/json")
	pusher := new(Pusher)

	if err := c.BodyParser(pusher); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	errors := ValidateStruct(*pusher)
	if errors != nil {
		return c.JSON(errors)
	}

	s.Storage.Push(pusher.Channel, pusher.Message)
	return c.SendStatus(200)
}

func (s *FiberServer) PullHandler(c *fiber.Ctx) error {
	c.Accepts("application/json")
	puller := new(Puller)

	if err := c.BodyParser(puller); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	errors := ValidateStruct(*puller)
	if errors != nil {
		return c.JSON(errors)
	}

	message, err := s.Storage.Pull(puller.Channel)
	if err != nil {
		return c.JSON(err.Error())
	}

	return c.Status(200).JSON(message)
}

func (s *FiberServer) PublishHandler(c *fiber.Ctx) error {
	c.Accepts("application/json")
	pusher := new(Pusher)

	if err := c.BodyParser(pusher); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	errors := ValidateStruct(*pusher)
	if errors != nil {
		return c.JSON(errors)
	}

	s.Repo.InsertMessage([]interface{}{pusher.Channel, pusher.Message}...)

	return c.SendStatus(200)
}

func (s *FiberServer) ConsumeHandler(c *fiber.Ctx) error {
	c.Accepts("application/json")
	puller := new(Puller)

	if err := c.BodyParser(puller); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	errors := ValidateStruct(*puller)
	if errors != nil {
		return c.JSON(errors)
	}

	clientId := s.Repo.InsertClient([]interface{}{c.IP(), puller.Channel}...)

	messageId, messagePayload := s.Repo.SelectMessage([]interface{}{puller.Channel, clientId}...)

	s.Repo.InsertClientMessage([]interface{}{clientId, messageId}...)

	return c.Status(200).JSON(messagePayload)
}

func (s *FiberServer) AddTaskHandler(c *fiber.Ctx) error {
	c.Accepts("application/json")
	addTask := new(AddTask)

	if err := c.BodyParser(addTask); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	errors := ValidateStruct(*addTask)
	if errors != nil {
		return c.JSON(errors)
	}

	s.Repo.InsertTask([]interface{}{addTask.Name, addTask.Channel, addTask.Message, addTask.Type, addTask.Time}...)

	return c.SendStatus(200)
}

func ValidateStruct(s interface{}) []*ErrorResponse {
	var errors []*ErrorResponse
	validate := validator.New()
	err := validate.Struct(s)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}

	return errors
}
