package main

import (
	"flag"
	ossignal "github.com/akley-MK4/go-tools-box/signal"
	"github.com/akley-MK4/simulater-datasource/echoservice"
	"log"
	"os"
	"syscall"
)

// -srcPort=9991 -msgQueueCapacity=5

func main() {
	srcPort := flag.Int("srcPort", 9990, "srcPort=9990")
	dstAddr := flag.String("dstAddr", "", "dstAddr=127.0.0.1:9991")

	flag.Parse()

	server, errServer := echoservice.NewUDPEchoServer(uint16(*srcPort), *dstAddr)
	if errServer != nil {
		log.Println("Failed to new the server ", errServer.Error())
		os.Exit(1)
	}

	if err := server.Start(); err != nil {
		log.Println("Failed to start the server ", err.Error())
		os.Exit(1)
	}

	signalHandler := &ossignal.Handler{}
	if err := signalHandler.InitSignalHandler(1); err != nil {
		log.Printf("Failed to initialize process signal handler, %v\n", err)
		os.Exit(1)
		return
	}
	for _, sig := range []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT} {
		signalHandler.RegisterSignal(sig, func() {
			signalHandler.CloseSignalHandler()
		})
	}

	log.Println("The app is running")
	signalHandler.ListenSignal()
}
