// +build ignore

package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/all"
)

func main() {
	flag.Parse()
	if err := embd.InitSPI(); err != nil {
		panic(err)
	}

	bus := embd.NewSPIBus(embd.SpiMode0, 0, 1000000, 8, 0)
	defer clean(bus)

	for i := 0; i < 30; i++ {
		time.Sleep(1 * time.Second)
		val, _ := getSensorValue(bus)
		fmt.Printf("value is: %v\n", val)
	}

}

func clean(bus embd.SPIBus) {
	bus.Close()
	embd.CloseSPI()
}

func getSensorValue(bus embd.SPIBus) (uint16, error) {
	data := make([]uint8, 3)
	data[0] = 1
	data[1] = 128
	data[2] = 0
	var err error
	err = bus.TransferAndRecieveData(data)
	if err != nil {
		return uint16(0), err
	}
	return uint16(data[1]&0x03)<<8 | uint16(data[2]), nil
}
