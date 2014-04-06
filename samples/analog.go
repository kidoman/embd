// +build ignore

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/kidoman/embd"

	_ "github.com/kidoman/embd/host/all"
)

func main() {
	flag.Parse()

	if err := embd.InitGPIO(); err != nil {
		panic(err)
	}
	defer embd.CloseGPIO()

	pin, err := embd.NewAnalogPin(0)
	if err != nil {
		panic(err)
	}
	defer pin.Close()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	defer signal.Stop(quit)

	for {
		select {
		case <-time.After(100 * time.Millisecond):
			val, err := pin.Read()
			if err != nil {
				panic(err)
			}
			fmt.Printf("reading: %v\n", val)
		case <-quit:
			return
		}
	}
}
