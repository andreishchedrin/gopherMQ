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
	defer db.Close()

	go server.WebsocketListen()

	server.Start(&wg)
	defer server.Stop()

	storage.Start(&wg)
	//storage.Test(&wg)
	//storage.Print(&wg)
	wg.Wait()
}
