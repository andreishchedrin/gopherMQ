package service

import "andreishchedrin/gopherMQ/storage"

type AbstractMessageService interface {
	Push(channel string, message string)
	Pull(channel string) (interface{}, error)
}

type MessageService struct {
	Storage storage.AbstractStorage
}

func (ms *MessageService) Push(channel string, message string) {
	ms.Storage.Push(channel, message)
}

func (ms *MessageService) Pull(channel string) (interface{}, error) {
	return ms.Storage.Pull(channel)
}
