package host

import (
	"errors"

	"github.com/kidoman/embd/gpio"
	"github.com/kidoman/embd/i2c"
)

type Descriptor struct {
	GPIO func() gpio.GPIO
	I2C  func() i2c.I2C
}

type Describer func(rev int) *Descriptor

var Describers = map[Host]Describer{}

func Describe() (*Descriptor, error) {
	host, rev, err := Detect()
	if err != nil {
		return nil, err
	}

	describer, ok := Describers[host]
	if !ok {
		return nil, errors.New("host: invalid host")
	}

	return describer(rev), nil
}
