package main

import (
	"log"
	"time"

	"github.com/kidoman/embd/i2c"
	"github.com/kidoman/embd/sensor/bmp180"
)

func main() {
	if err := i2c.Open(); err != nil {
		panic(err)
	}
	defer i2c.Close()

	bus := i2c.NewBus(1)

	baro := bmp180.New(bus)
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
