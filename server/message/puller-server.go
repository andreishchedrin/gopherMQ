package message

import (
	"andreishchedrin/gopherMQ/repository"
	"andreishchedrin/gopherMQ/service"
	"context"
	"google.golang.org/grpc/peer"
)

type NewPullerServer struct {
	MessageService service.AbstractMessageService
	Repo           repository.AbstractRepository
}

func (n *NewPullerServer) Pull(ctx context.Context, req *PullStruct) (*PullResponse, error) {
	message, err := n.MessageService.Pull(req.GetChannel())
	if err != nil {
		return &PullResponse{Code: 400, Payload: err.Error()}, err
	}

	res := &PullResponse{Code: 200, Payload: message}

	return res, nil
}

func (n *NewPullerServer) Consume(ctx context.Context, req *PullStruct) (*PullResponse, error) {
	p, _ := peer.FromContext(ctx)
	clientId := n.Repo.InsertClient([]interface{}{p.Addr.String(), req.GetChannel()}...)
	messageId, messagePayload := n.Repo.SelectMessage([]interface{}{req.GetChannel(), clientId}...)
	n.Repo.InsertClientMessage([]interface{}{clientId, messageId}...)

	res := &PullResponse{Code: 200, Payload: messagePayload}

	return res, nil
}

func (n *NewPullerServer) Ws(req *PullStruct, stream Puller_WsServer) error {
	if pair, ok := Exchange[req.GetChannel()]; !ok {
		return nil
	} else {
		for value := range pair {

			if err := stream.Send(&PullResponse{Code: 200, Payload: value.GetMessage()}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (n *NewPullerServer) mustEmbedUnimplementedPullerServer() {}
