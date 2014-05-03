// +build ignore

package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/kidoman/embd/convertors/mcp3008"
)

func main() {
	flag.Parse()
	fmt.Println("This is a sample code for mcp3008 10bit 8 channel ADC")

	adc, err := mcp3008.New(mcp3008.SingleMode, 0, 1000000)
	if err != nil {
		panic(err)
	}
	defer adc.Close()

	for i := 0; i < 20; i++ {
		time.Sleep(1 * time.Second)
		val, err := adc.AnalogValueAt(0)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Analog value is: %v\n", val)
	}

}
