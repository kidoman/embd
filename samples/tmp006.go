package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/kid0m4n/go-rpi/i2c"
	"github.com/kid0m4n/go-rpi/sensor/tmp006"
)

func main() {
	bus := i2c.NewBus(1)

	sensor := tmp006.New(bus, 0x40)
	if status, err := sensor.Present(); err != nil || !status {
		log.Print("tmp006: not found")
		log.Print(err)
		return
	}
	defer sensor.Close()

	sensor.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)

	for {
		select {
		case temp := <-sensor.ObjTemps():
			log.Printf("tmp006: got obj temp %.2f", temp)
		case temp := <-sensor.RawDieTemps():
			log.Printf("tmp006: got die temp %.2f", temp)
		case <-stop:
			return
		}
	}
}
