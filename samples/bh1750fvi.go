// +build ignore

package main

import (
	"log"
	"time"

	"github.com/kidoman/embd"
	"github.com/kidoman/embd/sensor/bh1750fvi"
)

func main() {
	if err := embd.InitI2C(); err != nil {
		panic(err)
	}
	defer embd.CloseI2C()

	bus := embd.NewI2CBus(1)

	sensor := bh1750fvi.New(bh1750fvi.High, bus)
	defer sensor.Close()

	for {
		lighting, err := sensor.Lighting()
		if err != nil {
			panic(err)
		}
		log.Printf("Lighting is %v lx", lighting)

		time.Sleep(500 * time.Millisecond)
	}
}
