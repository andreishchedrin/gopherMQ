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

func NewApp() *App {
	loggerInstance := &logger.Logger{File: os.Getenv("LOG_FILE")}

	enableDbLog, err := strconv.Atoi(os.Getenv("ENABLE_DB_LOG"))
	if err != nil {
		panic("can't parse params")
	}

	dbInstance := &db.Sqlite{
		ConnectInstance: db.Connect(os.Getenv("DB_DRIVER_NAME"), os.Getenv("DB_DATA_SOURCE_NAME")),
		Logger:          loggerInstance,
		Debug:           enableDbLog,
	}

	repoInstance := &repository.SqliteRepository{SqliteDb: dbInstance, Logger: loggerInstance}

	period := os.Getenv("PERSISTENT_TTL_DAYS")
	cleanerInstance := &cleaner.Cleaner{Repo: repoInstance, Period: period}

	enableStorageLog, err := strconv.Atoi(os.Getenv("ENABLE_STORAGE_LOG"))
	if err != nil {
		panic("can't parse params")
	}

	storageInstance := &storage.QueueStorage{
		Data:   make(map[string]*queue.Queue),
		Logger: loggerInstance,
		Debug:  enableStorageLog,
	}

	messageService := &service.MessageService{Storage: storageInstance}

	grpcServer := &server.Grpc{
		Port:           os.Getenv("GRPC_PORT"),
		Logger:         loggerInstance,
		MessageService: messageService,
		Repo:           repoInstance,
	}

	httpServer := &server.FiberServer{
		App:            fiber.New(),
		Port:           os.Getenv("HTTP_PORT"),
		Logger:         loggerInstance,
		Repo:           repoInstance,
		Storage:        storageInstance,
		MessageService: messageService,
	}

	timeout, err := strconv.Atoi(os.Getenv("SCHEDULER_TIMEOUT"))
	if err != nil {
		panic("can't parse params")
	}

	schedulerInstance := &scheduler.Scheduler{
		Repo:    repoInstance,
		Storage: storageInstance,
		Timeout: timeout,
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
