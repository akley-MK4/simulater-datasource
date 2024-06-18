package echoservice

type Message struct {
	MsgType uint16
	SeqNum  uint64
	Payload []byte
}

const (
	MessageTypeReqConnect uint16 = iota + 1
	MessageTypeRespConnect
	MessageTypeEcho
)
