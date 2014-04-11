// +build ignore

// Simple LED blinker, works OOTB on a RPi. However, it does not clean up
// after itself. So might leave the LED On. The RPi is not harmed though.

package main

import (
	"time"

	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi"
)

func main() {
	for {
		embd.LEDToggle(0)
		time.Sleep(250 * time.Millisecond)
	}
}
