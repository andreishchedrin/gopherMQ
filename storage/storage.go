package storage

import (
	"andreishchedrin/gopherMQ/logger"
	"fmt"
	"github.com/golang-collections/collections/queue"
	"sync"
	"time"
)

type AbstractStorage interface {
	Set(key Key) *queue.Queue
	Get(key Key) (*queue.Queue, error)
	Delete(key Key) (bool, error)
	Flush()
}

var Storage AbstractStorage
var PushData = make(chan Message)

func init() {
	Storage = &QueueStorage{make(map[string]*queue.Queue)}
}

func Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case item := <-PushData:
				q := Storage.Set(item.Key)
				q.Enqueue(item.Value)
				logger.Write(fmt.Sprintf("Put to queue: %s - %s.", item.Key.Name, item.Value.Text))
			case <-time.After(time.Second * 5):
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

func Get(key Key) (*queue.Queue, error) {
	return Storage.Get(key)
}

func Set(key Key) *queue.Queue {
	return Storage.Set(key)
}

func Delete(key Key) (bool, error) {
	return Storage.Delete(key)
}

func Flush() {
	Storage.Flush()
}
