package storage

import (
	"andreishchedrin/gopherMQ/logger"
	"bufio"
	"bytes"
	"github.com/golang-collections/collections/queue"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
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

func BenchmarkQueueStruct(b *testing.B) {
	b.ReportAllocs()

	loggerInstance := &logger.Logger{}
	storageInstance := QueueStorage{
		Data:   make(map[string]*queue.Queue),
		Logger: loggerInstance,
		Debug:  0,
	}

	var wg sync.WaitGroup
	storageInstance.Start(&wg)

	wg.Add(1)
	go func(q *QueueStorage) {
		defer wg.Done()
		for i := 0; i < b.N; i++ {
			wg.Add(1)
			go func(i int, q *QueueStorage) {
				defer wg.Done()
				var s1 strings.Builder
				s1.WriteString("test_channel")
				s1.WriteString(strconv.Itoa(i))
				testChannel := s1.String()

				for i := 0; i < b.N; i++ {
					var s2 strings.Builder
					s2.WriteString("test_message")
					s2.WriteString(strconv.Itoa(i))
					testMessage := s2.String()
					q.Push(testChannel, testMessage)
				}
			}(i, q)
		}
	}(&storageInstance)

	time.Sleep(10 * time.Second)
	wg.Add(1)
	go func(q *QueueStorage) {
		defer wg.Done()
		for i := 0; i < b.N; i++ {
			wg.Add(1)
			go func(i int, q *QueueStorage) {
				defer wg.Done()
				var s1 strings.Builder
				s1.WriteString("test_channel")
				s1.WriteString(strconv.Itoa(i))
				testChannel := s1.String()

				for i := 0; i < b.N; i++ {
					q.Pull(testChannel)
				}
			}(i, q)
		}
	}(&storageInstance)

	//hardcoded storage.Start defer wg.Done()
	wg.Done()
	wg.Wait()
}

func BenchmarkQueueBuffer(b *testing.B) {
	b.ReportAllocs()

	data := make(map[string]*bytes.Buffer)

	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(1)
	go func(data map[string]*bytes.Buffer, mu *sync.Mutex) {
		defer wg.Done()
		for i := 0; i < b.N; i++ {
			wg.Add(1)
			go func(i int, data map[string]*bytes.Buffer, mu *sync.Mutex) {
				defer wg.Done()
				var s1 strings.Builder
				s1.WriteString("test_channel")
				s1.WriteString(strconv.Itoa(i))
				testChannel := s1.String()

				var s2 strings.Builder
				buf := new(bytes.Buffer)
				for i := 0; i < b.N; i++ {
					s2.WriteString("test_message")
					s2.WriteString(strconv.Itoa(i))
					s2.WriteString("\n")
					testMessage := s2.String()

					buf.WriteString(testMessage)
					mu.Lock()
					data[testChannel] = buf
					mu.Unlock()
				}
			}(i, data, mu)
		}
	}(data, &mu)

	var mu2 sync.Mutex

	time.Sleep(10 * time.Second)
	wg.Add(1)
	go func(data map[string]*bytes.Buffer, mu *sync.Mutex) {
		defer wg.Done()
		for i := 0; i < b.N; i++ {
			wg.Add(1)
			go func(i int, data map[string]*bytes.Buffer, mu *sync.Mutex) {
				defer wg.Done()
				var s1 strings.Builder
				s1.WriteString("test_channel")
				s1.WriteString(strconv.Itoa(i))
				testChannel := s1.String()

				for i := 0; i < b.N; i++ {
					writer := data[testChannel]
					str := strings.NewReader(writer.String())
					r := bufio.NewReader(str)

					result, _, _ := r.ReadLine()

					newStr := writer.String()
					rpl := strings.Replace(newStr, string(result)+"\n", "", 1)
					newWriter := new(bytes.Buffer)
					newWriter.WriteString(rpl)

					mu.Lock()
					data[testChannel] = newWriter
					mu.Unlock()
				}
			}(i, data, mu)
		}
	}(data, &mu2)

	wg.Wait()
}
