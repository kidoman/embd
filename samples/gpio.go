package main

import (
	"time"

	"github.com/kidoman/embd/gpio"
)

func main() {
	io := gpio.New()
	defer io.Close()

	pin, err := io.Pin("MOSI")
	if err != nil {
		panic(err)
	}

	if err := pin.Output(); err != nil {
		panic(err)
	}
	if err := pin.ActiveLow(); err != nil {
		panic(err)
	}
	if err := pin.Low(); err != nil {
		panic(err)
	}

	time.Sleep(1 * time.Second)

	if err := pin.Input(); err != nil {
		panic(err)
	}
}
