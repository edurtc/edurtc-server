package server

import (
	"fmt"
	"github.com/edurtc/edurtc-server/edurtc/config"
	"github.com/nats-io/nats.go"
	"log"
	"os"
	"os/signal"
)

func StartServer(conf *config.Config)  {
	fmt.Println("started")
	conf.NatsConn.QueueSubscribe("connection", "job_workers", func(msg *nats.Msg) {
		fmt.Println(string(msg.Data))
	})
	conf.NatsConn.Flush()
	if err := conf.NatsConn.LastError(); err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println()
	log.Printf("Draining...")
	conf.NatsConn.Drain()
	log.Fatalf("Exiting")
}


