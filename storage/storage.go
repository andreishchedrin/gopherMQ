package storage

import (
	"github.com/golang-collections/collections/queue"
	"sync"
)

type AbstractStorage interface {
	Set(key Key) *queue.Queue
	Get(key Key) (*queue.Queue, error)
	Delete(key Key) (bool, error)
	Flush()
	Start(wg *sync.WaitGroup)
	Push(name string, message string)
	Pull(name string) (string, error)
}
