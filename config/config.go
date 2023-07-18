package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	LogFile           string
	EnableDbLog       int
	DbDriverName      string
	DbDataSourceName  string
	PersistentTtlDays string
	EnableStorageLog  int
	GrpcPort          string
	HttpPort          string
	SchedulerTimeout  int
}

func NewConfig(path string) (*Config, error) {
	err := godotenv.Load(path)
	if err != nil {
		panic("Error loading .env file")
	}

	enableDbLog, err := strconv.Atoi(os.Getenv("ENABLE_DB_LOG"))
	if err != nil {
		panic("can't parse params")
	}

	enableStorageLog, err := strconv.Atoi(os.Getenv("ENABLE_STORAGE_LOG"))
	if err != nil {
		panic("can't parse params")
	}

	timeout, err := strconv.Atoi(os.Getenv("SCHEDULER_TIMEOUT"))
	if err != nil {
		panic("can't parse params")
	}

	return &Config{
		LogFile:           os.Getenv("LOG_FILE"),
		EnableDbLog:       enableDbLog,
		DbDriverName:      os.Getenv("DB_DRIVER_NAME"),
		DbDataSourceName:  os.Getenv("DB_DATA_SOURCE_NAME"),
		PersistentTtlDays: os.Getenv("PERSISTENT_TTL_DAYS"),
		EnableStorageLog:  enableStorageLog,
		GrpcPort:          os.Getenv("GRPC_PORT"),
		HttpPort:          os.Getenv("HTTP_PORT"),
		SchedulerTimeout:  timeout,
	}, nil
}
