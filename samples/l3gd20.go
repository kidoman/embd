package main

import (
	"log"
	"time"

	"github.com/kid0m4n/go-rpi/i2c"
	"github.com/kid0m4n/go-rpi/sensor/l3gd20"
)

func main() {
	bus, err := i2c.NewBus(1)
	if err != nil {
		log.Panic(err)
	}
	gyro := l3gd20.New(bus, l3gd20.R250DPS)
	defer gyro.Close()

	x, y, z := 0.0, 0.0, 0.0
	dt := 0.02

	for {
		dx, dy, dz, err := gyro.Orientation()
		if err != nil {
			log.Panic(err)
		}

		x += dx * dt
		y += dy * dt
		z += dz * dt

		log.Printf("Orientation is (%v, %v, %v)", x, y, z)

		temp, err := gyro.Temperature()
		if err != nil {
			log.Panic(err)
		}

		log.Printf("Temperature is %v", temp)

		time.Sleep(time.Duration(dt*1000) * time.Millisecond)
	}
}
