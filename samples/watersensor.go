package main

import (
	"log"
	"time"

	"github.com/kid0m4n/go-rpi/sensor/watersensor"
)

func main() {
	fluidSensor := watersensor.New(7)

	for {
		isWater, err := fluidSensor.IsWet()
		if err != nil {
			log.Panic(err)
		}
        if isWater {
		log.Printf("Bot is dry")
        } else {
		log.Printf("Bot is Wet")
        }
		time.Sleep(500 * time.Millisecond)
	}
}
