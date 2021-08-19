package main

import (
	"github.com/edurtc/edurtc-server/edurtc-sfu/config"
	"github.com/edurtc/edurtc-server/edurtc-sfu/pkg/server"
	"github.com/pion/ion-sfu/pkg/sfu"
	"log"
)

func main() {
	conf, err := config.New("config", "toml", "./")
	if err != nil {
		log.Fatal("error reading config", err)
	}
	conf.SfuConfig.WebRTC.SDPSemantics = "unified-plan-with-fallback"

	s := sfu.NewSFU(conf.SfuConfig)
	s.NewDatachannel(sfu.APIChannelLabel)
	server.StartServer(conf)
}
