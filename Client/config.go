package main

import (
	"github.com/go-ini/ini"
	"log"
)

type ClientConfig struct {
	ServerAddr string `ini:"ServerAddr"`
	SocksAddr  string `ini:"SocksAddr"`
	AesKey     string `ini:"AesKey"`
}

func readConfig(configPath string) ClientConfig {

	cfg, err := ini.Load(configPath)
	if err != nil {
		log.Fatalf("[-] Error reading config file: %s", err)
		return ClientConfig{}
	}

	var clientConfig ClientConfig
	if err := cfg.Section("Client").MapTo(&clientConfig); err != nil {
		log.Fatalf("[-] Error mapping Server section to struct: %s", err)
		return ClientConfig{}
	}
	return clientConfig
}
