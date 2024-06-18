package echoservice

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/akley-MK4/go-tools-box/netaic"
	"log"
	"net"
	"sync/atomic"
	"time"
)

func NewUDPEchoClient(srcPort uint16, dstAddr string, outputInterface string) (*UDPEchoClient, error) {
	srcIp := "0.0.0.0"
	if outputInterface != "" {
		lIp, errLIp := netaic.GetNetInterfaceAddr(outputInterface)
		if errLIp != nil {
			return nil, errLIp
		}
		srcIp = lIp
	}

	lAddr, errLAddr := net.ResolveUDPAddr("udp", fmt.Sprintf("%v:%d", srcIp, srcPort))
	if errLAddr != nil {
		return nil, errLAddr
	}
	rAddr, errRAddr := net.ResolveUDPAddr("udp", dstAddr)
	if errRAddr != nil {
		return nil, errRAddr
	}

	conn, errConn := net.DialUDP("udp", lAddr, rAddr)
	if errConn != nil {
		return nil, errConn
	}

	client := &UDPEchoClient{
		conn: conn,
	}

	return client, nil
}

type UDPEchoClient struct {
	reqSeqNum uint64
	connected bool
	conn      *net.UDPConn
}

func (t *UDPEchoClient) Start() error {
	go t.receiveMessageForever()

	return nil
}

func (t *UDPEchoClient) receiveMessageForever() {
	readBuf := make([]byte, defaultReadBuffSize)
	for {
		n, rAddr, errRead := t.conn.ReadFromUDP(readBuf)
		if errRead != nil {
			break
		}

		msg := &Message{}
		if err := json.Unmarshal(readBuf[:n], msg); err != nil {
			continue
		}

		if msg.MsgType == MessageTypeRespConnect {
			t.connected = true
			log.Println("Successfully connected to the server")
			continue
		}

		log.Printf("Received a message from remote address %v, SeqNumber: %d, MessageType: %v, MessageSize: %d\n",
			rAddr.String(), msg.SeqNum, msg.MsgType, n)
	}
}

func (t *UDPEchoClient) Connect() error {
	msg := Message{
		MsgType: MessageTypeReqConnect,
	}
	d, _ := json.Marshal(msg)
	_, errWrite := t.conn.Write(d)
	if errWrite != nil {
		return errWrite
	}

	time.Sleep(time.Second * 3)
	if !t.connected {
		return errors.New("timeout")
	}

	return nil
}

func (t *UDPEchoClient) SendMessage(msgType uint16, msgData []byte) (retErr error) {
	if !t.connected {
		return errors.New("disconnected State")
	}

	msg := &Message{
		SeqNum:  atomic.AddUint64(&t.reqSeqNum, 1),
		MsgType: msgType,
		Payload: msgData,
	}
	msgSize := 0

	defer func() {
		if retErr != nil {
			log.Printf("Failed to send a message, SeqNumber: %d, MessageType: %v\n",
				msg.SeqNum, msg.MsgType)
			return
		}
		log.Printf("Sent a message to remote address %v, SeqNumber: %d, MessageType: %v, MessageSize: %d\n",
			t.conn.RemoteAddr().String(), msg.SeqNum, msg.MsgType, msgSize)
	}()

	d, _ := json.Marshal(msg)
	msgSize = len(d)
	_, errWrite := t.conn.Write(d)
	if errWrite != nil {
		return errWrite
	}

	return
}
