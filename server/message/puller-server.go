package message

import (
	"andreishchedrin/gopherMQ/service"
	"context"
)

type NewPullerServer struct {
	MessageService service.AbstractMessageService
}

func (n *NewPullerServer) Pull(ctx context.Context, req *PullStruct) (*PullResponse, error) {
	message, err := n.MessageService.Pull(req.GetChannel())
	if err != nil {
		return &PullResponse{Code: 400, Payload: err.Error()}, err
	}

	res := &PullResponse{Code: 200, Payload: message}

	return res, nil
}

func (n *NewPullerServer) mustEmbedUnimplementedPullerServer() {}
