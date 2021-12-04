package repository

import (
	"andreishchedrin/gopherMQ/db"
	"andreishchedrin/gopherMQ/logger"
	"database/sql"
	"fmt"
)

type SqliteRepository struct {
	SqliteDb db.AbstractDb
	Logger   logger.AbstractLogger
}

type AbstractRepository interface {
	InsertMessage(params ...interface{}) int64
	SelectMessage(params ...interface{}) (int64, string)
	InsertClient(params ...interface{}) int64
	InsertClientMessage(params ...interface{})
	InsertTask(params ...interface{}) int64
	DeleteTask(params ...interface{})
	GetTasksForWorker() []Task
}

type Task struct {
	Id         int64
	Name       string
	Channel    string
	Message    string
	Type       string
	Time       string
	Repeatable int
}

func (repo *SqliteRepository) InsertMessage(params ...interface{}) int64 {
	query := "INSERT INTO message (channel, payload) VALUES (?, ?)"
	return repo.SqliteDb.ExecuteWithParams(query, params...)
}

func (repo *SqliteRepository) SelectMessage(params ...interface{}) (int64, string) {
	query := "SELECT m.id, m.payload FROM message m WHERE m.channel = ? AND NOT EXISTS (SELECT 1 FROM client_message cm WHERE cm.message_id = m.id AND cm.client_id = ?); ORDER BY created_at DESC LIMIT 1"
	row := repo.SqliteDb.QueryRow(query, params...)

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

func (repo *SqliteRepository) InsertClient(params ...interface{}) int64 {
	querySelect := "SELECT id FROM client WHERE ip = ? AND channel = ?"
	row := repo.SqliteDb.QueryRow(querySelect, params...)

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
	return repo.SqliteDb.ExecuteWithParams(queryInsert, params...)
}

func (repo *SqliteRepository) InsertClientMessage(params ...interface{}) {
	query := "INSERT INTO client_message (client_id, message_id) VALUES (?, ?)"
	repo.SqliteDb.ExecuteWithParams(query, params...)
}

func (repo *SqliteRepository) InsertTask(params ...interface{}) int64 {
	query := "INSERT INTO task (name, channel, message, type, time) VALUES (?, ?, ?, ?, ?)"
	return repo.SqliteDb.ExecuteWithParams(query, params...)
}

func (repo *SqliteRepository) DeleteTask(params ...interface{}) {
	query := "DELETE FROM task WHERE name = ?"
	repo.SqliteDb.ExecuteWithParams(query, params...)
}

func (repo *SqliteRepository) GetTasksForWorker() []Task {
	query := "SELECT * type FROM task WHERE time <= time('now')"
	sqlite := repo.SqliteDb.GetConnectInstance()

	rows, err := sqlite.Query(query)

	if err != nil {
		repo.Logger.Log(fmt.Sprintf("SQL statement %s error: %q", query, err))
		return nil
	}

	var tasks []Task
	for rows.Next() {
		task := Task{}
		err = rows.Scan(&task.Id, &task.Name, &task.Channel, &task.Message, &task.Type, &task.Time, &task.Repeatable)
		if err != nil {
			repo.Logger.Log(err)
		}
		tasks = append(tasks, task)
	}

	return tasks
}
