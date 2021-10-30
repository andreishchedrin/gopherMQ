package server

import (
	"fmt"
	"github.com/gofiber/websocket/v2"
	"os"
	"strconv"
)

var channels = make(map[string]map[*websocket.Conn]Client)
var clients = make(map[*websocket.Conn]Client)
var register = make(chan *websocket.Conn)
var ws = make(chan *websocket.Conn)
var unregister = make(chan *websocket.Conn)
var messageErrors = make(chan error)
var broadcastMessage = make(chan *Pusher)

func (s *FiberServer) WebsocketListen() {
	debug, _ := strconv.Atoi(os.Getenv("ENABLE_WS_LOG"))

	for {
		select {
		case connection := <-register:
			clients[connection] = Client{}
			if debug == 1 {
				s.Logger.Log(fmt.Sprintf("Connection registered %v", connection))
				s.Logger.Log(fmt.Sprintf("Clients pool now is:  %v", clients))
			}
		case connection := <-unregister:
			delete(clients, connection)
			if debug == 1 {
				s.Logger.Log("Connection unregistered")
			}
		case subscribe := <-ws:
			_, ok := clients[subscribe]
			if ok {
				messageType, messageBody, err := subscribe.ReadMessage()
				if err != nil {
					messageErrors <- err
				}

				if messageType == websocket.TextMessage {
					channels[string(messageBody)] = make(map[*websocket.Conn]Client)
					channels[string(messageBody)][subscribe] = Client{}
				}

				if debug == 1 {
					s.Logger.Log(fmt.Sprintf("Websocket message received of type: %v", messageType))
					s.Logger.Log(fmt.Sprintf("Message received: %s", string(messageBody)))
				}
			}

		case message := <-broadcastMessage:
			for connection := range channels[message.Name] {
				if err := connection.WriteMessage(websocket.TextMessage, []byte(message.Message)); err != nil {
					s.Logger.Log(fmt.Sprintf("write error: %v", err))
					connection.WriteMessage(websocket.CloseMessage, []byte{})
					connection.Close()
				}
			}
		}
	}
}
