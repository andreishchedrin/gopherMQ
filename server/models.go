package server

import (
	"andreishchedrin/gopherMQ/logger"
	"andreishchedrin/gopherMQ/repository"
	"andreishchedrin/gopherMQ/service"
	"andreishchedrin/gopherMQ/storage"
	"github.com/gofiber/fiber/v2"
)

type Push struct {
	Channel string `json:"channel" validate:"required"`
	Message string `json:"message" validate:"required"`
}

type Pull struct {
	Channel string `json:"channel" validate:"required"`
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

type Client struct{}

type FiberServer struct {
	App            *fiber.App
	Port           string
	Logger         logger.AbstractLogger
	Repo           repository.AbstractRepository
	Storage        storage.AbstractStorage
	MessageService service.AbstractMessageService
	WsExit         chan bool
}

type AddTask struct {
	Name    string `json:"name" validate:"required"`
	Channel string `json:"channel" validate:"required"`
	Message string `json:"message" validate:"required"`
	Type    string `json:"type" validate:"required,oneof=broadcast queue persist"`
	Time    string `json:"time" validate:"required"`
}

type DeleteTask struct {
	Name string `json:"name" validate:"required"`
}
