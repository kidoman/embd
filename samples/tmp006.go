// +build ignore

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/kidoman/embd"
	"github.com/kidoman/embd/sensor/tmp006"

	_ "github.com/kidoman/embd/host/all"
)

func main() {
	flag.Parse()

	if err := embd.InitI2C(); err != nil {
		panic(err)
	}
	defer embd.CloseI2C()

	bus := embd.NewI2CBus(1)

	sensor := tmp006.New(bus, 0x40)
	if status, err := sensor.Present(); err != nil || !status {
		fmt.Println("tmp006: not found")
		fmt.Println(err)
		return
	}
	defer sensor.Close()

	sensor.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)

	for {
		select {
		case temp := <-sensor.ObjTemps():
			fmt.Printf("tmp006: got obj temp %.2f\n", temp)
		case temp := <-sensor.RawDieTemps():
			fmt.Printf("tmp006: got die temp %.2f\n", temp)
		case <-stop:
			return
		}
	}
}
