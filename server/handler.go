package server

import (
	"andreishchedrin/gopherMQ/storage"
	"fmt"
	"time"
)

type AbstractHandler interface {
	SetterHandler(setter *Setter)
	PusherHandler(pusher *Pusher)
	PullerHandler(puller *Puller)
}

type Handler struct{}

var handler AbstractHandler

func init() {
	handler = &Handler{}
}

func (h *Handler) SetterHandler(setter *Setter) {
	fmt.Println(setter)
	fmt.Println("Set channel type in storage!")
}

func (h *Handler) PusherHandler(pusher *Pusher) {
	storage.IncomeData <- storage.Message{Key: storage.Key{Name: pusher.Name}, Value: storage.Value{Text: pusher.Message, CreatedAt: time.Now()}}
}

func (h *Handler) PullerHandler(puller *Puller) {
	fmt.Println(puller)
	fmt.Println("Pull!")
}
