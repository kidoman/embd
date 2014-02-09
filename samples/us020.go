package main

import (
	"log"
	"time"

	"github.com/kidoman/embd/sensor/us020"
	"github.com/stianeikeland/go-rpio"
)

func main() {
	rpio.Open()
	defer rpio.Close()

	rf := us020.New(10, 9, nil)
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
