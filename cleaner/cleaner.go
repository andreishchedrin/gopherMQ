package cleaner

import (
	"andreishchedrin/gopherMQ/repository"
	"time"
)

type AbstractCleaner interface {
	StartCleaner()
	StopCleaner()
}

type Cleaner struct {
	Repo   repository.AbstractRepository
	Period string
}

var CleanerExit = make(chan bool)

func (c *Cleaner) StartCleaner() {
	go func() {
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
