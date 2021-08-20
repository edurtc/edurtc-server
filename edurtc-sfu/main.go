package main

import (
	"fmt"
	"github.com/edurtc/edurtc-server/edurtc-sfu/config"
	"github.com/edurtc/edurtc-server/edurtc-sfu/pkg/server"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/pion/ion-sfu/pkg/sfu"
	"log"
	"time"
)

func main() {
	// Connect Options.
	opts := []nats.Option{nats.Name("NATS Sample Queue Subscriber")}
	opts = setupConnOptions(opts)
	nc, err := nats.Connect(nats.DefaultURL, opts...)
	if err != nil {
		fmt.Errorf("nats connection failed: ", err)
		return
	}
	servername := uuid.New()

	conf, err := config.New("config", "toml", "./", servername)
	if err != nil {
		log.Fatal("error reading config", err)
	}
	conf.SfuConfig.WebRTC.SDPSemantics = "unified-plan-with-fallback"

	s := sfu.NewSFU(conf.SfuConfig)
	s.NewDatachannel(sfu.APIChannelLabel)
	ns := server.NewNatsSignal(s, nc, conf)
	ns.StartServer()
}

func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Printf("Disconnected due to: %s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatalf("Exiting: %v", nc.LastError())
	}))
	return opts
}

