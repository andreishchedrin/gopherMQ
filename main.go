package main

import (
	_ "andreishchedrin/gopherMQ/config"
	"andreishchedrin/gopherMQ/storage"
	"sync"
)

var wg sync.WaitGroup

func main() {
	storage.Start(&wg)
	storage.Test(&wg)
	storage.Print(&wg)
	wg.Wait()
}
