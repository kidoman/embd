// +build ignore

package main

import (
	"log"
	"time"

	"github.com/kidoman/embd/i2c"
	"github.com/kidoman/embd/sensor/bh1750fvi"
)

func main() {
	if err := i2c.Open(); err != nil {
		panic(err)
	}
	defer i2c.Close()

	bus := i2c.NewBus(1)

	sensor := bh1750fvi.New(bh1750fvi.High, bus)
	defer sensor.Close()

	for {
		lighting, err := sensor.Lighting()
		if err != nil {
			log.Panic(err)
		}
		log.Printf("Lighting is %v lx", lighting)

		time.Sleep(500 * time.Millisecond)
	}
}
