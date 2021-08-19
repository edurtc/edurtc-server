package config

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/pion/ion-sfu/pkg/sfu"
	"github.com/spf13/viper"
	"log"
	"time"
)

type Config struct {
	SfuConfig sfu.Config
	NatsConn *nats.Conn
}

func New(filename string, filetype string, filepath string) (*Config, error) {
	conf := sfu.Config{}
	viper.SetConfigName(filename)
	viper.SetConfigType(filetype)
	viper.AddConfigPath(filepath)
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Errorf("config file read failed", err)
		return nil, err
	}
	err = viper.GetViper().Unmarshal(&conf)
	if err != nil {
		fmt.Errorf("sfu config file loaded failed", err)
		return nil, err
	}

	// Connect Options.
	opts := []nats.Option{nats.Name("NATS Sample Queue Subscriber")}
	opts = setupConnOptions(opts)

	nc, err := nats.Connect(nats.DefaultURL, opts...)
	if err != nil {
		fmt.Errorf("nats connection failed: ", err)
		return nil, err
	}
	return &Config{SfuConfig: conf, NatsConn: nc}, nil
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
