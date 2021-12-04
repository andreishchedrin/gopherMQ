package server

import (
	"andreishchedrin/gopherMQ/logger"
	"andreishchedrin/gopherMQ/server/message"
	"andreishchedrin/gopherMQ/service"
	"errors"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

type Grpc struct {
	Port           string
	Logger         logger.AbstractLogger
	MessageService service.AbstractMessageService
}

func (g *Grpc) Serve() error {
	l, err := net.Listen("tcp", ":"+g.Port)
	if err != nil {
		log.Fatal(err)
	}

	pusherServer := message.NewPusherServer{MessageService: g.MessageService}
	pullerServer := message.NewPullerServer{MessageService: g.MessageService}

	grpcServer := grpc.NewServer()
	message.RegisterPusherServer(grpcServer, &pusherServer)
	message.RegisterPullerServer(grpcServer, &pullerServer)

	err = grpcServer.Serve(l)
	return err
}

func (g *Grpc) Shutdown() error {
	//@TODO
	return errors.New("temp")
}

func (g *Grpc) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		g.Serve()
	}()
}

func (g *Grpc) Stop() error {
	//@TODO
	return errors.New("temp")
}

func (g *Grpc) WebsocketListen() {}
