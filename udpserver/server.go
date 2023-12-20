package udpserver

import (
	"fmt"
	"log"
	"net"
	"runtime/debug"
)

const (
	defIOQueueThresholdLimit = 10 // queue capacity / IOQueueThreshold
	readBufSize              = 1024 * 82
)

type UDPPacket2 struct {
	SrcAddr string
	DstAddr string
	Payload interface{}
}

type UDPServer struct {
	listenAddr *net.UDPAddr
	conn       *net.UDPConn
	readBuf    []byte
}

func (t *UDPServer) Initialize(addr string) error {
	lAddr, err := net.ResolveUDPAddr("", addr)
	if err != nil {
		return err
	}

	t.readBuf = make([]byte, readBufSize)
	t.listenAddr = lAddr
	return nil
}

func (t *UDPServer) Start() error {
	if err := t.Listen(); err != nil {
		return err
	}
	t.watchReader()

	return nil
}

func (t *UDPServer) Stop() error {
	return nil
}

func (t *UDPServer) Listen() error {
	conn, err := net.ListenUDP("udp", t.listenAddr)
	if err != nil {
		return fmt.Errorf("failed to listen on %v, %v", t.listenAddr.String(), err)
	}

	t.conn = conn
	return nil
}

func (t *UDPServer) watchReader() {
	go t.readForever()
}

func (t *UDPServer) readForever() {
	for {
		if t.read() {
			break
		}
	}

}

func (t *UDPServer) read() bool {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Catch the exception, recover: %v, stack: %v\n", r, string(debug.Stack()))
		}
	}()

	_, _, readErr := t.conn.ReadFromUDP(t.readBuf)
	if readErr != nil {
		log.Printf("The server connection has failed to read, %v\n", readErr)
		return false
	}

	return false
}
