package sender

import "time"

type ISender interface {
	Start() error
	Stop() error
}

type NewSenderFunc func(mgr *Mgr, interval time.Duration, args ...interface{}) (ISender, error)

var (
	registerNewSenderFunc = map[string]NewSenderFunc{
		"udp": NewUDPSender,
	}
)

func GetNewSenderFunc(network string) NewSenderFunc {
	return registerNewSenderFunc[network]
}
