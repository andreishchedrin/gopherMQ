package main

import (
	_ "andreishchedrin/gopherMQ/config"
	"andreishchedrin/gopherMQ/db"
	"andreishchedrin/gopherMQ/server"
	"andreishchedrin/gopherMQ/storage"
	"sync"
)

var wg sync.WaitGroup

func main() {
	db.Prepare()

	go server.WebsocketListen()

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
