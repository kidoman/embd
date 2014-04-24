// +build ignore

// LED example, works OOTB on a BBB.

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/kidoman/embd"

	_ "github.com/kidoman/embd/host/bbb"
)

func main() {
	flag.Parse()

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
		case <-time.After(250 * time.Millisecond):
			if err := led.Toggle(); err != nil {
				panic(err)
			}
			fmt.Println("Toggled")
		case <-quit:
			return
		}
	}
}
