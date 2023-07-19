package app

import (
	"andreishchedrin/gopherMQ/cleaner"
	"andreishchedrin/gopherMQ/config"
	"andreishchedrin/gopherMQ/db"
	"andreishchedrin/gopherMQ/logger"
	"andreishchedrin/gopherMQ/repository"
	"andreishchedrin/gopherMQ/scheduler"
	"andreishchedrin/gopherMQ/server"
	"andreishchedrin/gopherMQ/service"
	"andreishchedrin/gopherMQ/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-collections/collections/queue"
)

type App struct {
	Logger         logger.AbstractLogger
	Db             db.AbstractDb
	Repo           repository.AbstractRepository
	Cleaner        cleaner.AbstractCleaner
	Scheduler      scheduler.AbstractScheduler
	Storage        storage.AbstractStorage
	MessageService service.AbstractMessageService
	//TODO separate this interface
	HttpServer *server.FiberServer
	GrpcServer *server.Grpc
}

func NewApp(config *config.Config) *App {
	loggerInstance := &logger.Logger{File: config.LogFile}

	dbInstance := &db.Sqlite{
		ConnectInstance: db.Connect(config.DbDriverName, config.DbDataSourceName),
		Logger:          loggerInstance,
		Debug:           config.EnableDbLog,
	}

	repoInstance := &repository.SqliteRepository{SqliteDb: dbInstance, Logger: loggerInstance}
	cleanerInstance := &cleaner.Cleaner{Repo: repoInstance, Period: config.PersistentTtlDays}

	storageInstance := &storage.QueueStorage{
		Data:        make(map[string]*queue.Queue),
		Logger:      loggerInstance,
		Debug:       config.EnableStorageLog,
		StorageExit: make(chan bool),
	}

	messageService := &service.MessageService{Storage: storageInstance}

	grpcServer := &server.Grpc{
		Port:           config.GrpcPort,
		Logger:         loggerInstance,
		MessageService: messageService,
		Repo:           repoInstance,
	}

	httpServer := &server.FiberServer{
		App:            fiber.New(),
		Port:           config.HttpPort,
		Logger:         loggerInstance,
		Repo:           repoInstance,
		Storage:        storageInstance,
		MessageService: messageService,
		WsExit:         make(chan bool),
	}

	schedulerInstance := &scheduler.Scheduler{
		Repo:    repoInstance,
		Storage: storageInstance,
		Timeout: config.SchedulerTimeout,
		//TODO remove it
		ServerMode: "",
	}

	return &App{
		Logger:         loggerInstance,
		Db:             dbInstance,
		Repo:           repoInstance,
		Cleaner:        cleanerInstance,
		Scheduler:      schedulerInstance,
		Storage:        storageInstance,
		MessageService: messageService,
		GrpcServer:     grpcServer,
		HttpServer:     httpServer,
	}
}

func (a *App) Start() {
	a.Db.Prepare()
	a.Cleaner.StartCleaner()
	a.Storage.Start()
	a.HttpServer.Start()
	go a.HttpServer.WebsocketListen()
	a.GrpcServer.Start()
	a.Scheduler.StartScheduler()
}

func (a *App) Shutdown() {
	a.Db.Close()
	a.Cleaner.StopCleaner()
	a.Storage.StopStorage()
	a.HttpServer.Stop()
	a.GrpcServer.Stop()
	a.Scheduler.StopScheduler()
}
