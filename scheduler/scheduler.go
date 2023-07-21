package scheduler

import (
	"andreishchedrin/gopherMQ/repository"
	"andreishchedrin/gopherMQ/storage"
	"time"
)

type AbstractScheduler interface {
	StartScheduler()
}

type Scheduler struct {
	Repo       repository.AbstractRepository
	Storage    storage.AbstractStorage
	Timeout    int
	ServerMode string
}

func (s *Scheduler) StartScheduler() {
	go func() {
		for {
			select {
			default:
				tasks := s.Repo.GetTasksForWorker()
				if tasks != nil {
					s.ProcessTasks(tasks)
				}
				time.Sleep(time.Duration(s.Timeout) * time.Second)
			}
		}
	}()
}

func (s *Scheduler) ProcessTasks(tasks []repository.Task) {
	for _, task := range tasks {
		switch task.Name {
		//case "broadcast":
		//	if s.ServerMode == "grpc" {
		//		continue
		//	}
		//	pusher := &server.Push{Channel: task.Channel, Message: task.Message}
		//	server.BroadcastMessage <- pusher
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
