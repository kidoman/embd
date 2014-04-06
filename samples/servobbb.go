// +build ignore

package main

import (
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/kidoman/embd"
	"github.com/kidoman/embd/motion/servo"

	_ "github.com/kidoman/embd/host/all"
)

func main() {
	flag.Parse()

	embd.InitGPIO()
	defer embd.CloseGPIO()

	pwm, err := embd.NewPWMPin("P9_14")
	if err != nil {
		panic(err)
	}
	defer pwm.Close()

	servo := servo.New(pwm)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

	turnTimer := time.Tick(500 * time.Millisecond)
	left := true

	servo.SetAngle(90)
	defer func() {
		servo.SetAngle(90)
	}()

	for {
		select {
		case <-turnTimer:
			left = !left
			switch left {
			case true:
				servo.SetAngle(70)
			case false:
				servo.SetAngle(110)
			}
		case <-quit:
			return
		}
	}
}
