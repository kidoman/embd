// Pin mapping support.

package embd

import (
	"fmt"
	"strconv"
)

const (
	// CapDigital represents the digital IO capability.
	CapDigital int = 1 << iota

	// CapI2C represents pins with the I2C capability.
	CapI2C

	// CapUART represents pins with the UART capability.
	CapUART

	// CapSPI represents pins with the SPI capability.
	CapSPI

	// CapGPMS represents pins with the GPMC capability.
	CapGPMC

	// CapLCD represents pins used to carry LCD data.
	CapLCD

	// CapPWM represents pins with PWM capability.
	CapPWM

	// CapAnalog represents pins with analog IO capability.
	CapAnalog
)

// PinDesc represents a pin descriptor.
type PinDesc struct {
	ID      string
	Aliases []string
	Caps    int

	DigitalLogical int
	AnalogLogical  int
}

// PinMap type represents a collection of pin descriptors.
type PinMap []*PinDesc

// Lookup returns a pin descriptor matching the provided key and capability
// combination. This allows the same keys to be used across pins with differing
// capabilities. For example, it is perfectly fine to have:
//
//	pin1: {Aliases: [10, GPIO10], Cap: CapDigital}
//	pin2: {Aliases: [10, AIN0], Cap: CapAnalog}
//
// Searching for 10 with CapDigital will return pin1 and searching for
// 10 with CapAnalog will return pin2. This makes for a very pleasant to use API.
func (m PinMap) Lookup(k interface{}, cap int) (*PinDesc, bool) {
	var ks string
	switch key := k.(type) {
	case int:
		ks = strconv.Itoa(key)
	case string:
		ks = key
	case fmt.Stringer:
		ks = key.String()
	default:
		return nil, false
	}

	for i := range m {
		pd := m[i]

		if pd.ID == ks {
			return pd, true
		}

		for j := range pd.Aliases {
			if pd.Aliases[j] == ks && pd.Caps&cap != 0 {
				return pd, true
			}
		}
	}

	return nil, false
}
