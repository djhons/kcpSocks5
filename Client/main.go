package main

import (
	"github.com/pkg/errors"
	"github.com/xtaci/kcp-go/v5"
	"github.com/xtaci/smux"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	configPath := "client.ini"
	config := readConfig(configPath)
	createSocks(config)
}
func createSocks(config ClientConfig) {
	socksConn, _ := net.Listen("tcp", config.SocksAddr)
	for {
		socks, err := socksConn.Accept()
		if err != nil {
			continue
		}
		go handleClient(config, socks)
	}
}
func handleClient(config ClientConfig, p1 net.Conn) {
	mux, err := createUdp(config.AesKey, config.ServerAddr)
	if err != nil {
		return
	}
	p2, err := mux.OpenStream()
	defer p1.Close()
	defer mux.Close()
	if err != nil {
		log.Println(err)
		return
	}
	defer p2.Close()
	streamCopy := func(dst io.Writer, src io.ReadCloser) {
		if _, err := Copy(dst, src); err != nil {
		}
		p1.Close()
		p2.Close()
	}
	go streamCopy(p1, p2)
	streamCopy(p2, p1)
}
func Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	// If the reader has a WriteTo method, use it to do the copy.
	// Avoids an allocation and a copy.
	if wt, ok := src.(io.WriterTo); ok {
		return wt.WriteTo(dst)
	}
	// Similarly, if the writer has a ReadFrom method, use it to do the copy.
	if rt, ok := dst.(io.ReaderFrom); ok {
		return rt.ReadFrom(src)
	}

	// fallback to standard io.CopyBuffer
	buf := make([]byte, 4096)
	return io.CopyBuffer(dst, src, buf)
}
func createUdp(aesKey, serverAddr string) (*smux.Session, error) {
	block, _ := kcp.NewAESBlockCrypt([]byte(aesKey))
	kcpconn, err := kcp.DialWithOptions(serverAddr, block, 10, 3)
	if err != nil {
		return nil, errors.Wrap(err, "dial()")
	}
	kcpconn.SetStreamMode(true)
	kcpconn.SetWriteDelay(false)
	kcpconn.SetNoDelay(0, 30, 2, 1)
	kcpconn.SetWindowSize(128, 512)
	kcpconn.SetMtu(1350)
	kcpconn.SetACKNoDelay(true)

	if err := kcpconn.SetDSCP(0); err != nil {
	}
	if err := kcpconn.SetReadBuffer(4194304); err != nil {
	}
	if err := kcpconn.SetWriteBuffer(4194304); err != nil {
	}

	// stream multiplex
	smuxConfig := smux.DefaultConfig()
	smuxConfig.MaxReceiveBuffer = 4194304
	smuxConfig.MaxStreamBuffer = 2097152
	smuxConfig.KeepAliveInterval = time.Duration(5) * time.Second
	mux, err := smux.Client(kcpconn, smuxConfig)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return mux, nil
}
