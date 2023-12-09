package main

import (
	"github.com/go-ini/ini"
	"log"
)

type Config struct {
	Name       string
	DeleteMe   bool
	AesKey     string
	UserName   string
	PassWord   string
	ServerAddr string
}

func readConfig(configPath string) Config {

	cfg, err := ini.Load(configPath)
	if err != nil {
		log.Fatalf("[-] Error reading config file: %s", err)
		return Config{}
	}

	var clientConfig Config
	if err := cfg.Section("Server").MapTo(&clientConfig); err != nil {
		log.Fatalf("[-] Error mapping Server section to struct: %s", err)
		return Config{}
	}
	return clientConfig
}
