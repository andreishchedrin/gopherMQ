package logger

import (
	"log"
	"os"
	"sync"
)

type Logger struct {
	File string
	mu   sync.Mutex
}

type AbstractLogger interface {
	Log(text interface{})
}

func (l *Logger) Log(text interface{}) {
	l.mu.Lock()
	f, err := os.OpenFile(l.File, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println(text)
	l.mu.Unlock()
}
