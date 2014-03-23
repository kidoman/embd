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
	N() int

	Write(val int) error
	Read() (int, error)

	SetDirection(dir Direction) error
	ActiveLow(b bool) error

	PullUp() error
	PullDown() error

	Close() error
}

type AnalogPin interface {
	N() int

	Write(val int) error
	Read() (int, error)

	Close() error
}

type GPIODriver interface {
	DigitalPin(key interface{}) (DigitalPin, error)
	AnalogPin(key interface{}) (AnalogPin, error)

	Close() error
}

var gpioDriverInstance GPIODriver

func InitGPIO() error {
	desc, err := DescribeHost()
	if err != nil {
		return err
	}

	if desc.GPIODriver == nil {
		return ErrFeatureNotSupport
	}

	gpioDriverInstance = desc.GPIODriver()

	return nil
}

func CloseGPIO() error {
	return gpioDriverInstance.Close()
}

func NewDigitalPin(key interface{}) (DigitalPin, error) {
	return gpioDriverInstance.DigitalPin(key)
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

func NewAnalogPin(key interface{}) (AnalogPin, error) {
	return gpioDriverInstance.AnalogPin(key)
}

func AnalogWrite(key interface{}, val int) error {
	pin, err := NewAnalogPin(key)
	if err != nil {
		return err
	}

	return pin.Write(val)
}

func AnalogRead(key interface{}) (int, error) {
	pin, err := NewAnalogPin(key)
	if err != nil {
		return 0, err
	}

	return pin.Read()
}
