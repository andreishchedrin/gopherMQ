package logger

import (
	"log"
	"os"
)

type Logger struct {
	file string
}

type AbstractLogger interface {
	Log(text interface{})
}

func (l *Logger) Log(text interface{}) {
	f, err := os.OpenFile(l.file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println(text)
}

var logger AbstractLogger

func init() {
	logger = &Logger{os.Getenv("LOG_FILE")}
}

func Write(text interface{}) {
	logger.Log(text)
}
