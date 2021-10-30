package server

import (
	"andreishchedrin/gopherMQ/db"
	"andreishchedrin/gopherMQ/logger"
	"andreishchedrin/gopherMQ/storage"
	"github.com/gofiber/fiber/v2"
)

type Pusher struct {
	Name    string `json:"name" validate:"required"`
	Message string `json:"message" validate:"required"`
}

type Puller struct {
	Name string `json:"name" validate:"required"`
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
