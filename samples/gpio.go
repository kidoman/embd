// +build ignore

package main

import (
	"time"

	"github.com/kidoman/embd/gpio"
)

func main() {
	if err := gpio.Open(); err != nil {
		panic(err)
	}
	defer gpio.Close()

	led, err := gpio.NewDigitalPin(10)
	if err != nil {
		panic(err)
	}
	defer led.Close()

	if err := led.SetDirection(gpio.Out); err != nil {
		panic(err)
	}
	if err := led.Write(gpio.High); err != nil {
		panic(err)
	}

	time.Sleep(1 * time.Second)

	if err := led.SetDirection(gpio.In); err != nil {
		panic(err)
	}
}
