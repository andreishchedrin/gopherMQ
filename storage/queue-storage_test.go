package storage

import (
	"andreishchedrin/gopherMQ/logger"
	"github.com/golang-collections/collections/queue"
	"sync"
	"testing"
)

func TestQueueStorage(t *testing.T) {
	loggerInstance := &logger.Logger{}
	storageInstance := QueueStorage{
		Data:   make(map[string]*queue.Queue),
		Logger: loggerInstance,
		Debug:  0,
	}

	testChannel := "test_channel123"
	testMessage := "test_message"
	q := storageInstance.Set(Key{testChannel})
	q.Enqueue(testMessage)

	q, _ = storageInstance.Get(Key{testChannel})
	message := q.Dequeue()

	if message != testMessage {
		t.Errorf("got %v, want %v", message, testMessage)
	}

	message = q.Dequeue()

	if message != nil {
		t.Errorf("got %v, want %v", message, nil)
	}

	result, _ := storageInstance.Delete(Key{testChannel})

	if result != true {
		t.Errorf("got %v, want %v", result, true)
	}

	q = storageInstance.Set(Key{testChannel})
	q.Enqueue(testMessage)

	storageInstance.Flush()
	q, _ = storageInstance.Get(Key{testChannel})

	if q != nil {
		t.Errorf("got %v, want %v", q, nil)
	}

	var wg sync.WaitGroup
	storageInstance.Start(&wg)

	storageInstance.Push(testChannel, testMessage)
	message, _ = storageInstance.Pull(testChannel)

	if message != testMessage {
		t.Errorf("got %v, want %v", message, testMessage)
	}
}
