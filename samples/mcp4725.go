// +build ignore

package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"

	"github.com/kidoman/embd"
	"github.com/kidoman/embd/controller/mcp4725"
)

func main() {
	if err := embd.InitI2C(); err != nil {
		panic(err)
	}
	defer embd.CloseI2C()

	bus := embd.NewI2CBus(1)

	dac := mcp4725.New(bus, 0x62)
	defer dac.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)

	for {
		select {
		case <-stop:
			return
		default:
			voltage := rand.Intn(4096)
			if err := dac.SetVoltage(voltage); err != nil {
				fmt.Printf("mcp4725: %v\n", err)
			}
		}
	}
}
