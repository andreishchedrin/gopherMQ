package server

import (
	"sync"
)

type AbstractServer interface {
	Serve() error
	Shutdown() error
	WebsocketListen()
	Start(wg *sync.WaitGroup)
	Stop() error
}
