package types

type ChanDir int

const (
	SendRecv ChanDir = iota
	SendOnly
	RecvOnly
)

type Chan struct {
	Type any
	Dir  ChanDir
}
