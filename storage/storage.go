package storage

import (
	"github.com/golang-collections/collections/queue"
)

type AbstractStorage interface {
	Set(key Key) *queue.Queue
	Get(key Key) (*queue.Queue, error)
	Delete(key Key) (bool, error)
	Flush()
	Start()
	Push(name string, message string)
	Pull(name string) (string, error)
}
