package db

import (
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strconv"
)

type AbstractDb interface {
	Prepare()
	Execute(query string)
	ExecuteWithParams(query string, params ...interface{})
	Close()
}

var Db AbstractDb

func init() {
	debug, _ := strconv.Atoi(os.Getenv("ENABLE_DB_LOG"))
	Db = &Sqlite{Connect(), debug}
}

func Prepare() {
	Db.Prepare()
}

func Execute(query string) {
	Db.Execute(query)
}

func ExecuteWithParams(query string, params ...interface{}) {
	Db.ExecuteWithParams(query, params...)
}

func Close() {
	Db.Close()
}
