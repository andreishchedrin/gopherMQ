package storage

import (
	"andreishchedrin/gopherMQ/logger"
	"fmt"
	"github.com/golang-collections/collections/queue"
	"strconv"
	"sync"
	"time"
)

type AbstractStorage interface {
	Set(key Key) *queue.Queue
	Get(key Key) (*queue.Queue, error)
	Delete(key Key) (bool, error)
	FlushStorage()
}

func (qs *QueueStorage) Set(key Key) *queue.Queue {
	_, ok := qs.Data[key.Name]
	if ok {
		q := qs.Data[key.Name]
		return q
	}
	qs.Data[key.Name] = queue.New()
	q := qs.Data[key.Name]
	return q
}

func (qs *QueueStorage) Get(key Key) (*queue.Queue, error) {
	_, ok := qs.Data[key.Name]
	if ok {
		q := qs.Data[key.Name]
		return q, nil
	}
	return nil, fmt.Errorf("queue not found")
}

func (qs *QueueStorage) Delete(key Key) (bool, error) {
	_, ok := qs.Data[key.Name]
	if ok {
		delete(qs.Data, key.Name)
		return true, nil
	}
	return false, fmt.Errorf("queue not found")
}

func (qs *QueueStorage) FlushStorage() {
	for k := range qs.Data {
		delete(qs.Data, k)
	}
}

var Storage AbstractStorage
var PushData chan Message

func init() {
	Storage = &QueueStorage{make(map[string]*queue.Queue)}
	PushData = make(chan Message)
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

func Test(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			PushData <- Message{Key{"queue1"}, Value{"text" + strconv.Itoa(i), time.Now()}}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 150; i++ {
			PushData <- Message{Key{"queue2"}, Value{"text" + strconv.Itoa(i), time.Now()}}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			PushData <- Message{Key{"queue3"}, Value{"text" + strconv.Itoa(i), time.Now()}}
		}
	}()
}

func Print(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(5 * time.Second)
		q1, _ := Storage.Get(Key{"queue1"})
		fmt.Println("queue1: ")
		fmt.Println(q1.Len())
		q2, _ := Storage.Get(Key{"queue2"})
		fmt.Println("queue2: ")
		fmt.Println(q2.Len())
		q3, _ := Storage.Get(Key{"queue3"})
		fmt.Println("queue3: ")
		fmt.Println(q3.Len())
	}()
}
