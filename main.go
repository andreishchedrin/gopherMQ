package main

import (
	_ "andreishchedrin/gopherMQ/config"
	"andreishchedrin/gopherMQ/server"
	"andreishchedrin/gopherMQ/storage"
	"sync"
)

var wg sync.WaitGroup

func main() {

	go server.Listen()

	server.Start(&wg)

	// @TODO
	defer func() {
		err := server.Stop()
		if err != nil {

		}
	}()

	storage.Start(&wg)
	//storage.Test(&wg)
	//storage.Print(&wg)
	wg.Wait()
}
