package gpio

type Direction int

const (
	In Direction = iota
	Out
)

const (
	Low int = iota
	High
)

type DigitalPin interface {
	Write(val int) error
	Read() (int, error)

	SetDir(dir Direction) error
	ActiveLow(b bool) error

	Close() error
}

type GPIO interface {
	DigitalPin(key interface{}) (DigitalPin, error)

	Close() error
}
