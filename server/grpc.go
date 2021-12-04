package server

import (
	api "andreishchedrin/gopherMQ/server/message"
	"context"
	"errors"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
)

type GrpcServer struct {
	Port string
}

func (g *GrpcServer) Serve() error {
	l, err := net.Listen("tcp", ":"+g.Port)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	api.
		err = grpcServer.Serve(l)
	return err
}

func (g *GrpcServer) Shutdown() error {
	return errors.New("temp")
}

func (g *GrpcServer) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		g.Serve()
	}()
}

func (g *GrpcServer) Stop() error {
	return errors.New("temp")
}

func (g *GrpcServer) WebsocketListen() {}

func (g *GrpcServer) Push(ctx context.Context, req *api.PushStruct) (*api.Response, error) {
	return &api.Response{Code: 200}, nil
}
