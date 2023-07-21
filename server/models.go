package server

import (
	"andreishchedrin/gopherMQ/logger"
	"andreishchedrin/gopherMQ/repository"
	"andreishchedrin/gopherMQ/service"
	"andreishchedrin/gopherMQ/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
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
	Ws             *FiberServerWs
	WsExit         chan bool
}

type FiberServerWs struct {
	Channels         map[string]map[*websocket.Conn]Client
	Clients          map[*websocket.Conn]Client
	Register         chan *websocket.Conn
	Ws               chan *websocket.Conn
	Unregister       chan *websocket.Conn
	MessageErrors    chan error
	BroadcastMessage chan *Push
	EnableWsLog      int
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
