// Generic GPIO driver.

package embd

import (
	"errors"
	"fmt"
)

type pin interface {
	Close() error
}

type gpioDriver struct {
	pinMap PinMap

	dpf func(n int) DigitalPin
	apf func(n int) AnalogPin

	initializedPins map[string]pin
}

func newGPIODriver(pinMap PinMap, dpf func(n int) DigitalPin, apf func(n int) AnalogPin) GPIODriver {
	return &gpioDriver{
		pinMap: pinMap,
		dpf:    dpf,
		apf:    apf,

		initializedPins: map[string]pin{},
	}
}

func (io *gpioDriver) DigitalPin(key interface{}) (DigitalPin, error) {
	if io.dpf == nil {
		return nil, errors.New("gpio: digital io not supported on this host")
	}

	pd, found := io.pinMap.Lookup(key, CapDigital)
	if !found {
		return nil, fmt.Errorf("gpio: could not find pin matching %v", key)
	}

	p := io.dpf(pd.DigitalLogical)
	io.initializedPins[pd.ID] = p

	return p, nil
}

func (io *gpioDriver) AnalogPin(key interface{}) (AnalogPin, error) {
	if io.apf == nil {
		return nil, errors.New("gpio: analog io not supported on this host")
	}

	pd, found := io.pinMap.Lookup(key, CapAnalog)
	if !found {
		return nil, fmt.Errorf("gpio: could not find pin matching %v", key)
	}

	p := io.apf(pd.AnalogLogical)
	io.initializedPins[pd.ID] = p

	return p, nil
}

func (io *gpioDriver) Close() error {
	for _, p := range io.initializedPins {
		if err := p.Close(); err != nil {
			return err
		}
	}

	return nil
}
