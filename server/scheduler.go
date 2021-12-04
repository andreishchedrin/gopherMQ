package server

import (
	"andreishchedrin/gopherMQ/repository"
	"os"
	"strconv"
	"sync"
	"time"
)

var SchedulerExit = make(chan bool)

func (s *FiberServer) StartScheduler(wg *sync.WaitGroup) {
	wg.Add(1)

	timeout, _ := strconv.Atoi(os.Getenv("SCHEDULER_TIMEOUT"))
	go func() {
		defer wg.Done()
		for {
			select {
			case <-SchedulerExit:
				return
			default:
				tasks := s.Repo.GetTasksForWorker()

				if tasks != nil {
					s.ProcessTasks(tasks)
				}

				time.Sleep(time.Duration(timeout) * time.Second)
			}
		}
	}()
}

func (s *FiberServer) ProcessTasks(tasks []repository.Task) {
	for _, task := range tasks {
		switch task.Name {
		case "broadcast":
			pusher := &Pusher{Channel: task.Channel, Message: task.Message}
			broadcastMessage <- pusher
		case "queue":
			s.Storage.Push(task.Channel, task.Message)
		case "persist":
			s.Repo.InsertMessage([]interface{}{task.Channel, task.Message}...)
		default:
			continue
		}

		if task.Repeatable == 0 {
			s.Repo.DeleteTask([]interface{}{task.Name}...)
		}
	}
}
