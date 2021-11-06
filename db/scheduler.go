package db

import (
	"os"
	"strconv"
	"sync"
	"time"
)

func (sqlite *Sqlite) StartScheduler(wg *sync.WaitGroup) {
	wg.Add(1)

	timeout, _ := strconv.Atoi(os.Getenv("SCHEDULER_TIMEOUT"))
	go func(timeout int) {
		defer wg.Done()
		for {
			select {
			case <-sqlite.SchedulerExit:
				return
			default:
				//
				time.Sleep(time.Duration(timeout) * time.Second)
			}
		}
	}(timeout)
}

func (sqlite *Sqlite) StopScheduler() {
	sqlite.SchedulerExit <- true
}
