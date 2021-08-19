package config

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/pion/ion-sfu/pkg/sfu"
	"github.com/spf13/viper"
)

type Config struct {
	SfuConfig sfu.Config
	ServerName uuid.UUID
}

func New(filename string, filetype string, filepath string, servername uuid.UUID) (*Config, error) {
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


	return &Config{SfuConfig: conf, ServerName: servername}, nil
}
