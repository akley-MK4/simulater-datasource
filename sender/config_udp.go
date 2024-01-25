package sender

import (
	"encoding/json"
	PCD "github.com/akley-MK4/pep-coroutine/define"
	PCI "github.com/akley-MK4/pep-coroutine/implement"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"time"
)

type Config struct {
	MaxConnNum          int               `json:"MaxConnNum"`
	GroupSendConfigList []GroupSendConfig `json:"GroupSendConfigList"`
}

type GroupSendConfig struct {
	BeforeSendWaitRangeMs int `json:"beforeSendWaitRangeMs"`
	//Count                 uint64 `json:"count"`
	SrcInterface         string `json:"srcInterface"`
	DstIp                string `json:"dstIp"`
	DstPort              uint16 `json:"dstPort"`
	PayloadSize          uint16 `json:"payloadSize"`
	MaxSendPacketsNumSec int    `json:"maxSendPacketsNumSec"`
}

func NewConfigUDPSender(mgr *Mgr, interval time.Duration, args ...interface{}) (retSenderList []ISender, retErr error) {
	configPath := args[0].(string)
	data, readErr := ioutil.ReadFile(configPath)
	if readErr != nil {
		retErr = readErr
		return
	}

	var conf Config
	if err := json.Unmarshal(data, &conf); err != nil {
		retErr = err
		return
	}

	randVal := rand.New(rand.NewSource(time.Now().UnixNano()))

	for idx, groupConf := range conf.GroupSendConfigList {
		if groupConf.DstPort <= 0 || groupConf.PayloadSize <= 0 {
			log.Printf("invalid parameter, Idx: %d\n", idx)
			continue
		}

		if interval <= 0 {
			log.Printf("invalid interval, Idx: %d\n", idx)
			continue
		}

		confUDPSender := &ConfigUDPSender{
			interval: interval,
		}
		ioWriter, createErr := CreateIOWriter("udp", groupConf.SrcInterface, groupConf.DstIp, groupConf.DstPort)
		if createErr != nil {
			log.Printf("Failed to IO Writer, %v\n", createErr)
			continue
		}

		co, createCoErr := PCI.CreateCoroutine(PCD.TimerCoroutineType, "UDPSender-sendPeriodically", interval, confUDPSender.sendPeriodically)
		if createCoErr != nil {
			log.Printf("Failed to create co, %v", createErr)
			continue
		}

		confUDPSender.co = co
		confUDPSender.mgr = mgr
		confUDPSender.ioWriter = ioWriter
		confUDPSender.payloadSize = groupConf.PayloadSize
		confUDPSender.maxSendPacketsNumSec = groupConf.MaxSendPacketsNumSec
		confUDPSender.payload = make([]byte, groupConf.PayloadSize)
		confUDPSender.beforeSendWaitTime = time.Millisecond * time.Duration(randVal.Intn(groupConf.BeforeSendWaitRangeMs)+1)

		retSenderList = append(retSenderList, confUDPSender)
	}

	return
}

type ConfigUDPSender struct {
	//srcIP   string
	//dstIp   net.IP
	//dstPort uint16
	mgr *Mgr

	beforeSendWaitTime   time.Duration
	maxSendPacketsNumSec int
	payloadSize          uint16
	payload              []byte
	interval             time.Duration
	ioWriter             io.Writer
	co                   PCI.ICoroutine
}

func (t *ConfigUDPSender) Start() error {
	if err := PCI.StartCoroutine(t.co); err != nil {
		return err
	}

	return nil
}

func (t *ConfigUDPSender) Stop() error {
	if err := PCI.CloseCoroutine(t.co); err != nil {
		return err
	}
	return nil
}

func (t *ConfigUDPSender) sendPeriodically(_ PCD.CoId, args ...interface{}) bool {
	statsHandler := t.mgr.GetStatsHandler()
	for i := 0; i < t.maxSendPacketsNumSec; i++ {

		if t.beforeSendWaitTime > 0 {
			time.Sleep(t.beforeSendWaitTime)
		}

		n, writeErr := t.ioWriter.Write(t.payload)
		if writeErr != nil {
			statsHandler.AddTotalFailedSendPacketsNum(1)
			log.Println("Failed to send UDP data, ", writeErr.Error())
			return true
		}

		statsHandler.AddTotalSentPacketsNum(1)
		statsHandler.AddTotalSentPacketsBytes(uint64(n))
	}

	return true
}
