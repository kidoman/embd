// Generic LED driver.

package embd

import (
	"errors"
	"fmt"
	"strconv"
)

// LEDMap type represents a LED mapping for a host.
type LEDMap map[string][]string

type ledFactory func(string) LED

type ledDriver struct {
	ledMap LEDMap

	lf ledFactory

	initializedLEDs map[string]LED
}

// NewLEDDriver returns a LEDDriver interface which allows control
// over the LED subsystem.
func NewLEDDriver(ledMap LEDMap, lf ledFactory) LEDDriver {
	return &ledDriver{
		ledMap: ledMap,
		lf:     lf,

		initializedLEDs: map[string]LED{},
	}
}

func (d *ledDriver) lookup(k interface{}) (string, error) {
	var ks string
	switch key := k.(type) {
	case int:
		ks = strconv.Itoa(key)
	case string:
		ks = key
	case fmt.Stringer:
		ks = key.String()
	default:
		return "", errors.New("led: invalid key type")
	}

	for id := range d.ledMap {
		for _, alias := range d.ledMap[id] {
			if alias == ks {
				return id, nil
			}
		}
	}

	return "", fmt.Errorf("led: no match found for %q", k)
}

func (d *ledDriver) LED(k interface{}) (LED, error) {
	id, err := d.lookup(k)
	if err != nil {
		return nil, err
	}

	led := d.lf(id)
	d.initializedLEDs[id] = led

	return led, nil
}

func (d *ledDriver) Close() error {
	for _, led := range d.initializedLEDs {
		if err := led.Close(); err != nil {
			return err
		}
	}

	return nil
}
