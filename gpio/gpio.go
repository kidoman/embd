package gpio

import "github.com/kidoman/embd/host"

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

	SetDirection(dir Direction) error
	ActiveLow(b bool) error

	Close() error
}

type gpio interface {
	DigitalPin(key interface{}) (DigitalPin, error)

	Close() error
}

var instance gpio

func Open() error {
	desc, err := host.Describe()
	if err != nil {
		return err
	}

	instance = desc.GPIO().(gpio)

	return nil
}

func Close() error {
	return instance.Close()
}

func NewDigitalPin(key interface{}) (DigitalPin, error) {
	return instance.DigitalPin(key)
}

func DigitalWrite(key interface{}, val int) error {
	pin, err := NewDigitalPin(key)
	if err != nil {
		return err
	}

	return pin.Write(val)
}

func DigitalRead(key interface{}) (int, error) {
	pin, err := NewDigitalPin(key)
	if err != nil {
		return 0, err
	}

	return pin.Read()
}

func SetDirection(key interface{}, dir Direction) error {
	pin, err := NewDigitalPin(key)
	if err != nil {
		return err
	}

	return pin.SetDirection(dir)
}

func ActiveLow(key interface{}, b bool) error {
	pin, err := NewDigitalPin(key)
	if err != nil {
		return err
	}

	return pin.ActiveLow(b)
}
