package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strconv"
)

type AbstractDb interface {
	Prepare()
	Execute(query string) int64
	ExecuteWithParams(query string, params ...interface{}) int64
	QueryRow(query string, params ...interface{}) *sql.Row
	QueryRows(query string, params ...interface{}) (*sql.Rows, error)
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

func Close() {
	Db.Close()
}

func InsertMessage(params ...interface{}) int64 {
	query := "INSERT INTO message (channel, payload) VALUES (?, ?)"
	return Db.ExecuteWithParams(query, params...)
}

func SelectMessage(params ...interface{}) (int64, string) {
	query := "SELECT m.id, m.payload FROM message m WHERE m.channel = ? AND NOT EXISTS (SELECT 1 FROM client_message cm WHERE cm.message_id = m.id AND cm.client_id = ?); ORDER BY created_at DESC LIMIT 1"
	row := Db.QueryRow(query, params...)

	return func(r *sql.Row) (int64, string) {
		var id int64
		var payload string
		err := r.Scan(&id, &payload)
		if err != nil {
			return 0, "no message found."
		}

		return id, payload
	}(row)
}

func InsertClient(params ...interface{}) int64 {
	querySelect := "SELECT id FROM client WHERE ip = ? AND channel = ?"
	row := Db.QueryRow(querySelect, params...)

	clientId := func(r *sql.Row) int64 {
		var id int64
		err := r.Scan(&id)
		if err != nil {
			return 0
		}

		return id
	}(row)

	if clientId != 0 {
		return clientId
	}

	queryInsert := "INSERT INTO client (ip, channel) VALUES (?, ?)"
	return Db.ExecuteWithParams(queryInsert, params...)
}

func InsertClientMessage(params ...interface{}) {
	query := "INSERT INTO client_message (client_id, message_id) VALUES (?, ?)"
	Db.ExecuteWithParams(query, params...)
}
