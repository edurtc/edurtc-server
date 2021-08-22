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
	peers map[string]*sfu.PeerLocal
	natsconn *nats.Conn
	mutex *sync.Mutex
	emission.Emitter
	wg sync.WaitGroup
	config *config.Config
}

func NewNatsSignal(s *sfu.SFU, n *nats.Conn, c *config.Config) *NatsSignal {
	return &NatsSignal{
		sfu:     s,
		peers: make(map[string]*sfu.PeerLocal),
		natsconn: n,
		mutex:   new(sync.Mutex),
		Emitter: *emission.NewEmitter(),
		wg:      sync.WaitGroup{},
		config: c,
	}
}

func (n *NatsSignal) StartServer()  {
	n.StartListen()
	fmt.Println("started", n.config.ServerName)
	n.wg.Add(1)
	defer n.natsconn.Close()
	n.listenQueue()
	go func() {
		n.listenSubscribe()
		n.wg.Done()
	}()
	n.wg.Wait()
}

