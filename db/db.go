package db

import (
	_ "github.com/mattn/go-sqlite3"
)

type AbstractDb interface {
	Prepare()
	Execute(query string, debug int)
	Close()
}

var Db AbstractDb

func init() {
	Db = &Sqlite{Connect()}
}

func Prepare() {
	Db.Prepare()
}

func Close() {
	Db.Close()
}
