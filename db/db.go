package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type AbstractDb interface {
	Prepare()
	Execute(query string) int64
	ExecuteWithParams(query string, params ...interface{}) int64
	QueryRow(query string, params ...interface{}) *sql.Row
	QueryRows(query string, params ...interface{}) (*sql.Rows, error)
	Close()
	GetConnectInstance() *sql.DB
}

func Connect(driverName string, dataSourceName string) *sql.DB {
	db, err := sql.Open(driverName, dataSourceName)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
