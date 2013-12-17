package main

import (
	"log"
	"time"

	"github.com/kid0m4n/go-rpi/i2c"
	"github.com/kid0m4n/go-rpi/sensor/bh1750Fvi"
)

func main() {
	bus, err := i2c.NewBus(1)
	if err != nil {
		log.Panic(err)
	}

	lightingSensor := bh1750Fvi.New("H", bus)

	defer lightingSensor.Close()

	for {
		lighting, err := lightingSensor.Lighting()
		if err != nil {
			log.Panic(err)
		}
		log.Printf("Lighting is %v", lighting, "lx")
		time.Sleep(500 * time.Millisecond)
	}
}
