package embd

import (
	"errors"

	"github.com/kidoman/embd/gpio"
	"github.com/kidoman/embd/host/rpi"
	"github.com/kidoman/embd/i2c"
)

type descriptor interface {
	GPIO() gpio.GPIO
	I2C() i2c.I2C
}

func describeHost() (descriptor, error) {
	host, rev, err := DetectHost()
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
