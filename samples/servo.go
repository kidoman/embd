package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/kid0m4n/go-rpi/controller/pca9685"
	"github.com/kid0m4n/go-rpi/i2c"
	"github.com/kid0m4n/go-rpi/motion/servo"
)

func main() {
	bus, err := i2c.NewBus(1)
	if err != nil {
		log.Panic(err)
	}

	pwm := pca9685.New(bus, 0x42, 50)
	defer pwm.Close()

	servo := servo.New(pwm, 50, 0, 1, 2.5)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	turnTimer := time.Tick(500 * time.Millisecond)
	left := true

	servo.SetAngle(90)
	defer func() {
		servo.SetAngle(90)
	}()

	for {
		select {
		case <-turnTimer:
			left = !left
			switch left {
			case true:
				servo.SetAngle(70)
			case false:
				servo.SetAngle(110)
			}
		case <-c:
			return
		}
	}
}
