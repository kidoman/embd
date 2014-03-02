// +build ignore

package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/kidoman/embd"
	"github.com/kidoman/embd/sensor/tmp006"
)

func main() {
	if err := embd.InitI2C(); err != nil {
		panic(err)
	}
	defer embd.CloseI2C()

	bus := embd.NewI2CBus(1)

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
