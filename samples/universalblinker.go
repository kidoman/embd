// +build ignore

// Universal LED blinker, works OOTB on a RPi / BBB.

package main

import (
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/all"
)

func main() {
	flag.Parse()

	panicIf(embd.InitLED())
	defer embd.CloseLED()

	led, err := embd.NewLED(ledToBlink())
	panicIf(err)
	defer led.Off()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	defer signal.Stop(quit)

	for {
		select {
		case <-time.After(200 * time.Millisecond):
			panicIf(led.Toggle())
		case <-quit:
			return
		}
	}
}

func ledToBlink() string {
	host, _, err := embd.DetectHost()
	panicIf(err)

	switch host {
	case embd.HostRPi:
		return "LED0"
	case embd.HostBBB:
		return "USR3"
	}

	panic("Unsupported host")
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
