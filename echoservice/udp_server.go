package echoservice

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
)

type Client struct {
	Id   uint32
	conn *net.UDPConn
}

func (t *Client) start() error {
	go t.receiveMessageForever()
	return nil
}

func (t *Client) receiveMessageForever() {
	readBuf := make([]byte, defaultReadBuffSize)
	for {
		n, _, errRead := t.conn.ReadFromUDP(readBuf)
		if errRead != nil {
			break
		}

		msg := &Message{}
		if err := json.Unmarshal(readBuf[:n], msg); err != nil {
			continue
		}

		log.Printf("Receive a message from remote address %v, ClientId: %v, SeqNumber: %d, MessageType: %v\n",
			t.conn.RemoteAddr().String(), t.Id, msg.SeqNum, msg.MsgType)

		if err := t.handleMessage(msg); err != nil {
			log.Printf("Failed to handle the message, ClientId: %v, SeqNumber: %d, MessageType: %v\n",
				t.Id, msg.SeqNum, msg.MsgType)
		}
	}
}

func (t *Client) handleMessage(msg *Message) error {
	var respData []byte
	switch msg.MsgType {
	case MessageTypeEcho:
		respData, _ = json.Marshal(msg)
		break
	default:
		return fmt.Errorf("unkonw message type %v", msg.MsgType)
	}

	_, errWrite := t.conn.Write(respData)
	if errWrite != nil {
		return errWrite
	}

	log.Printf("Responded the message to remote address %v, ClientId: %v, SeqNumber: %d, MessageType: %v\n",
		t.conn.RemoteAddr().String(), t.Id, msg.SeqNum, msg.MsgType)
	return nil
}

func NewUDPEchoServer(srcPort uint16, dstAddr string) (*UDPEchoServer, error) {
	lAddr, errLAddr := net.ResolveUDPAddr("udp", fmt.Sprintf("0.0.0.0:%d", srcPort))
	if errLAddr != nil {
		return nil, errLAddr
	}
	//rAddr, errRAddr := net.ResolveUDPAddr("udp", dstAddr)
	//if errRAddr != nil {
	//	return nil, errRAddr
	//}

	conn, errConn := net.ListenUDP("udp", lAddr)
	if errConn != nil {
		return nil, errConn
	}
	//conn, errConn := net.DialUDP("udp", lAddr, rAddr)
	//if errConn != nil {
	//	return nil, errConn
	//}

	server := &UDPEchoServer{
		lAddr:      lAddr,
		listenConn: conn,
		clientMap:  make(map[string]*Client),
	}

	return server, nil
}

type UDPEchoServer struct {
	lAddr      *net.UDPAddr
	listenConn *net.UDPConn

	incClientId    uint32
	clientMap      map[string]*Client
	clientMapMutex sync.RWMutex
}

func (t *UDPEchoServer) Start() error {
	go t.accept()

	return nil
}

func (t *UDPEchoServer) accept() {
	readBuf := make([]byte, 1024*1024)
	for {
		n, addr, errRead := t.listenConn.ReadFromUDP(readBuf)
		if errRead != nil {
			break
		}

		msg := &Message{}
		if err := json.Unmarshal(readBuf[:n], msg); err != nil {
			continue
		}

		if _, err := t.listenConn.WriteToUDP(readBuf[:n], addr); err != nil {
			log.Printf("Failed to Respond the message to remote address %v, SeqNumber: %d, MessageType: %v, Err: %v\n",
				addr.String(), msg.SeqNum, msg.MsgType, err)
			continue
		}

		log.Printf("Responded the message to remote address %v, SeqNumber: %d, MessageType: %v\n",
			addr.String(), msg.SeqNum, msg.MsgType)

		continue

		if msg.MsgType != MessageTypeReqConnect {
			continue
		}

		t.incClientId++
		newClient, errClient := t.addClient(t.incClientId, addr)
		if errClient != nil {
			log.Println("Failed to add a client, ", errClient.Error())
			continue
		}
		if err := newClient.start(); err != nil {
			log.Println("Failed to start the client, ", err.Error())
			continue
		}
		log.Printf("Started the client, Id: %v, Addr: %v\n", newClient.Id, addr.String())

		respMsg := Message{
			MsgType: MessageTypeRespConnect,
		}
		d, _ := json.Marshal(respMsg)
		if _, err := newClient.conn.Write(d); err != nil {
			log.Printf("Failed to respond to a connection message, Id: %v, Addr: %v\n",
				newClient.Id, addr.String())
		}
	}
}

func (t *UDPEchoServer) addClient(id uint32, addr *net.UDPAddr) (*Client, error) {
	t.clientMapMutex.Lock()
	defer t.clientMapMutex.Unlock()

	existClient, exist := t.clientMap[addr.String()]
	if exist {
		return existClient, nil
	}

	conn, errConn := net.DialUDP("udp", t.lAddr, addr)
	if errConn != nil {
		return nil, errConn
	}

	t.clientMap[addr.String()] = &Client{
		Id:   id,
		conn: conn,
	}

	return t.clientMap[addr.String()], nil
}
