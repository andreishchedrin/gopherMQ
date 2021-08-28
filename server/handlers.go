package server

import (
	"andreishchedrin/gopherMQ/logger"
	"andreishchedrin/gopherMQ/storage"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
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

	//send to broadcast channel

	return c.SendStatus(200)
}

func PusherHandler(c *fiber.Ctx) error {
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

func PullerHandler(c *fiber.Ctx) error {
	c.Accepts("application/json")
	puller := new(Puller)

	if err := c.BodyParser(puller); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	errors := ValidateStruct(*puller)
	if errors != nil {
		return c.JSON(errors)
	}

	q, err := storage.Storage.Get(storage.Key{Name: puller.Name})
	if err != nil {
		return c.JSON(err.Error())
	}

	return c.Status(200).JSON(q.Dequeue())
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

var channels = make(map[string]map[*websocket.Conn]Client)
var clients = make(map[*websocket.Conn]Client)
var register = make(chan *websocket.Conn)
var broadcast = make(chan string)
var unregister = make(chan *websocket.Conn)

func Listen() {
	for {
		select {
		case connection := <-register:
			clients[connection] = Client{}
			logger.Write(fmt.Sprintf("Connection registered %v", connection))
			logger.Write(fmt.Sprintf("Clients pool now is:  %v", clients))
		case message := <-broadcast:
			logger.Write(fmt.Sprintf("Message received: %s", message))

			//channels[message][] =
			//puller := new(Puller)
			//
			//if err := c.BodyParser(puller); err != nil {
			//	logger.Write(fmt.Sprintf("Context BodyParser message error: %s", err.Error()))
			//}
			//
			//errors := ValidateStruct(*puller)
			//if errors != nil {
			//	return c.JSON(errors)
			//}

			// Send the message to all clients
			//for connection := range clients {
			//	if err := connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			//		logger.Write(fmt.Sprintf("Write error: %v", err))
			//
			//		connection.WriteMessage(websocket.CloseMessage, []byte{})
			//		connection.Close()
			//		delete(clients, connection)
			//	}
			//}

		case connection := <-unregister:
			// Remove the client from the hub
			delete(clients, connection)
			logger.Write("Connection unregistered")
		}
	}
}
