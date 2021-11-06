package server

import (
	"andreishchedrin/gopherMQ/db"
	"andreishchedrin/gopherMQ/logger"
	"andreishchedrin/gopherMQ/storage"
	"github.com/gofiber/fiber/v2"
)

type Pusher struct {
	Channel string `json:"channel" validate:"required"`
	Message string `json:"message" validate:"required"`
}

type Puller struct {
	Channel string `json:"channel" validate:"required"`
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

type Client struct{}

type FiberServer struct {
	App     *fiber.App
	Port    string
	Logger  logger.AbstractLogger
	Db      db.AbstractDb
	Storage storage.AbstractStorage
}

type Task struct {
	Name    string `json:"name" validate:"required"`
	Channel string `json:"channel" validate:"required"`
	Message string `json:"message" validate:"required"`
	Type    string `json:"type" validate:"required,oneof=broadcast queue persist"`
	Time    string `json:"time" validate:"required"`
}
