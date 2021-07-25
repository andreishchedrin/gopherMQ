package main

import (
	_ "andreishchedrin/gopherMQ/config"
	"andreishchedrin/gopherMQ/server"
	"andreishchedrin/gopherMQ/storage"
	"sync"
)

var wg sync.WaitGroup

func main() {
	server.Start(&wg)
	storage.Start(&wg)
	storage.Test()
	storage.Print()
	wg.Wait()
}
