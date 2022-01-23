package container

import (
	"andreishchedrin/gopherMQ/cleaner"
	"andreishchedrin/gopherMQ/db"
	"andreishchedrin/gopherMQ/logger"
	"andreishchedrin/gopherMQ/repository"
	"andreishchedrin/gopherMQ/scheduler"
	"andreishchedrin/gopherMQ/server"
	"andreishchedrin/gopherMQ/service"
	"andreishchedrin/gopherMQ/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-collections/collections/queue"
	"os"
	"strconv"
)

var LoggerInstance logger.AbstractLogger
var DbInstance db.AbstractDb
var RepoInstance repository.AbstractRepository
var CleanerInstance cleaner.AbstractCleaner
var SchedulerInstance scheduler.AbstractScheduler
var StorageInstance storage.AbstractStorage
var MessageService service.AbstractMessageService
var ServerInstance server.AbstractServer

func init() {
	LoggerInstance = &logger.Logger{File: os.Getenv("LOG_FILE")}

	enableDbLog, _ := strconv.Atoi(os.Getenv("ENABLE_DB_LOG"))
	DbInstance = &db.Sqlite{
		ConnectInstance: db.Connect(os.Getenv("DB_DRIVER_NAME"), os.Getenv("DB_DATA_SOURCE_NAME")),
		Logger:          LoggerInstance,
		Debug:           enableDbLog,
	}

	RepoInstance = &repository.SqliteRepository{SqliteDb: DbInstance, Logger: LoggerInstance}

	period := os.Getenv("PERSISTENT_TTL_DAYS")
	CleanerInstance = &cleaner.Cleaner{Repo: RepoInstance, Period: period}

	enableStorageLog, _ := strconv.Atoi(os.Getenv("ENABLE_STORAGE_LOG"))
	StorageInstance = &storage.QueueStorage{
		Data:   make(map[string]*queue.Queue),
		Logger: LoggerInstance,
		Debug:  enableStorageLog,
	}

	MessageService = &service.MessageService{Storage: StorageInstance}

	serverMode := os.Getenv("MODE")

	if serverMode == "grpc" {
		ServerInstance = &server.Grpc{
			Port:           os.Getenv("SERVER_PORT"),
			Logger:         LoggerInstance,
			MessageService: MessageService,
		}
	} else {
		ServerInstance = &server.FiberServer{
			App:            fiber.New(),
			Port:           os.Getenv("SERVER_PORT"),
			Logger:         LoggerInstance,
			Repo:           RepoInstance,
			Storage:        StorageInstance,
			MessageService: MessageService,
		}
	}

	timeout, _ := strconv.Atoi(os.Getenv("SCHEDULER_TIMEOUT"))
	SchedulerInstance = &scheduler.Scheduler{
		Repo:       RepoInstance,
		Storage:    StorageInstance,
		Timeout:    timeout,
		ServerMode: serverMode,
	}
}
