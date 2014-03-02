// +build ignore

package main

import "github.com/kidoman/embd/gpio"

func main() {
	gpio.Open()
	defer gpio.Close()

	gpio.SetDirection(10, gpio.Out)
	gpio.DigitalWrite(10, gpio.High)
}
