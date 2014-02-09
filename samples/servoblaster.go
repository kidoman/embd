package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/kidoman/embd/controller/servoblaster"
	"github.com/kidoman/embd/motion/servo"
)

func main() {
	sb := servoblaster.New()
	sb.Debug = true
	defer sb.Close()

	servo := servo.New(sb, 0)
	servo.Debug = true

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

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
			var err error
			switch left {
			case true:
				err = servo.SetAngle(45)
			case false:
				err = servo.SetAngle(135)
			}
			if err != nil {
				log.Panic(err)
			}
		case <-c:
			return
		}
	}
}
