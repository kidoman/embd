package main

import (
	"time"

	"github.com/kidoman/embd"
)

func main() {
	h, _, err := embd.DetectHost()
	if err != nil {
		return
	}

	var pinNo interface{}

	switch h {
	case embd.HostBBB:
		pinNo = "P9_31"
	case embd.HostRPi:
		pinNo = 10
	default:
		panic("host not supported (yet :P)")
	}

	gpio, err := embd.NewGPIO()
	if err != nil {
		panic(err)
	}
	defer gpio.Close()

	led, err := gpio.DigitalPin(pinNo)
	if err != nil {
		panic(err)
	}
	defer led.Close()

	if err := led.SetDir(embd.Out); err != nil {
		panic(err)
	}
	if err := led.Write(embd.High); err != nil {
		panic(err)
	}

	time.Sleep(1 * time.Second)

	if err := led.SetDir(embd.In); err != nil {
		panic(err)
	}
}
