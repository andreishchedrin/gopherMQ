package message

import (
	"andreishchedrin/gopherMQ/repository"
	"andreishchedrin/gopherMQ/service"
	"context"
	"io"
)

type NewPusherServer struct {
	MessageService service.AbstractMessageService
	Repo           repository.AbstractRepository
}

func (n *NewPusherServer) Push(ctx context.Context, req *PushStruct) (*PushResponse, error) {
	n.MessageService.Push(req.GetChannel(), req.GetMessage())
	return &PushResponse{Code: 200}, nil
}

func (n *NewPusherServer) Publish(ctx context.Context, req *PushStruct) (*PushResponse, error) {
	n.Repo.InsertMessage([]interface{}{req.GetChannel(), req.GetMessage()}...)
	return &PushResponse{Code: 200}, nil
}

func (n *NewPusherServer) Broadcast(stream Pusher_BroadcastServer) error {
	for {
		value, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&PushResponse{Code: 200})
		}

		if err != nil {
			return err
		}

		if pair, ok := Exchange[value.GetChannel()]; ok {
			pair <- value
		} else {
			Exchange[value.GetChannel()] = make(chan *PushStruct, 1000)
			pair = Exchange[value.GetChannel()]
			pair <- value
		}
	}
}

func (n *NewPusherServer) mustEmbedUnimplementedPusherServer() {}
