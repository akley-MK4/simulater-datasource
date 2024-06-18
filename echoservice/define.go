package echoservice

const (
	defaultReadBuffSize = 1024 * 1024
)

type IEchoService interface {
	Start() error
}
