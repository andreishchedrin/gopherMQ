package db

import (
	"os"
	"sync"
	"time"
)

var CleanerExit = make(chan bool)

func (sqlite *Sqlite) StartCleaner(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-CleanerExit:
				return
			default:
				sqlite.deleteOverdueMessages()
				time.Sleep(30 * time.Minute)
			}
		}
	}()
}

func (sqlite *Sqlite) StopCleaner() {
	CleanerExit <- true
}

func (sqlite *Sqlite) deleteOverdueMessages() {
	query := "DELETE FROM message WHERE created_at <= datetime('now', '-" + os.Getenv("PERSISTENT_TTL_DAYS") + " days')"
	sqlite.Execute(query)
}
