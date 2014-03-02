package main

import (
	"fmt"
	"time"

	"github.com/kidoman/embd/gpio"
	"github.com/kidoman/embd/interface/keypad/matrix4x3"
)

func main() {
	rowPins := []int{4, 17, 27, 22}
	colPins := []int{23, 24, 25}

	if err := gpio.Open(); err != nil {
		panic(err)
	}
	defer gpio.Close()

	keypad, err := matrix4x3.New(rowPins, colPins)
	if err != nil {
		panic(err)
	}

	for {
		key, err := keypad.PressedKey()
		if err != nil {
			panic(err)
		}
		if key != matrix4x3.KNone {
			fmt.Printf("Key Pressed = %v\n", key)
		}

		time.Sleep(500 * time.Millisecond)
	}
}
