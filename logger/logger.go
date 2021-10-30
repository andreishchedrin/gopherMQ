package logger

import (
	"log"
	"os"
)

type Logger struct {
	File string
}

type AbstractLogger interface {
	Log(text interface{})
}

func (l *Logger) Log(text interface{}) {
	f, err := os.OpenFile(l.File, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println(text)
}
