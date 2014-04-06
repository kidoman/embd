// +build ignore

package main

import (
	"flag"
	"time"

	"github.com/kidoman/embd"

	_ "github.com/kidoman/embd/host/all"
)

func main() {
	flag.Parse()

	if err := embd.InitGPIO(); err != nil {
		panic(err)
	}
	defer embd.CloseGPIO()

	led, err := embd.NewDigitalPin(10)
	if err != nil {
		panic(err)
	}
	defer led.Close()

	if err := led.SetDirection(embd.Out); err != nil {
		panic(err)
	}
	if err := led.Write(embd.High); err != nil {
		panic(err)
	}

	time.Sleep(1 * time.Second)

	if err := led.SetDirection(embd.In); err != nil {
		panic(err)
	}
}
