package server

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"strings"
)

func (n *NatsSignal) listenQueue() {
	n.wg.Add(10)

	// Create a queue subscription on "updates" with queue name "workers"
	if _, err := n.natsconn.QueueSubscribe("connection.*", "job_workers", func(m *nats.Msg) {
		fmt.Println(string(m.Data))
		room := getRoom(m.Subject)
		fmt.Println("room: ", room)
		n.natsconn.Publish("node_name", []byte(n.config.ServerName.String()))
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

func getMethod(channel string) string {
	delimiter := "."
	result := strings.Join(strings.Split(channel, delimiter)[3:], delimiter)
	return result
}
func getRoom(room string) string {
	delimiter := "."
	result := strings.Split(room, delimiter)[1]
	return result
}
