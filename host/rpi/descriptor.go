package rpi

import (
	lgpio "github.com/kidoman/embd/driver/linux/gpio"
	li2c "github.com/kidoman/embd/driver/linux/i2c"
	"github.com/kidoman/embd/gpio"
	"github.com/kidoman/embd/i2c"
)

type descriptor struct {
	rev int
}

func (d *descriptor) GPIO() gpio.GPIO {
	var pins = rev1Pins
	if d.rev > 1 {
		pins = rev2Pins
	}

	return lgpio.New(pins)
}

func (d *descriptor) I2C() i2c.I2C {
	return li2c.New()
}

func Descriptor(rev int) *descriptor {
	return &descriptor{rev}
}
