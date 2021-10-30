package container

import (
	"andreishchedrin/gopherMQ/db"
	"andreishchedrin/gopherMQ/logger"
	"andreishchedrin/gopherMQ/server"
	"andreishchedrin/gopherMQ/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-collections/collections/queue"
	"os"
	"strconv"
)

var LoggerInstance logger.AbstractLogger
var DbInstance db.AbstractDb
var StorageInstance storage.AbstractStorage
var ServerInstance server.AbstractServer

func init() {
	LoggerInstance = &logger.Logger{File: os.Getenv("LOG_FILE")}

	enableDbLog, _ := strconv.Atoi(os.Getenv("ENABLE_DB_LOG"))
	DbInstance = &db.Sqlite{
		db.Connect(os.Getenv("DB_DRIVER_NAME"), os.Getenv("DB_DATA_SOURCE_NAME")),
		enableDbLog,
		LoggerInstance,
	}

	StorageInstance = &storage.QueueStorage{make(map[string]*queue.Queue), LoggerInstance}
	ServerInstance = &server.FiberServer{fiber.New(), os.Getenv("SERVER_PORT"), LoggerInstance, DbInstance, StorageInstance}
}
