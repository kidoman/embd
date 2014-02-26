package embd

import (
	"github.com/kidoman/embd/gpio"
	"github.com/kidoman/embd/i2c"
)

const (
	In  = gpio.In
	Out = gpio.Out
)

const (
	Low  = gpio.Low
	High = gpio.High
)

func NewGPIO() (gpio.GPIO, error) {
	desc, err := describeHost()
	if err != nil {
		return nil, err
	}

	return desc.GPIO(), nil
}

func NewI2C() (i2c.I2C, error) {
	desc, err := describeHost()
	if err != nil {
		return nil, err
	}

	return desc.I2C(), nil
}
