package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"

	"github.com/kid0m4n/go-rpi/controller/mcp4725"
	"github.com/kid0m4n/go-rpi/i2c"
)

func main() {
	bus, err := i2c.NewBus(1)
	if err != nil {
		log.Panic(err)
	}

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
				log.Printf("mcp4725: %v", err)
			}
		}
	}
}
