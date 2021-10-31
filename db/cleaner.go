package db

import (
	"sync"
	"time"
)

func (sqlite *Sqlite) StartCleaner(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-sqlite.CleanerExit:
				return
			default:
				sqlite.deleteOverdueMessages()
				time.Sleep(30 * time.Minute)
			}
		}
	}()
}

func (sqlite *Sqlite) StopCleaner() {
	sqlite.CleanerExit <- true
}
