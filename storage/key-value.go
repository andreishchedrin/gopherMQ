package storage

import (
	"andreishchedrin/gopherMQ/logger"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

type KeyValueStorage struct {
	workerPoolSize int
}

type AbstractStorage interface {
	Push()
	Pull()
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

var data map[Key]Value
var storage KeyValueStorage
var incomeData chan Message
var incomeErrors chan error

func init() {
	data = make(map[Key]Value)
	size, _ := strconv.Atoi(os.Getenv("KEY_VALUE_WORKERS"))
	storage = KeyValueStorage{size}
	incomeData = make(chan Message, storage.workerPoolSize)
	incomeErrors = make(chan error, 1000)
}

func Push() {

}

func Pull() {

}

func Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case item := <-incomeData:
				data[item.key] = item.value
				logger.Write(fmt.Sprintf("Put to data: %s - %s", item.key.name, item.value.text))
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

	//@TODO close channel when shutdown?
	//defer close(incomeData)
	//defer close(incomeErrors)
}

func Test() {
	incomeData <- Message{Key{"key1"}, Value{"text1", time.Now()}}
	fmt.Printf("%v, %T\n", incomeData, incomeData)
	time.Sleep(2 * time.Second)
	incomeData <- Message{Key{"key2"}, Value{"text2", time.Now()}}
	fmt.Printf("%v, %T\n", incomeData, incomeData)
	time.Sleep(2 * time.Second)
	incomeData <- Message{Key{"key3"}, Value{"text3", time.Now()}}
	fmt.Printf("%v, %T\n", incomeData, incomeData)
}

func Print() {
	fmt.Println(data)
}
