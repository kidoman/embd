package host

import (
	"errors"

	"github.com/kidoman/embd/gpio"
	"github.com/kidoman/embd/host/rpi"
	"github.com/kidoman/embd/i2c"
)

type Descriptor interface {
	GPIO() gpio.GPIO
	I2C() i2c.I2C
}

func Describe() (Descriptor, error) {
	host, rev, err := Detect()
	if err != nil {
		return nil, err
	}

	switch host {
	case RPi:
		return rpi.Descriptor(rev), nil
	default:
		return nil, errors.New("host: invalid host")
	}
}
