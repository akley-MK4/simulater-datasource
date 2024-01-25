package main

import (
	"flag"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
	"os"
	"time"
)

func main() {

	filterStr := flag.String("filter_str", "", "filter_str=")
	intervalPrintStatsSec := flag.Int("interval_print_stats_sec", 5, "interval_print_stats_sec=")
	flag.Parse()

	handler, errOpen := pcap.OpenLive("eth0", 65536, false, pcap.BlockForever)
	if errOpen != nil {
		log.Println("Failed to open live, ", errOpen)
		os.Exit(1)
	}

	if err := handler.SetBPFFilter(*filterStr); err != nil {
		log.Println("Failed to set BPF Filter, ", err)
		os.Exit(1)
	}

	totalReceivedPacketsNum := uint64(0)
	totalReceivedPacketsNumPrinted := uint64(0)

	go func() {
		for {
			time.Sleep(time.Second * time.Duration(*intervalPrintStatsSec))
			if totalReceivedPacketsNum == totalReceivedPacketsNumPrinted {
				continue
			}
			totalReceivedPacketsNumPrinted = totalReceivedPacketsNum
			log.Printf("TotalReceivedPacketsNum=%d\n", totalReceivedPacketsNumPrinted)

		}
	}()

	packetSource := gopacket.NewPacketSource(handler, handler.LinkType())
	for packet := range packetSource.Packets() {
		if 1 == 0 {
			println(packet)
		}
		totalReceivedPacketsNum++
	}
}
