// +build ignore

package main

import (
	"time"

	"github.com/golang/glog"
	"github.com/kidoman/embd/gpio"
	"github.com/kidoman/embd/sensor/watersensor"
)

func main() {
	if err := gpio.Open(); err != nil {
		panic(err)
	}
	defer gpio.Close()

	pin, err := gpio.NewDigitalPin(7)
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
