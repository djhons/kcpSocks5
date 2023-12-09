package main

import (
	"github.com/armon/go-socks5"
	kcp "github.com/xtaci/kcp-go/v5"
	"github.com/xtaci/smux"
	"log"
	"time"
)

func main() {
	configPath := "server.ini"
	config := readConfig(configPath)
	createUdp(config)
}
func createUdp(config Config) {
	block, _ := kcp.NewAESBlockCrypt([]byte(config.AesKey))
	lis, err := kcp.ListenWithOptions(config.ServerAddr, block, 10, 3)
	if err != nil {
		log.Println("[-] ", err)
		return
	}
	log.Println("[+] Listen udp ", config.ServerAddr)
	if err := lis.SetDSCP(0); err != nil {
	}
	if err := lis.SetReadBuffer(4194304); err != nil {
	}
	if err := lis.SetWriteBuffer(4194304); err != nil {
	}

	for {
		if conn, err := lis.AcceptKCP(); err == nil {
			log.Println("acc a coonect")
			conn.SetStreamMode(true)
			conn.SetWriteDelay(false)
			conn.SetNoDelay(0, 30, 2, 1)
			conn.SetMtu(1350)
			conn.SetWindowSize(1024, 1024)
			conn.SetACKNoDelay(true)
			cfg := &socks5.Config{}
			if config.UserName != "" || config.PassWord != "" {
				cfg.Credentials = socks5.StaticCredentials(map[string]string{config.UserName: config.PassWord})
			}
			smuxConfig := smux.DefaultConfig()
			smuxConfig.Version = 1
			smuxConfig.MaxReceiveBuffer = 4194304
			smuxConfig.MaxStreamBuffer = 2097152
			smuxConfig.KeepAliveInterval = time.Duration(5) * time.Second

			if err := smux.VerifyConfig(smuxConfig); err != nil {
				log.Fatalf("%+v", err)
			}

			mux, err := smux.Server(conn, smuxConfig)
			defer mux.Close()
			defer conn.Close()
			server, err := socks5.New(cfg)
			if err != nil {
				continue
			}
			stream, err := mux.AcceptStream()
			go func() {
				err = server.ServeConn(stream)
				if err != nil {
					return
				}
			}()
		} else {
			log.Printf("%+v", err)
		}

	}
}
