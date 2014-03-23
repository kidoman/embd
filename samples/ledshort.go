// +build ignore

// Short LED example, works OOTB on a BBB.

package main

import (
	"time"

	"github.com/kidoman/embd"
)

func main() {
	embd.InitLED()
	defer embd.CloseLED()

	embd.LEDOn(3)
	time.Sleep(1 * time.Second)
	embd.LEDOff(3)
}
