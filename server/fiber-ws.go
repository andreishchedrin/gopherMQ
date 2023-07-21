package server

import (
	"fmt"
	"github.com/gofiber/websocket/v2"
)

func (s *FiberServer) WebsocketListen() {
	for {
		select {
		case <-s.WsExit:
			return
		case connection := <-s.Ws.Register:
			s.Ws.Clients[connection] = Client{}
			if s.Ws.EnableWsLog == 1 {
				s.Logger.Log(fmt.Sprintf("Connection registered %v", connection))
				s.Logger.Log(fmt.Sprintf("Clients pool now is:  %v", s.Ws.Clients))
			}
		case connection := <-s.Ws.Unregister:
			delete(s.Ws.Clients, connection)
			if s.Ws.EnableWsLog == 1 {
				s.Logger.Log("Connection unregistered")
			}
		case subscribe := <-s.Ws.Ws:
			_, ok := s.Ws.Clients[subscribe]
			if ok {
				messageType, messageBody, err := subscribe.ReadMessage()
				if err != nil {
					s.Ws.MessageErrors <- err
				}

				if messageType == websocket.TextMessage {
					s.Ws.Channels[string(messageBody)] = make(map[*websocket.Conn]Client)
					s.Ws.Channels[string(messageBody)][subscribe] = Client{}
				}

				if s.Ws.EnableWsLog == 1 {
					s.Logger.Log(fmt.Sprintf("Websocket message received of type: %v", messageType))
					s.Logger.Log(fmt.Sprintf("Message received: %s", string(messageBody)))
				}
			}

		case message := <-s.Ws.BroadcastMessage:
			for connection := range s.Ws.Channels[message.Channel] {
				if err := connection.WriteMessage(websocket.TextMessage, []byte(message.Message)); err != nil {
					s.Logger.Log(fmt.Sprintf("write error: %v", err))
					connection.WriteMessage(websocket.CloseMessage, []byte{})
					connection.Close()
				}
			}
		}
	}
}
