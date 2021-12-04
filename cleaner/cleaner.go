package cleaner

import (
	"andreishchedrin/gopherMQ/repository"
	"sync"
	"time"
)

type AbstractCleaner interface {
	StartCleaner(wg *sync.WaitGroup)
	StopCleaner()
}

type Cleaner struct {
	Repo   repository.AbstractRepository
	Period string
}

var CleanerExit = make(chan bool)

func (c *Cleaner) StartCleaner(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-CleanerExit:
				return
			default:
				c.Repo.DeleteOverdueMessages(c.Period)
				time.Sleep(30 * time.Minute)
			}
		}
	}()
}

func (c *Cleaner) StopCleaner() {
	CleanerExit <- true
}
