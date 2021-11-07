package storage

import (
	"andreishchedrin/gopherMQ/logger"
	"github.com/golang-collections/collections/queue"
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
	Data   map[string]*queue.Queue
	Logger logger.AbstractLogger
	Debug  int
}
