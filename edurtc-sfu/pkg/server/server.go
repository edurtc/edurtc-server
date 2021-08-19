package server

import (
	"fmt"
	"github.com/chuckpreslar/emission"
	"github.com/edurtc/edurtc-server/edurtc-sfu/config"
	"github.com/nats-io/nats.go"
	"github.com/pion/ion-sfu/pkg/sfu"
	"log"
	"sync"
)

type NatsSignal struct {
	sfu *sfu.SFU
	NatsConn *nats.Conn
	mutex *sync.Mutex
	emission.Emitter
	wg sync.WaitGroup
}

func NewNatsSignal(sfu *sfu.SFU, n *nats.Conn) *NatsSignal {
	return &NatsSignal{
		sfu:     sfu,
		NatsConn: n,
		mutex:   new(sync.Mutex),
		Emitter: *emission.NewEmitter(),
		wg:      sync.WaitGroup{},
	}
}

func (n *NatsSignal) StartServer(conf *config.Config)  {
	fmt.Println("started", conf.ServerName)
	defer n.NatsConn.Close()
	n.wg.Add(10)

	// Create a queue subscription on "updates" with queue name "workers"
	if _, err := n.NatsConn.QueueSubscribe("connection", "job_workers", func(m *nats.Msg) {
		fmt.Println(string(m.Data))
		n.NatsConn.Publish("node_name", []byte(conf.ServerName.String()))
		n.wg.Done()
	}); err != nil {
		log.Fatal(err)
	}


	n.NatsConn.Flush()
	if err := n.NatsConn.LastError(); err != nil {
		log.Fatal(err)
	}

	// Wait for messages to come in
	n.wg.Wait()
}

