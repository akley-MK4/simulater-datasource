package sender

import (
	"errors"
	"fmt"
	"log"
	"time"
)

func NewMgr(maxSenderNum int, network string, interval time.Duration, args ...interface{}) (*Mgr, error) {
	newSenderFunc := GetNewSenderFunc(network)
	if newSenderFunc == nil {
		return nil, fmt.Errorf("invalid network type %v", network)
	}

	if maxSenderNum <= 0 || interval <= 0 {
		return nil, errors.New("invalid parameter")
	}

	mgr := &Mgr{
		avgStats:     &Stats{},
		maxAvgStats:  &Stats{},
		statsHandler: NewStatsHandler(),
		rateUnit:     interval,
	}
	for i := 0; i < maxSenderNum; i++ {
		s, newErr := newSenderFunc(mgr, interval, args...)
		if newErr != nil {
			return nil, newErr
		}
		mgr.senderList = append(mgr.senderList, s)
	}

	return mgr, nil
}

type Mgr struct {
	rateUnit     time.Duration
	senderList   []ISender
	statsHandler *StatsHandler
	avgStats     *Stats
	maxAvgStats  *Stats
}

func (t *Mgr) Start() error {
	for _, s := range t.senderList {
		if err := s.Start(); err != nil {
			return err
		}
	}

	go t.countRateStats()
	go t.outputStats()
	return nil
}

func (t *Mgr) GetStatsHandler() *StatsHandler {
	return t.statsHandler
}

func (t *Mgr) countRateStats() {
	statsHandler := t.GetStatsHandler()
	//rateStatsHandler := NewStatsHandler()
	rateStats := &Stats{}
	for {
		time.Sleep(time.Second)
		stats := statsHandler.GetStats()

		sentPacketsNum := stats.TotalSentPacketsNum - rateStats.TotalSentPacketsNum
		sentPacketsBytes := stats.TotalSentPacketsBytes - rateStats.TotalSentPacketsBytes
		failedSendPacketsNum := stats.TotalFailedSendPacketsNum - rateStats.TotalFailedSendPacketsNum

		rateStats.TotalSentPacketsNum += sentPacketsNum
		rateStats.TotalSentPacketsBytes += sentPacketsBytes
		rateStats.TotalFailedSendPacketsNum += failedSendPacketsNum

		t.avgStats.TotalSentPacketsNum = sentPacketsNum
		t.avgStats.TotalSentPacketsBytes = sentPacketsBytes
		t.avgStats.TotalFailedSendPacketsNum = failedSendPacketsNum

		if t.avgStats.TotalSentPacketsNum > t.maxAvgStats.TotalSentPacketsNum {
			t.maxAvgStats.TotalSentPacketsNum = t.avgStats.TotalSentPacketsNum
		}
		if t.avgStats.TotalSentPacketsBytes > t.maxAvgStats.TotalSentPacketsBytes {
			t.maxAvgStats.TotalSentPacketsBytes = t.avgStats.TotalSentPacketsBytes
		}
		if t.avgStats.TotalSentPacketsNum > t.maxAvgStats.TotalSentPacketsNum {
			t.maxAvgStats.TotalSentPacketsNum = t.avgStats.TotalSentPacketsNum
		}
	}
}

func (t *Mgr) outputStats() {
	for {
		time.Sleep(time.Second * 5)
		log.Println("==========Total Stats / Second==========")
		totalStats := t.statsHandler.GetStats()
		fmt.Printf("SentPacketsNum: %d, SentPacketsBytes: %d, FailedSendPacketsNum: %d\n",
			totalStats.TotalSentPacketsNum, totalStats.TotalSentPacketsBytes, totalStats.TotalFailedSendPacketsNum)
		log.Print("=======================================\n\n")

		log.Println("==========Rate Stats / Second==========")
		fmt.Printf("SentPacketsNum: %d, SentPacketsBytes: %d, FailedSendPacketsNum: %d\n",
			t.avgStats.TotalSentPacketsNum, t.avgStats.TotalSentPacketsBytes, t.avgStats.TotalFailedSendPacketsNum)
		log.Print("=======================================\n\n")

		log.Println("==========Max Rate Stats / Second==========")
		fmt.Printf("SentPacketsNum: %d, SentPacketsBytes: %d, FailedSendPacketsNum: %d\n",
			t.maxAvgStats.TotalSentPacketsNum, t.maxAvgStats.TotalSentPacketsBytes, t.maxAvgStats.TotalFailedSendPacketsNum)
		log.Print("=======================================\n\n")

	}
}
