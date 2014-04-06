// +build ignore

package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/kidoman/embd"
	"github.com/kidoman/embd/sensor/bh1750fvi"

	_ "github.com/kidoman/embd/host/all"
)

func main() {
	flag.Parse()

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
		fmt.Printf("Lighting is %v lx\n", lighting)

		time.Sleep(500 * time.Millisecond)
	}
}
