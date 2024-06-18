package main

import (
	"flag"
	"github.com/akley-MK4/simulater-datasource/echoservice"
	"log"
	"os"
	"time"
)

// -srcPort=9990 -dstAddr=127.0.0.1:9991 -maxSendMsgNumLimit=5 -msgPayloadSize=200

func main() {

	srcPort := flag.Int("srcPort", 9990, "srcPort=9990")
	outputInterface := flag.String("outputInterface", "", "outputInterface=9990")
	dstAddr := flag.String("dstAddr", "", "dstAddr=9990")
	maxSendMsgNumLimit := flag.Int("maxSendMsgNumLimit", 1, "maxSendMsgNumLimit=100")
	msgPayloadSize := flag.Int("msgPayloadSize", 100, "msgPayloadSize=100")

	flag.Parse()

	client, errClient := echoservice.NewUDPEchoClient(uint16(*srcPort), *dstAddr, *outputInterface)
	if errClient != nil {
		log.Println("Failed to new the client ", errClient.Error())
		os.Exit(1)
	}
	if err := client.Start(); err != nil {
		log.Println("Failed to start the client ", errClient.Error())
		os.Exit(1)
	}
	time.Sleep(time.Millisecond * 10)

	//if err := client.Connect(); err != nil {
	//	log.Println("Failed to connect the server ", err.Error())
	//	os.Exit(1)
	//}

	for i := 0; i < *maxSendMsgNumLimit; i++ {
		if err := client.SendMessage(echoservice.MessageTypeEcho, make([]byte, *msgPayloadSize)); err != nil {

		}
		time.Sleep(time.Second * 1)
	}
	time.Sleep(time.Hour)
}
