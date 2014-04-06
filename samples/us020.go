// +build ignore

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/kidoman/embd"
	"github.com/kidoman/embd/sensor/us020"

	_ "github.com/kidoman/embd/host/all"
)

func main() {
	flag.Parse()

	if err := embd.InitGPIO(); err != nil {
		panic(err)
	}
	defer embd.CloseGPIO()

	echoPin, err := embd.NewDigitalPin("P9_21")
	if err != nil {
		panic(err)
	}

	triggerPin, err := embd.NewDigitalPin("P9_22")
	if err != nil {
		panic(err)
	}

	rf := us020.New(echoPin, triggerPin, nil)
	defer rf.Close()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

	for {
		select {
		default:
			distance, err := rf.Distance()
			if err != nil {
				panic(err)
			}
			fmt.Printf("Distance is %v\n", distance)

			time.Sleep(500 * time.Millisecond)
		case <-quit:
			return
		}
	}
}
