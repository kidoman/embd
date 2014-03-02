// +build ignore

package main

import "github.com/kidoman/embd"

func main() {
	embd.InitGPIO()
	defer embd.CloseGPIO()

	embd.SetDirection(10, embd.Out)
	embd.DigitalWrite(10, embd.High)
}
