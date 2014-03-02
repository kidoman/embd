// +build ignore

package main

import (
	"log"
	"time"

	"github.com/kidoman/embd"
	"github.com/kidoman/embd/sensor/bmp085"
)

func main() {
	if err := embd.InitI2C(); err != nil {
		panic(err)
	}
	defer embd.CloseI2C()

	bus := embd.NewI2CBus(1)

	baro := bmp085.New(bus)
	defer baro.Close()

	for {
		temp, err := baro.Temperature()
		if err != nil {
			log.Panic(err)
		}
		log.Printf("Temp is %v", temp)
		pressure, err := baro.Pressure()
		if err != nil {
			log.Panic(err)
		}
		log.Printf("Pressure is %v", pressure)
		altitude, err := baro.Altitude()
		if err != nil {
			log.Panic(err)
		}
		log.Printf("Altitude is %v", altitude)

		time.Sleep(500 * time.Millisecond)
	}
}
