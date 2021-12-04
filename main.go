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

	container.CleanerInstance.StartCleaner(&wg)
	defer container.CleanerInstance.StopCleaner()

	container.StorageInstance.Start(&wg)

	go container.ServerInstance.WebsocketListen()

	container.ServerInstance.Start(&wg)
	defer container.ServerInstance.Stop()

	container.SchedulerInstance.StartScheduler(&wg)
	defer container.SchedulerInstance.StopScheduler()

	wg.Wait()
}
