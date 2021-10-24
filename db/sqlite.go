package db

import (
	"andreishchedrin/gopherMQ/logger"
	"database/sql"
	"fmt"
	"log"
)

type Sqlite struct {
	ConnectInstance *sql.DB
	Debug           int
}

func Connect() *sql.DB {
	db, err := sql.Open("sqlite3", "persistent.db")

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func (sqlite *Sqlite) Close() {
	sqlite.ConnectInstance.Close()
}

func (sqlite *Sqlite) Prepare() {
	queryMessage := "CREATE TABLE IF NOT EXISTS message (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, channel TEXT NOT NULL, payload TEXT NOT NULL, created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP);"
	queryClient := "CREATE TABLE IF NOT EXISTS client (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, ip TEXT NOT NULL, channel TEXT NOT NULL, connected_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, CONSTRAINT unique_ip_channel_index UNIQUE (ip, channel) ON CONFLICT ROLLBACK);"
	queryClientMessage := "CREATE TABLE IF NOT EXISTS client_message (client_id INTEGER NOT NULL, message_id INTEGER NOT NULL, FOREIGN KEY (client_id) REFERENCES client(id) ON DELETE CASCADE, FOREIGN KEY (message_id) REFERENCES message(id) ON DELETE CASCADE);"

	sqlite.Execute(queryMessage)
	sqlite.Execute(queryClient)
	sqlite.Execute(queryClientMessage)
}

func (sqlite *Sqlite) Execute(query string) int64 {
	res, err := sqlite.ConnectInstance.Exec(query)

	if err != nil && sqlite.Debug == 1 {
		logger.Write(fmt.Sprintf("SQL statement %s error: %q", query, err))
	}

	if res != nil {
		lastInsertId, _ := res.LastInsertId()
		return lastInsertId
	}

	return 0
}

func (sqlite *Sqlite) ExecuteWithParams(query string, params ...interface{}) int64 {
	stmt, err := sqlite.ConnectInstance.Prepare(query)

	if err != nil {
		logger.Write(fmt.Sprintf("SQL statement %s error: %q", query, err))
	}

	defer stmt.Close()

	res, err := stmt.Exec(params...)

	if err != nil {
		logger.Write(fmt.Sprintf("SQL statement %s error: %q", query, err))
	}

	if res != nil {
		lastInsertId, _ := res.LastInsertId()
		return lastInsertId
	}

	return 0
}

func (sqlite *Sqlite) QueryRow(query string, params ...interface{}) *sql.Row {
	stmt, err := sqlite.ConnectInstance.Prepare(query)

	if err != nil {
		logger.Write(fmt.Sprintf("SQL statement %s error: %q", query, err))
	}

	defer stmt.Close()

	return stmt.QueryRow(params...)
}

func (sqlite *Sqlite) QueryRows(query string, params ...interface{}) (*sql.Rows, error) {
	stmt, err := sqlite.ConnectInstance.Prepare(query)

	if err != nil {
		logger.Write(fmt.Sprintf("SQL statement %s error: %q", query, err))
	}

	defer stmt.Close()

	return stmt.Query(params...)
}
