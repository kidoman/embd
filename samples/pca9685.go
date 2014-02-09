package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/kidoman/embd/controller/pca9685"
	"github.com/kidoman/embd/i2c"
)

func main() {
	bus := i2c.NewBus(1)

	pca9685 := pca9685.New(bus, 0x41, 1000)
	pca9685.Debug = true
	defer pca9685.Close()

	if err := pca9685.SetPwm(15, 0, 2000); err != nil {
		log.Panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	timer := time.Tick(2 * time.Second)
	sleeping := false

	for {
		select {
		case <-timer:
			sleeping = !sleeping
			if sleeping {
				pca9685.Sleep()
			} else {
				pca9685.Wake()
			}
		case <-c:
			return
		}
	}
}
