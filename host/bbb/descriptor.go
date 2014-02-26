package bbb

import (
	"github.com/kidoman/embd/gpio"
	lgpio "github.com/kidoman/embd/host/generic/linux/gpio"
	li2c "github.com/kidoman/embd/host/generic/linux/i2c"
	"github.com/kidoman/embd/i2c"
)

type descriptor struct {
}

func (d *descriptor) GPIO() gpio.GPIO {
	return lgpio.New(pins)
}

func (d *descriptor) I2C() i2c.I2C {
	return li2c.New()
}

func Descriptor(rev int) *descriptor {
	return &descriptor{}
}
