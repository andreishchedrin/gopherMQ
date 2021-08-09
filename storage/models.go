package storage

import "time"

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
