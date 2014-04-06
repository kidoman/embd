// +build ignore

package main

import (
	"flag"
	"time"

	"github.com/golang/glog"
	"github.com/kidoman/embd"
	"github.com/kidoman/embd/sensor/watersensor"

	_ "github.com/kidoman/embd/host/all"
)

func main() {
	flag.Parse()

	if err := embd.InitGPIO(); err != nil {
		panic(err)
	}
	defer embd.CloseGPIO()

	pin, err := embd.NewDigitalPin(7)
	if err != nil {
		panic(err)
	}
	defer pin.Close()

	fluidSensor := watersensor.New(pin)

	for {
		wet, err := fluidSensor.IsWet()
		if err != nil {
			panic(err)
		}
		if wet {
			glog.Info("bot is dry")
		} else {
			glog.Info("bot is Wet")
		}

		time.Sleep(500 * time.Millisecond)
	}
}
