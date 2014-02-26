package main

import (
	"log"
	"time"
	"github.com/kidoman/embd"
	"github.com/kidoman/embd/sensor/bh1750fvi"
)

func main() {
	i2c, err := embd.NewI2C()
	if err != nil {
		panic(err)
	}

	bus := i2c.Bus(1)

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
