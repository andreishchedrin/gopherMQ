package storage

import (
	"fmt"
	"github.com/golang-collections/collections/queue"
	"time"
)

var PushData = make(chan Message)

func (qs *QueueStorage) Start() {
	go func() {
		for {
			select {
			case item := <-PushData:
				qs.mu.Lock()
				q := qs.Set(item.Key)
				q.Enqueue(item.Value)
				qs.mu.Unlock()
				if qs.Debug == 1 {
					qs.Logger.Log(fmt.Sprintf("Put to queue: %s - %s.", item.Key.Name, item.Value.Text))
				}
			default:
				time.Sleep(1 * time.Millisecond)
			}
		}
	}()
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
	q, ok := qs.Data[key.Name]
	if ok {
		return q, nil
	}
	return nil, fmt.Errorf("queue not found")
}

func (qs *QueueStorage) Delete(key Key) (bool, error) {
	qs.mu.Lock()
	defer qs.mu.Unlock()
	_, ok := qs.Data[key.Name]
	if ok {
		delete(qs.Data, key.Name)
		return true, nil
	}
	return false, fmt.Errorf("queue not found")
}

func (qs *QueueStorage) Flush() {
	qs.mu.Lock()
	defer qs.mu.Unlock()
	for k := range qs.Data {
		delete(qs.Data, k)
	}
}

func (qs *QueueStorage) Push(name string, message string) {
	PushData <- Message{Key{name}, Value{message, time.Now()}}
}

func (qs *QueueStorage) Pull(name string) (string, error) {
	qs.mu.Lock()
	q, err := qs.Get(Key{name})
	if err != nil {
		return err.Error(), err
	}

	if q == nil {
		return "Queue is not exists.", nil
	}

	if q.Len() == 0 {
		return "Queue is empty.", nil
	}

	res := q.Dequeue()
	qs.mu.Unlock()

	return res.(Value).Text, nil
}
