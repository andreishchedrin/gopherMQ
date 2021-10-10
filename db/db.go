package db

import (
	"andreishchedrin/gopherMQ/logger"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strconv"
)

func Prepare() {
	Db, err := sql.Open("sqlite3", "./db/persistent.db")

	if err != nil {
		log.Fatal(err)
	}

	defer Db.Close()

	debug, _ := strconv.Atoi(os.Getenv("ENABLE_DB_LOG"))

	stmtMessages := "CREATE TABLE IF NOT EXISTS messages (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, channel TEXT NOT NULL, payload TEXT NOT NULL, created_at TEXT NOT NULL);"

	Execute(Db, stmtMessages, debug)
}

func Execute(db *sql.DB, stmt string, debug int) {
	var err error
	_, err = db.Exec(stmt)
	if err != nil && debug == 1 {
		logger.Write(fmt.Sprintf("SQL statement %s error: %q", stmt, err))
	}
}
