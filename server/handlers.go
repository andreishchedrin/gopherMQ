package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func (s *FiberServer) BroadcastHandler(c *fiber.Ctx) error {
	c.Accepts("application/json")
	push := new(Push)

	if err := c.BodyParser(push); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	errors := ValidateStruct(*push)
	if errors != nil {
		return c.JSON(errors)
	}

	broadcastMessage <- push

	return c.SendStatus(200)
}

func (s *FiberServer) PushHandler(c *fiber.Ctx) error {
	c.Accepts("application/json")
	push := new(Push)

	if err := c.BodyParser(push); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	errors := ValidateStruct(*push)
	if errors != nil {
		return c.JSON(errors)
	}

	s.MessageService.Push(push.Channel, push.Message)
	return c.SendStatus(200)
}

func (s *FiberServer) PullHandler(c *fiber.Ctx) error {
	c.Accepts("application/json")
	pull := new(Pull)

	if err := c.BodyParser(pull); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	errors := ValidateStruct(*pull)
	if errors != nil {
		return c.JSON(errors)
	}

	message, err := s.MessageService.Pull(pull.Channel)
	if err != nil {
		return c.JSON(err.Error())
	}

	return c.Status(200).JSON(message)
}

func (s *FiberServer) PublishHandler(c *fiber.Ctx) error {
	c.Accepts("application/json")
	push := new(Push)

	if err := c.BodyParser(push); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	errors := ValidateStruct(*push)
	if errors != nil {
		return c.JSON(errors)
	}

	s.Repo.InsertMessage([]interface{}{push.Channel, push.Message}...)

	return c.SendStatus(200)
}

func (s *FiberServer) ConsumeHandler(c *fiber.Ctx) error {
	c.Accepts("application/json")
	pull := new(Pull)

	if err := c.BodyParser(pull); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	errors := ValidateStruct(*pull)
	if errors != nil {
		return c.JSON(errors)
	}

	clientId := s.Repo.InsertClient([]interface{}{c.IP(), pull.Channel}...)

	messageId, messagePayload := s.Repo.SelectMessage([]interface{}{pull.Channel, clientId}...)

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

func (s *FiberServer) DeleteTaskHandler(c *fiber.Ctx) error {
	c.Accepts("application/json")
	deleteTask := new(DeleteTask)

	if err := c.BodyParser(deleteTask); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	errors := ValidateStruct(*deleteTask)
	if errors != nil {
		return c.JSON(errors)
	}

	s.Repo.DeleteTask([]interface{}{deleteTask.Name}...)

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
