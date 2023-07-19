package storage

import (
	"andreishchedrin/gopherMQ/logger"
	"github.com/golang-collections/collections/queue"
	"sync"
	"time"
)

type Message struct {
	Key   Key
	Value Value
}

type Key struct {
	Name string
}

type Value struct {
	Text      string
	CreatedAt time.Time
}

type QueueStorage struct {
	mu          sync.RWMutex
	Data        map[string]*queue.Queue
	Logger      logger.AbstractLogger
	Debug       int
	StorageExit chan bool
}
