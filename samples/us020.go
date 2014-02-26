package main

import (
	"log"
	"time"

	"github.com/kidoman/embd"
	"github.com/kidoman/embd/sensor/us020"
)

func main() {
	gpio, err := embd.NewGPIO()
	if err != nil {
		panic(err)
	}
	defer gpio.Close()

	echoPin, err := gpio.DigitalPin(10)
	if err != nil {
		panic(err)
	}
	triggerPin, err := gpio.DigitalPin(9)
	if err != nil {
		panic(err)
	}

	rf := us020.New(echoPin, triggerPin, nil)
	defer rf.Close()

	for {
		distance, err := rf.Distance()
		if err != nil {
			log.Panic(err)
		}
		log.Printf("Distance is %v", distance)

		time.Sleep(500 * time.Millisecond)
	}
}
