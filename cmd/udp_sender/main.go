package main

import (
	"flag"
	ossignal "github.com/akley-MK4/go-tools-box/signal"
	"github.com/akley-MK4/simulater-datasource/sender"
	"log"
	"os"
	"syscall"
	"time"
)

// 1万多个packets / 秒
//  upd-sender -dst_ip=172.0.3.22 -dst_port=9991 -max_send_packets_num_sec=600 -max_sender_num=20 -payload_size=256 upd-sender -dst_ip=172.0.3.22 -dst_port=9991 -max_send_packets_num_sec=600 -max_sender_num=20 -payload_size=256

func main() {
	// -max_sender_num=10 -dst_ip=127.0.0.1 -dst_port=8080 -payload_size=1000 -max_send_packets_num_sec=500
	maxSenderNum := flag.Int("max_sender_num", 0, "max_sender_num=1")
	network := flag.String("network", "udp", "network=udp")
	intervalMillisecond := flag.Int("interval_millisecond", 1000, "interval_millisecond=1")
	srcInterfaceName := flag.String("src_interface_name", "", "src_interface_name=eth0")
	dstIp := flag.String("dst_ip", "", "dst_ip=127.0.0.1")
	dstPort := flag.Int("dst_port", 0, "dst_port=8080")
	payloadSize := flag.Int("payload_size", 1024, "payload_size=1024")
	maxSendPacketsNumSec := flag.Int("max_send_packets_num_sec", 600, "max_send_packets_num_sec=600")
	flag.Parse()

	mgr, newMgrErr := sender.NewMgr(*maxSenderNum, *network, time.Duration(*intervalMillisecond)*time.Millisecond,
		*srcInterfaceName, *dstIp, uint16(*dstPort), uint16(*payloadSize), *maxSendPacketsNumSec)
	if newMgrErr != nil {
		log.Println("Failed to create a sender mgr, ", newMgrErr.Error())
		os.Exit(1)
	}

	if err := mgr.Start(); err != nil {
		log.Println("Failed to start the sender mgr, ", err.Error())
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
