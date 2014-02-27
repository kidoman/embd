package bbb

import (
	lgpio "github.com/kidoman/embd/driver/linux/gpio"
	li2c "github.com/kidoman/embd/driver/linux/i2c"
	"github.com/kidoman/embd/gpio"
	"github.com/kidoman/embd/host"
)

func init() {
	host.Describers[host.BBB] = describer
}

func describer(rev int) *host.Descriptor {
	return &host.Descriptor{
		GPIO: func() gpio.GPIO {
			return lgpio.New(pins)
		},
		I2C: li2c.New,
	}
}
