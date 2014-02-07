package main

import (
	"log"
	"time"

	"github.com/kid0m4n/go-rpi/i2c"
	"github.com/kid0m4n/go-rpi/sensor/bmp085"
)

func main() {
	bus := i2c.NewBus(1)

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
