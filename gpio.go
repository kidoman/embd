package embd

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

	PullUp() error
	PullDown() error

	Close() error
}

type GPIO interface {
	DigitalPin(key interface{}) (DigitalPin, error)

	Close() error
}

var gpioInstance GPIO

func InitGPIO() error {
	desc, err := DescribeHost()
	if err != nil {
		return err
	}

	gpioInstance = desc.GPIO()

	return nil
}

func CloseGPIO() error {
	return gpioInstance.Close()
}

func NewDigitalPin(key interface{}) (DigitalPin, error) {
	return gpioInstance.DigitalPin(key)
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
