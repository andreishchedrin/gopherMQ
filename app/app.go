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
	"sync"
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
		Data:   make(map[string]*queue.Queue),
		Logger: loggerInstance,
		Debug:  config.EnableStorageLog,
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
	var wg sync.WaitGroup

	a.Db.Prepare()
	defer a.Db.Close()

	a.Cleaner.StartCleaner(&wg)
	defer a.Cleaner.StopCleaner()

	a.Storage.Start(&wg)

	a.HttpServer.Start(&wg)
	go a.HttpServer.WebsocketListen()
	defer a.HttpServer.Stop()

	a.GrpcServer.Start(&wg)
	defer a.GrpcServer.Stop()

	a.Scheduler.StartScheduler(&wg)
	defer a.Scheduler.StopScheduler()

	wg.Wait()
}
