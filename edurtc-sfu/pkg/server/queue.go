package server

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"log"
	"strings"
)

func (n *NatsSignal) listenQueue() {
	// Create a queue subscription on "updates" with queue name "workers"
	//n.wg.Add(1)
	if _, err := n.natsconn.QueueSubscribe("connection.*", "job_workers", func(m *nats.Msg) {
		fmt.Println(string(m.Data))
		room := getRoom(m.Subject)
		fmt.Println("room: ", room)
		rid, _ := uuid.Parse(room)
		token, err := n.ClaimToken(rid)
		if err != nil {
			log.Fatal("failed to claim jwt: ", err)
		}
		fmt.Println("token: ", token)
		n.natsconn.Publish("node_name", []byte(n.config.ServerName.String()+":"+token))
		//n.wg.Done()
	}); err != nil {
		log.Fatal(err)
	}

	n.natsconn.Flush()
	if err := n.natsconn.LastError(); err != nil {
		log.Fatal(err)
	}

	//c := make(chan os.Signal, 1)
	//signal.Notify(c, os.Interrupt)
	//<-c
	//log.Println()
	//log.Printf("Draining...")
	//n.natsconn.Drain()
	//log.Fatalf("Exiting")

	// Wait for messages to come in
	//n.wg.Wait()
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
func getUid(uid string) string {
	delimiter := "."
	result := strings.Split(uid, delimiter)[2]
	return result
}
func getToken(room string) string {
	delimiter := "."
	result := strings.Split(room, delimiter)[3]
	return result
}
