package main

import (
	"flag"
	ossignal "github.com/akley-MK4/go-tools-box/signal"
	"github.com/akley-MK4/simulater-datasource/udpserver"
	"log"
	"os"
	"syscall"
)

func main() {
	// -max_sender_num=10 -dst_ip=127.0.0.1 -dst_port=8080 -payload_size=1000 -max_send_packets_num_sec=500
	addr := flag.String("addr", "0.0.0.0:9991", "addr=")
	flag.Parse()

	udpSvr := udpserver.UDPServer{}

	if err := udpSvr.Initialize(*addr); err != nil {
		log.Println("Failed to initialize UDPServer, ", err.Error())
		os.Exit(1)
	}
	if err := udpSvr.Start(); err != nil {
		log.Println("Failed to start UDPServer, ", err.Error())
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
