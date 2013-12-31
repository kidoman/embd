package main

import (
	"fmt"
	"time"

	"github.com/kid0m4n/go-rpi/sensor/keypad/matrix4x3"
)

func main() {
	rowPins := []int{4, 17, 27, 22}
	colPins := []int{23, 24, 25}

	keypad := matrix4x3.New(rowPins, colPins)

	for {
		if key, err := keypad.PressedKey(); err == nil {
			fmt.Printf("Key Pressed = %v\n", key)
		}

		time.Sleep(500 * time.Millisecond)
	}
}
