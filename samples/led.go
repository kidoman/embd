// +build ignore

// LED example, works OOTB on a BBB.

package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/kidoman/embd"
)

func main() {
	if err := embd.InitLED(); err != nil {
		panic(err)
	}
	defer embd.CloseLED()

	led, err := embd.NewLED(3)
	if err != nil {
		panic(err)
	}
	defer func() {
		led.Off()
		led.Close()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	defer signal.Stop(quit)

	for {
		select {
		case <-time.After(500 * time.Millisecond):
			if err := led.Toggle(); err != nil {
				panic(err)
			}
			fmt.Printf("Toggled\n")
		case <-quit:
			return
		}
	}
}
