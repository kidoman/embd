// LED support.

package embd

// The LED interface is used to control a led on the prototyping board.
type LED interface {
	// On switches the LED on.
	On() error

	// Off switches the LED off.
	Off() error

	// Toggle toggles the LED.
	Toggle() error

	// Close releases resources associated with the LED.
	Close() error
}

// LEDDriver interface interacts with the host descriptors to allow us
// control of the LEDs.
type LEDDriver interface {
	LED(key interface{}) (LED, error)

	Close() error
}

var ledDriverInitialized bool
var ledDriverInstance LEDDriver

// InitLED initializes the LED driver.
func InitLED() error {
	if ledDriverInitialized {
		return nil
	}

	desc, err := DescribeHost()
	if err != nil {
		return err
	}

	if desc.LEDDriver == nil {
		return ErrFeatureNotSupported
	}

	ledDriverInstance = desc.LEDDriver()
	ledDriverInitialized = true

	return nil
}

// CloseLED releases resources associated with the LED driver.
func CloseLED() error {
	return ledDriverInstance.Close()
}

// NewLED returns a LED interface which allows control over the LED.
func NewLED(key interface{}) (LED, error) {
	if err := InitLED(); err != nil {
		return nil, err
	}

	return ledDriverInstance.LED(key)
}

// LEDOn switches the LED on.
func LEDOn(key interface{}) error {
	led, err := NewLED(key)
	if err != nil {
		return err
	}

	return led.On()
}

// LEDOff switches the LED off.
func LEDOff(key interface{}) error {
	led, err := NewLED(key)
	if err != nil {
		return err
	}

	return led.Off()
}

// LEDToggle toggles the LED.
func LEDToggle(key interface{}) error {
	led, err := NewLED(key)
	if err != nil {
		return err
	}

	return led.Toggle()
}
