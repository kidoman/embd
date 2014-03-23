// LED support.

package embd

// The LED interface is used to control a led on the prototyping board.
type LED interface {
	On() error
	Off() error

	Toggle() error

	Close() error
}

// LEDDriver interface interacts with the host descriptors to allow us
// control of the LEDs.
type LEDDriver interface {
	LED(key interface{}) (LED, error)

	Close() error
}

var ledDriverInstance LEDDriver

// InitLED initializes the LED driver.
func InitLED() error {
	desc, err := DescribeHost()
	if err != nil {
		return err
	}

	if desc.LEDDriver == nil {
		return ErrFeatureNotSupported
	}

	ledDriverInstance = desc.LEDDriver()

	return nil
}

// CloseLED gracefully closes the LED driver.
func CloseLED() error {
	return ledDriverInstance.Close()
}

// NewLED returns a LED interface which allows control over the LED
// represented by key.
func NewLED(key interface{}) (LED, error) {
	return ledDriverInstance.LED(key)
}

// LEDOn switches the corresponding LED on.
func LEDOn(key interface{}) error {
	led, err := NewLED(key)
	if err != nil {
		return err
	}

	return led.On()
}

// LEDOff switches the corresponding LED off.
func LEDOff(key interface{}) error {
	led, err := NewLED(key)
	if err != nil {
		return err
	}

	return led.Off()
}

// LEDToggle toggles the corresponding LED.
func LEDToggle(key interface{}) error {
	led, err := NewLED(key)
	if err != nil {
		return err
	}

	return led.Toggle()
}
