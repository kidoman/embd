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

	if err := gpio.SetDirection(10, gpio.Out); err != nil {
		panic(err)
	}
	if err := gpio.DigitalWrite(10, gpio.High); err != nil {
		panic(err)
	}

	time.Sleep(1 * time.Second)

	if err := gpio.SetDirection(10, gpio.In); err != nil {
		panic(err)
	}
}
