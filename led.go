package embd

type LED interface {
	On() error
	Off() error

	Toggle() error

	Close() error
}

type LEDDriver interface {
	LED(key interface{}) (LED, error)

	Close() error
}

var ledDriverInstance LEDDriver

func InitLED() error {
	desc, err := DescribeHost()
	if err != nil {
		return err
	}

	if desc.LEDDriver == nil {
		return ErrFeatureNotSupport
	}

	ledDriverInstance = desc.LEDDriver()

	return nil
}

func CloseLED() error {
	return ledDriverInstance.Close()
}

func NewLED(key interface{}) (LED, error) {
	return ledDriverInstance.LED(key)
}

func LEDOn(key interface{}) error {
	led, err := NewLED(key)
	if err != nil {
		return err
	}

	return led.On()
}

func LEDOff(key interface{}) error {
	led, err := NewLED(key)
	if err != nil {
		return err
	}

	return led.Off()
}
