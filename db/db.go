package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"sync"
)

type AbstractDb interface {
	Prepare()
	Execute(query string) int64
	ExecuteWithParams(query string, params ...interface{}) int64
	QueryRow(query string, params ...interface{}) *sql.Row
	QueryRows(query string, params ...interface{}) (*sql.Rows, error)
	Close()
	InsertMessage(params ...interface{}) int64
	SelectMessage(params ...interface{}) (int64, string)
	InsertClient(params ...interface{}) int64
	InsertClientMessage(params ...interface{})
	StartCleaner(wg *sync.WaitGroup)
	StopCleaner()
}

func Connect(driverName string, dataSourceName string) *sql.DB {
	db, err := sql.Open(driverName, dataSourceName)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
