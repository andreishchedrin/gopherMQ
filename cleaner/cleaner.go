package cleaner

import (
	"andreishchedrin/gopherMQ/repository"
	"time"
)

type AbstractCleaner interface {
	StartCleaner()
}

type Cleaner struct {
	Repo   repository.AbstractRepository
	Period string
}

func (c *Cleaner) StartCleaner() {
	go func() {
		for {
			select {
			default:
				c.Repo.DeleteOverdueMessages(c.Period)
				time.Sleep(30 * time.Minute)
			}
		}
	}()
}
