package main

import (
	"log"
	"time"

	"github.com/kidoman/embd"
	"github.com/kidoman/embd/sensor/lsm303"
)

func main() {
	i2c, err := embd.NewI2C()
	if err != nil {
		panic(err)
	}
	defer i2c.Close()

	bus := i2c.Bus(1)

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
