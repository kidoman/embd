package main

import (
	"log"
	"time"

	"github.com/kid0m4n/go-rpi/i2c"
	"github.com/kid0m4n/go-rpi/sensor/lsm303"
)

func main() {
	bus := i2c.NewBus(1)

	mems := lsm303.New(bus)
	defer mems.Close()

	for {
		heading, err := mems.Heading()
		if err != nil {
			log.Panic(err)
		}
		log.Printf("Heading is %v", heading)

		time.Sleep(500 * time.Millisecond)
	}
}
