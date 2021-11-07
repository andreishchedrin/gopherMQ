package db

import (
	"andreishchedrin/gopherMQ/logger"
	"database/sql"
	"fmt"
	"os"
)

type Sqlite struct {
	ConnectInstance *sql.DB
	Logger          logger.AbstractLogger
	CleanerExit     chan bool
	SchedulerExit   chan bool
	Debug           int
}

func (sqlite *Sqlite) Close() {
	sqlite.ConnectInstance.Close()
}

func (sqlite *Sqlite) Prepare() {
	queryMessage := "CREATE TABLE IF NOT EXISTS message (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, channel TEXT NOT NULL, payload TEXT NOT NULL, created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP);"
	queryClient := "CREATE TABLE IF NOT EXISTS client (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, ip TEXT NOT NULL, channel TEXT NOT NULL, connected_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, CONSTRAINT unique_ip_channel_index UNIQUE (ip, channel) ON CONFLICT ROLLBACK);"
	queryClientMessage := "CREATE TABLE IF NOT EXISTS client_message (client_id INTEGER NOT NULL, message_id INTEGER NOT NULL, FOREIGN KEY (client_id) REFERENCES client(id) ON DELETE CASCADE, FOREIGN KEY (message_id) REFERENCES message(id) ON DELETE CASCADE);"
	queryScheduler := "CREATE TABLE IF NOT EXISTS task (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL, channel TEXT NOT NULL, message TEXT NOT NULL, type TEXT NOT NULL, time TEXT NOT NULL, CONSTRAINT unique_name_channel_index UNIQUE (name, channel) ON CONFLICT ROLLBACK);"

	sqlite.Execute(queryMessage)
	sqlite.Execute(queryClient)
	sqlite.Execute(queryClientMessage)
	sqlite.Execute(queryScheduler)
}

func (sqlite *Sqlite) Execute(query string) int64 {
	res, err := sqlite.ConnectInstance.Exec(query)

	if err != nil && sqlite.Debug == 1 {
		sqlite.Logger.Log(fmt.Sprintf("SQL statement %s error: %q", query, err))
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
		sqlite.Logger.Log(fmt.Sprintf("SQL statement %s error: %q", query, err))
	}

	defer stmt.Close()

	res, err := stmt.Exec(params...)

	if err != nil {
		sqlite.Logger.Log(fmt.Sprintf("SQL statement %s error: %q", query, err))
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
		sqlite.Logger.Log(fmt.Sprintf("SQL statement %s error: %q", query, err))
	}

	defer stmt.Close()

	return stmt.QueryRow(params...)
}

func (sqlite *Sqlite) QueryRows(query string, params ...interface{}) (*sql.Rows, error) {
	stmt, err := sqlite.ConnectInstance.Prepare(query)

	if err != nil {
		sqlite.Logger.Log(fmt.Sprintf("SQL statement %s error: %q", query, err))
	}

	defer stmt.Close()

	return stmt.Query(params...)
}

func (sqlite *Sqlite) InsertMessage(params ...interface{}) int64 {
	query := "INSERT INTO message (channel, payload) VALUES (?, ?)"
	return sqlite.ExecuteWithParams(query, params...)
}

func (sqlite *Sqlite) SelectMessage(params ...interface{}) (int64, string) {
	query := "SELECT m.id, m.payload FROM message m WHERE m.channel = ? AND NOT EXISTS (SELECT 1 FROM client_message cm WHERE cm.message_id = m.id AND cm.client_id = ?); ORDER BY created_at DESC LIMIT 1"
	row := sqlite.QueryRow(query, params...)

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

func (sqlite *Sqlite) InsertClient(params ...interface{}) int64 {
	querySelect := "SELECT id FROM client WHERE ip = ? AND channel = ?"
	row := sqlite.QueryRow(querySelect, params...)

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
	return sqlite.ExecuteWithParams(queryInsert, params...)
}

func (sqlite *Sqlite) InsertClientMessage(params ...interface{}) {
	query := "INSERT INTO client_message (client_id, message_id) VALUES (?, ?)"
	sqlite.ExecuteWithParams(query, params...)
}

func (sqlite *Sqlite) deleteOverdueMessages() {
	query := "DELETE FROM message WHERE created_at <= datetime('now', '-" + os.Getenv("PERSISTENT_TTL_DAYS") + " days')"
	sqlite.Execute(query)
}

func (sqlite *Sqlite) InsertTask(params ...interface{}) int64 {
	query := "INSERT INTO task (name, channel, message, type, time) VALUES (?, ?, ?, ?, ?)"
	return sqlite.ExecuteWithParams(query, params...)
}

func (sqlite *Sqlite) GetTasksForWorker() {
	//
}

func (sqlite *Sqlite) DeleteTask(name string) {
	//
}
