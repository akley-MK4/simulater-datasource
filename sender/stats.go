package sender

import "sync/atomic"

type Stats struct {
	TotalSentPacketsNum       uint64
	TotalSentPacketsBytes     uint64
	TotalFailedSendPacketsNum uint64
}

func NewStatsHandler() *StatsHandler {
	return &StatsHandler{}
}

type StatsHandler struct {
	stats Stats
}

func (t *StatsHandler) GetStats() Stats {
	return t.stats
}

func (t *StatsHandler) AddTotalSentPacketsNum(delta uint64) {
	atomic.AddUint64(&t.stats.TotalSentPacketsNum, delta)
}

func (t *StatsHandler) AddTotalSentPacketsBytes(delta uint64) {
	atomic.AddUint64(&t.stats.TotalSentPacketsBytes, delta)
}

func (t *StatsHandler) AddTotalFailedSendPacketsNum(delta uint64) {
	atomic.AddUint64(&t.stats.TotalFailedSendPacketsNum, delta)
}
