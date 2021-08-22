package config

import (
	"gopkg.in/ini.v1"
	"log"
	"os"
)

type JwtConfig struct {
	Secret []byte
}

func NewJwtConfig() *JwtConfig {
	return &JwtConfig{
		Secret: []byte(GetIni("jwt_secret", "JWT_SECRET", "awesome")),
	}
}

func GetIni(section, key, defaultValue string) string {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatal("Fail to read file: %v", err)
		os.Exit(1)
	}

	if value := cfg.Section(section).Key(key).String(); value != "" {
		return value
	}
	return defaultValue
}