// GPIO support.

package embd

import "time"

// The Direction type indicates the direction of a GPIO pin.
type Direction int

// The Edge trigger for the GPIO Interrupt
type Edge string

const (
	// In represents read mode.
	In Direction = iota

	// Out represents write mode.
	Out
)

const (
	// Low represents 0.
	Low int = iota

	// High represents 1.
	High
)

const (
	EdgeNone    Edge = "none"
	EdgeRising  Edge = "rising"
	EdgeFalling Edge = "falling"
	EdgeBoth    Edge = "both"
)

// InterruptPin implements access to an interrupt capable GPIO pin.
// The basic capability provided is to watch for a transition on the pin and
// generate a callback to a handler when a transition occurs.
// On Linux the underlying implementation generally uses epoll to receive the
// interrupts at user-level.
type InterruptPin interface {

	// Start watching this pin for interrupt
	Watch(edge Edge, handler func(DigitalPin)) error

	// Stop watching this pin for interrupt
	StopWatching() error
}

// DigitalPin implements access to a digital IO capable GPIO pin.
type DigitalPin interface {
	InterruptPin

	// N returns the logical GPIO number.
	N() int

	// Write writes the provided value to the pin.
	Write(val int) error

	// Read reads the value from the pin.
	Read() (int, error)

	// TimePulse measures the duration of a pulse on the pin.
	TimePulse(state int) (time.Duration, error)

	// SetDirection sets the direction of the pin (in/out).
	SetDirection(dir Direction) error

	// ActiveLow makes the pin active low. A low logical state is represented by
	// a high state on the physical pin, and vice-versa.
	ActiveLow(b bool) error

	// PullUp pulls the pin up.
	PullUp() error

	// PullDown pulls the pin down.
	PullDown() error

	// Close releases the resources associated with the pin.
	Close() error
}

// AnalogPin implements access to a analog IO capable GPIO pin.
type AnalogPin interface {
	// N returns the logical GPIO number.
	N() int

	// Read reads the value from the pin.
	Read() (int, error)

	// Close releases the resources associated with the pin.
	Close() error
}

// The Polarity type indicates the polarity of a pwm pin.
type Polarity int

const (
	// Positive represents (default) positive polarity.
	Positive Polarity = iota

	// Negative represents negative polarity.
	Negative
)

// PWMPin implements access to a pwm capable GPIO pin.
type PWMPin interface {
	// N returns the logical PWM id.
	N() string

	// SetPeriod sets the period of a pwm pin.
	SetPeriod(ns int) error

	// SetDuty sets the duty of a pwm pin.
	SetDuty(ns int) error

	// SetPolarity sets the polarity of a pwm pin.
	SetPolarity(pol Polarity) error

	// SetMicroseconds sends a command to the PWM driver to generate a us wide pulse.
	SetMicroseconds(us int) error

	// SetAnalog allows easy manipulation of the PWM based on a (0-255) range value.
	SetAnalog(value byte) error

	// Close releases the resources associated with the pin.
	Close() error
}

// GPIODriver implements a generic GPIO driver.
type GPIODriver interface {
	// PinMap returns the pinmap for this driver.
	PinMap() PinMap

	// Unregister unregisters the pin from the driver. Should be called when the pin is closed.
	Unregister(string) error

	// DigitalPin returns a pin capable of doing digital IO.
	DigitalPin(key interface{}) (DigitalPin, error)

	// AnalogPin returns a pin capable of doing analog IO.
	AnalogPin(key interface{}) (AnalogPin, error)

	// PWMPin returns a pin capable of generating PWM.
	PWMPin(key interface{}) (PWMPin, error)

	// Close releases the resources associated with the driver.
	Close() error
}

var gpioDriverInitialized bool
var gpioDriverInstance GPIODriver

// InitGPIO initializes the GPIO driver.
func InitGPIO() error {
	if gpioDriverInitialized {
		return nil
	}

	desc, err := DescribeHost()
	if err != nil {
		return err
	}

	if desc.GPIODriver == nil {
		return ErrFeatureNotSupported
	}

	gpioDriverInstance = desc.GPIODriver()
	gpioDriverInitialized = true

	return nil
}

// CloseGPIO releases resources associated with the GPIO driver.
func CloseGPIO() error {
	return gpioDriverInstance.Close()
}

// NewDigitalPin returns a DigitalPin interface which allows control over
// the digital GPIO pin.
func NewDigitalPin(key interface{}) (DigitalPin, error) {
	if err := InitGPIO(); err != nil {
		return nil, err
	}

	return gpioDriverInstance.DigitalPin(key)
}

// DigitalWrite writes val to the pin.
func DigitalWrite(key interface{}, val int) error {
	pin, err := NewDigitalPin(key)
	if err != nil {
		return err
	}

	return pin.Write(val)
}

// DigitalRead reads a value from the pin.
func DigitalRead(key interface{}) (int, error) {
	pin, err := NewDigitalPin(key)
	if err != nil {
		return 0, err
	}

	return pin.Read()
}

// SetDirection sets the direction of the pin (in/out).
func SetDirection(key interface{}, dir Direction) error {
	pin, err := NewDigitalPin(key)
	if err != nil {
		return err
	}

	return pin.SetDirection(dir)
}

// ActiveLow makes the pin active low. A low logical state is represented by
// a high state on the physical pin, and vice-versa.
func ActiveLow(key interface{}, b bool) error {
	pin, err := NewDigitalPin(key)
	if err != nil {
		return err
	}

	return pin.ActiveLow(b)
}

// PullUp pulls the pin up.
func PullUp(key interface{}) error {
	pin, err := NewDigitalPin(key)
	if err != nil {
		return err
	}

	return pin.PullUp()
}

// PullDown pulls the pin down.
func PullDown(key interface{}) error {
	pin, err := NewDigitalPin(key)
	if err != nil {
		return err
	}

	return pin.PullDown()
}

// NewAnalogPin returns a AnalogPin interface which allows control over
// the analog GPIO pin.
func NewAnalogPin(key interface{}) (AnalogPin, error) {
	if err := InitGPIO(); err != nil {
		return nil, err
	}

	return gpioDriverInstance.AnalogPin(key)
}

// AnalogWrite reads a value from the pin.
func AnalogRead(key interface{}) (int, error) {
	pin, err := NewAnalogPin(key)
	if err != nil {
		return 0, err
	}

	return pin.Read()
}

// NewPWMPin returns a PWMPin interface which allows PWM signal
// generation over a the PWM pin.
func NewPWMPin(key interface{}) (PWMPin, error) {
	if err := InitGPIO(); err != nil {
		return nil, err
	}

	return gpioDriverInstance.PWMPin(key)
}
