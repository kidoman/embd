// +build ignore

package main

import (
	"fmt"

	"github.com/kidoman/embd"
)

func main() {
	embd.InitGPIO()
	defer embd.CloseGPIO()

	val, _ := embd.AnalogRead(0)
	fmt.Printf("Reading: %v\n", val)
}
