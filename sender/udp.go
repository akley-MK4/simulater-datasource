package sender

import (
	"errors"
	PCD "github.com/akley-MK4/pep-coroutine/define"
	PCI "github.com/akley-MK4/pep-coroutine/implement"
	"io"
	"log"
	"time"
)

func NewUDPSender(mgr *Mgr, interval time.Duration, args ...interface{}) (ISender, error) {
	srcInterfaceName := args[0].(string)
	dstIp := args[1].(string)
	dstPort := args[2].(uint16)
	payloadSize := args[3].(uint16)
	maxSendPacketsNumSec := args[4].(int)
	totalMaxSentPacketsNum := args[5].(int)

	if dstPort <= 0 || payloadSize <= 0 {
		return nil, errors.New("invalid parameter")
	}

	if interval <= 0 {
		return nil, errors.New("invalid interval")
	}
	sender := &UDPSender{interval: interval}

	ioWriter, createErr := CreateIOWriter("udp", srcInterfaceName, dstIp, dstPort)
	if createErr != nil {
		return nil, createErr
	}

	co, createCoErr := PCI.CreateCoroutine(PCD.TimerCoroutineType, "UDPSender-sendPeriodically", interval, sender.sendPeriodically)
	if createCoErr != nil {
		return nil, createCoErr
	}
	sender.co = co
	sender.mgr = mgr
	sender.ioWriter = ioWriter
	sender.payloadSize = payloadSize
	sender.maxSendPacketsNumSec = maxSendPacketsNumSec
	sender.totalMaxSentPacketsNum = uint64(totalMaxSentPacketsNum)
	sender.payload = make([]byte, payloadSize)

	return sender, nil
}

type UDPSender struct {
	//srcIP   string
	//dstIp   net.IP
	//dstPort uint16
	mgr *Mgr

	maxSendPacketsNumSec   int
	payloadSize            uint16
	payload                []byte
	interval               time.Duration
	ioWriter               io.Writer
	co                     PCI.ICoroutine
	totalMaxSentPacketsNum uint64
	totalSentPacketsNum    uint64
}

func (t *UDPSender) Start() error {
	if err := PCI.StartCoroutine(t.co); err != nil {
		return err
	}

	return nil
}

func (t *UDPSender) Stop() error {
	if err := PCI.CloseCoroutine(t.co); err != nil {
		return err
	}
	return nil
}

func (t *UDPSender) sendPeriodically(_ PCD.CoId, args ...interface{}) bool {
	statsHandler := t.mgr.GetStatsHandler()
	for i := 0; i < t.maxSendPacketsNumSec; i++ {
		n, writeErr := t.ioWriter.Write(t.payload)
		if writeErr != nil {
			statsHandler.AddTotalFailedSendPacketsNum(1)
			log.Println("Failed to send UDP data, ", writeErr.Error())
			return true
		}

		statsHandler.AddTotalSentPacketsNum(1)
		statsHandler.AddTotalSentPacketsBytes(uint64(n))
		t.totalSentPacketsNum++
		if t.totalMaxSentPacketsNum > 0 && (t.totalSentPacketsNum >= t.totalMaxSentPacketsNum) {
			return false
		}
	}

	return true
}
