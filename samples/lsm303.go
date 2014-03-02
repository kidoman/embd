// +build ignore

package main

import (
	"log"
	"time"

	"github.com/kidoman/embd/i2c"
	"github.com/kidoman/embd/sensor/lsm303"
)

func main() {
	if err := i2c.Open(); err != nil {
		panic(err)
	}
	defer i2c.Close()

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
