// +build ignore

package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/all"
)

const (
	channel = 0
	speed   = 1000000
	bpw     = 8
	delay   = 0
)

func main() {
	flag.Parse()
	if err := embd.InitSPI(); err != nil {
		panic(err)
	}
	defer embd.CloseSPI()

	bus := embd.NewSPIBus(embd.SPIMode0, channel, speed, bpw, delay)
	defer bus.Close()

	for i := 0; i < 30; i++ {
		time.Sleep(1 * time.Second)
		val, err := getSensorValue(bus)
		if err != nil {
			panic(err)
		}
		fmt.Printf("value is: %v\n", val)
	}
}

func getSensorValue(bus embd.SPIBus) (uint16, error) {
	data := [3]uint8{1, 128, 0}

	var err error
	err = bus.TransferAndReceiveData(data[:])
	if err != nil {
		return uint16(0), err
	}
	return uint16(data[1]&0x03)<<8 | uint16(data[2]), nil
}
