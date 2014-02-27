package rpi

import (
	lgpio "github.com/kidoman/embd/driver/linux/gpio"
	li2c "github.com/kidoman/embd/driver/linux/i2c"
	"github.com/kidoman/embd/gpio"
	"github.com/kidoman/embd/i2c"
)

func init() {
	host.Describers[host.RPi] = describer
}

func describer(rev int) *host.Descriptor {
	var pins = rev1Pins
	if d.rev > 1 {
		pins = rev2Pins
	}

	return &host.Descriptor{
		GPIO: func() gpio.GPIO {
			return lgpio.New(pins)
		},
		I2C: li2c.New,
	}
}
