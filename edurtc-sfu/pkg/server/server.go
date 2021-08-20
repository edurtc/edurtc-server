package server

import (
	"fmt"
	"github.com/chuckpreslar/emission"
	"github.com/edurtc/edurtc-server/edurtc-sfu/config"
	"github.com/nats-io/nats.go"
	"github.com/pion/ion-sfu/pkg/sfu"
	"sync"
)

type NatsSignal struct {
	sfu *sfu.SFU
	natsconn *nats.Conn
	mutex *sync.Mutex
	emission.Emitter
	wg sync.WaitGroup
	config *config.Config
}

func NewNatsSignal(sfu *sfu.SFU, n *nats.Conn, c *config.Config) *NatsSignal {
	return &NatsSignal{
		sfu:     sfu,
		natsconn: n,
		mutex:   new(sync.Mutex),
		Emitter: *emission.NewEmitter(),
		wg:      sync.WaitGroup{},
		config: c,
	}
}

func (n *NatsSignal) StartServer()  {
	fmt.Println("started", n.config.ServerName)
	defer n.natsconn.Close()
	n.listenQueue()
}

