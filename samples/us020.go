package main

import (
	"log"
	"time"

	"github.com/kid0m4n/go-rpi/sensor/us020"
)

func main() {
	rangeFinder := us020.New(10, 9)
	defer rangeFinder.Close()

	for {
		distance, err := rangeFinder.Distance()
		if err != nil {
			log.Panic(err)
		}
		log.Printf("Distance is %v", distance)

		time.Sleep(500 * time.Millisecond)
	}
}
