package storage

import (
	"andreishchedrin/gopherMQ/logger"
	"fmt"
	"github.com/golang-collections/collections/queue"
	"os"
	"strconv"
	"sync"
	"time"
)

type QueueStorage struct {
	workerPoolSize int
	data           map[string]*queue.Queue
}

type Message struct {
	key   Key
	value Value
}

type Key struct {
	name string
}

type Value struct {
	text      string
	createdAt time.Time
}

type AbstractStorage interface {
	Set(key Key) *queue.Queue
	Get(key Key) (*queue.Queue, error)
	Delete(key Key) (bool, error)
	FlushStorage()
}

func (qs *QueueStorage) Set(key Key) *queue.Queue {
	_, ok := qs.data[key.name]
	if ok {
		q := qs.data[key.name]
		return q
	}
	qs.data[key.name] = &queue.Queue{}
	q := qs.data[key.name]
	return q
}

func (qs *QueueStorage) Get(key Key) (*queue.Queue, error) {
	_, ok := qs.data[key.name]
	if ok {
		q := qs.data[key.name]
		return q, nil
	}
	return nil, fmt.Errorf("queue not found")
}

func (qs *QueueStorage) Delete(key Key) (bool, error) {
	_, ok := qs.data[key.name]
	if ok {
		delete(qs.data, key.name)
		return true, nil
	}
	return false, fmt.Errorf("queue not found")
}

func (qs *QueueStorage) FlushStorage() {
	for k := range qs.data {
		delete(qs.data, k)
	}
}

var size int
var storage AbstractStorage
var incomeData chan Message
var incomeErrors chan error

func init() {
	size, _ = strconv.Atoi(os.Getenv("KEY_VALUE_WORKERS"))
	storage = &QueueStorage{size, make(map[string]*queue.Queue)}
	incomeData = make(chan Message)
	incomeErrors = make(chan error)
}

func Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case item := <-incomeData:
				q := storage.Set(item.key)
				q.Enqueue(item.value)
				logger.Write(fmt.Sprintf("Put to queue: %s - %s.", item.key.name, item.value.text))
			case <-time.After(time.Second * 5):
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case err := <-incomeErrors:
				logger.Write(fmt.Sprintf("Finished with income error: %s\n", err.Error()))
			case <-time.After(time.Second * 5):
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
}

func Test(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			incomeData <- Message{Key{"queue1"}, Value{"text" + strconv.Itoa(i), time.Now()}}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 150; i++ {
			incomeData <- Message{Key{"queue2"}, Value{"text" + strconv.Itoa(i), time.Now()}}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			incomeData <- Message{Key{"queue3"}, Value{"text" + strconv.Itoa(i), time.Now()}}
		}
	}()
}

func Print(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(5 * time.Second)
		q1, _ := storage.Get(Key{"queue1"})
		fmt.Println("queue1: ")
		fmt.Println(q1.Len())
		q2, _ := storage.Get(Key{"queue2"})
		fmt.Println("queue2: ")
		fmt.Println(q2.Len())
		q3, _ := storage.Get(Key{"queue3"})
		fmt.Println("queue3: ")
		fmt.Println(q3.Len())
	}()
}
