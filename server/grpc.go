package server

import (
	"andreishchedrin/gopherMQ/logger"
	"andreishchedrin/gopherMQ/repository"
	"andreishchedrin/gopherMQ/server/message"
	"andreishchedrin/gopherMQ/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Grpc struct {
	Port           string
	Logger         logger.AbstractLogger
	MessageService service.AbstractMessageService
	Repo           repository.AbstractRepository
	Server         *grpc.Server
}

func (g *Grpc) Serve() {
	l, err := net.Listen("tcp", ":"+g.Port)
	if err != nil {
		log.Fatal(err)
	}

	pusherServer := message.NewPusherServer{MessageService: g.MessageService, Repo: g.Repo}
	pullerServer := message.NewPullerServer{MessageService: g.MessageService, Repo: g.Repo}

	g.Server = grpc.NewServer()
	message.RegisterPusherServer(g.Server, &pusherServer)
	message.RegisterPullerServer(g.Server, &pullerServer)

	go func() {
		err = g.Server.Serve(l)
		if err != nil {
			panic(err)
		}
	}()
}

func (g *Grpc) Shutdown() error {
	g.Server.GracefulStop()
	return nil
}

func (g *Grpc) Start() {
	g.Serve()
}

func (g *Grpc) Stop() error {
	return g.Shutdown()
}
