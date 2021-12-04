package message

import (
	"andreishchedrin/gopherMQ/service"
	"context"
)

type NewPusherServer struct {
	MessageService service.AbstractMessageService
}

func (n *NewPusherServer) Push(ctx context.Context, req *PushStruct) (*PushResponse, error) {
	n.MessageService.Push(req.GetChannel(), req.GetMessage())
	return &PushResponse{Code: 200}, nil
}

func (n *NewPusherServer) mustEmbedUnimplementedPusherServer() {}
