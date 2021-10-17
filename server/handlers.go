package server

import (
	"andreishchedrin/gopherMQ/db"
	"andreishchedrin/gopherMQ/storage"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"time"
)

func BroadcastHandler(c *fiber.Ctx) error {
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

func PushHandler(c *fiber.Ctx) error {
	c.Accepts("application/json")
	pusher := new(Pusher)

	if err := c.BodyParser(pusher); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	errors := ValidateStruct(*pusher)
	if errors != nil {
		return c.JSON(errors)
	}

	storage.PushData <- storage.Message{Key: storage.Key{Name: pusher.Name}, Value: storage.Value{Text: pusher.Message, CreatedAt: time.Now()}}

	return c.SendStatus(200)
}

func PullHandler(c *fiber.Ctx) error {
	c.Accepts("application/json")
	puller := new(Puller)

	if err := c.BodyParser(puller); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	errors := ValidateStruct(*puller)
	if errors != nil {
		return c.JSON(errors)
	}

	q, err := storage.Get(storage.Key{Name: puller.Name})
	if err != nil {
		return c.JSON(err.Error())
	}

	return c.Status(200).JSON(q.Dequeue())
}

func PublishHandler(c *fiber.Ctx) error {
	c.Accepts("application/json")
	pusher := new(Pusher)

	if err := c.BodyParser(pusher); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	errors := ValidateStruct(*pusher)
	if errors != nil {
		return c.JSON(errors)
	}

	query := "INSERT INTO message (channel, payload) VALUES (?, ?)"

	params := []interface{}{pusher.Name, pusher.Message}

	db.ExecuteWithParams(query, params...)

	return c.SendStatus(200)
}

func ConsumeHandler(c *fiber.Ctx) error {
	c.Accepts("application/json")
	puller := new(Puller)

	if err := c.BodyParser(puller); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	errors := ValidateStruct(*puller)
	if errors != nil {
		return c.JSON(errors)
	}

	//@TODO

	return c.Status(200).JSON("a")
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
