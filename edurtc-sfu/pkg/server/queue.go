package server

import (
	"fmt"
	"github.com/edurtc/edurtc-server/edurtc-sfu/config"
	"github.com/nats-io/nats.go"
	"log"
)

func (n *NatsSignal) listenQueue(conf *config.Config) {
	n.wg.Add(10)

	// Create a queue subscription on "updates" with queue name "workers"
	if _, err := n.natsconn.QueueSubscribe("connection", "job_workers", func(m *nats.Msg) {
		fmt.Println(string(m.Data))
		n.natsconn.Publish("node_name", []byte(conf.ServerName.String()))
		n.wg.Done()
	}); err != nil {
		log.Fatal(err)
	}


	n.natsconn.Flush()
	if err := n.natsconn.LastError(); err != nil {
		log.Fatal(err)
	}

	// Wait for messages to come in
	n.wg.Wait()
}
