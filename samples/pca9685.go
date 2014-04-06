// +build ignore

package main

import (
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/kidoman/embd"
	"github.com/kidoman/embd/controller/pca9685"

	_ "github.com/kidoman/embd/host/all"
)

func main() {
	flag.Parse()

	if err := embd.InitI2C(); err != nil {
		panic(err)
	}
	defer embd.CloseI2C()

	bus := embd.NewI2CBus(1)

	pca9685 := pca9685.New(bus, 0x41)
	pca9685.Freq = 1000
	defer pca9685.Close()

	if err := pca9685.SetPwm(15, 0, 2000); err != nil {
		panic(err)
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
