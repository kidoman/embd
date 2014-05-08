// +build ignore

package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/kidoman/embd/convertors/mcp3008"

	"github.com/kidoman/embd"
)

func main() {
	flag.Parse()
	fmt.Println("this is a sample code for mcp3008 10bit 8 channel ADC")

	if err := embd.InitSPI(); err != nil {
		panic(err)
	}
	defer embd.CloseSPI()

	spiBus := embd.NewSPIBus(embd.SpiMode0, 0, 1000000, 8, 0)
	defer spiBus.Close()

	adc := mcp3008.New(mcp3008.SingleMode, spiBus)

	for i := 0; i < 20; i++ {
		time.Sleep(1 * time.Second)
		val, err := adc.AnalogValueAt(0)
		if err != nil {
			panic(err)
		}

		fmt.Printf("analog value is: %v\n", val)
	}

}
