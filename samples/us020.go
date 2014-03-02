// +build ignore

package main

import (
	"log"
	"time"

	"github.com/kidoman/embd"
	"github.com/kidoman/embd/sensor/us020"
)

func main() {
	if err := embd.InitGPIO(); err != nil {
		panic(err)
	}
	defer embd.CloseGPIO()

	echoPin, err := embd.NewDigitalPin(10)
	if err != nil {
		panic(err)
	}
	triggerPin, err := embd.NewDigitalPin(9)
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
