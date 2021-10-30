package main

import (
	_ "andreishchedrin/gopherMQ/config"
	"andreishchedrin/gopherMQ/container"
	"sync"
)

var wg sync.WaitGroup

func main() {
	container.DbInstance.Prepare()
	defer container.DbInstance.Close()

	go container.ServerInstance.WebsocketListen()

	container.ServerInstance.Start(&wg)
	defer container.ServerInstance.Stop()

	container.StorageInstance.Start(&wg)
	//storage.Test(&wg)
	//storage.Print(&wg)
	wg.Wait()
}
