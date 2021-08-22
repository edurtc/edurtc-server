package server

import (
	"github.com/nats-io/nats.go"
	"log"
)

func (n *NatsSignal) listenSubscribe() {
	ch := make(chan *nats.Msg, 64)
	sub, err := n.natsconn.ChanSubscribe("signal.*.*", ch)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		sub.Unsubscribe()
		sub.Drain()
	}()
	for  {
		msg, ok := <- ch
		if !ok {
			log.Fatal("redis pub/sub is down")
			break
		}
		n.Emit("msg", msg)
	}


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