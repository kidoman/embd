// +build ignore

package main

// Control a stepper motor (28BJY-48)
//
// Datasheet:
// http://www.raspberrypi-spy.co.uk/wp-content/uploads/2012/07/Stepper-Motor-28BJY-48-Datasheet.pdf
//
// this is a port of Matt Hawkins' example impl from
// http://www.raspberrypi-spy.co.uk/2012/07/stepper-motor-control-in-python/
// (link privides additional instructions for wiring your pi)

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi"
)

func main() {
	stepDelay := flag.Int("step-delay", 10, "milliseconds between steps")
	flag.Parse()

	if err := embd.InitGPIO(); err != nil {
		panic(err)
	}
	defer embd.CloseGPIO()

	// Physical pins 11,15,16,18 on rasp pi
	// GPIO17,GPIO22,GPIO23,GPIO24
	stepPinNums := []int{17, 22, 23, 24}

	stepPins := make([]embd.DigitalPin, 4)

	for i, pinNum := range stepPinNums {
		pin, err := embd.NewDigitalPin(pinNum)
		if err != nil {
			panic(err)
		}
		defer pin.Close()
		if err := pin.SetDirection(embd.Out); err != nil {
			panic(err)
		}
		if err := pin.Write(embd.Low); err != nil {
			panic(err)
		}
		defer pin.SetDirection(embd.In)

		stepPins[i] = pin
	}

	// Define sequence described in manufacturer's datasheet
	seq := [][]int{
		[]int{1, 0, 0, 0},
		[]int{1, 1, 0, 0},
		[]int{0, 1, 0, 0},
		[]int{0, 1, 1, 0},
		[]int{0, 0, 1, 0},
		[]int{0, 0, 1, 1},
		[]int{0, 0, 0, 1},
		[]int{1, 0, 0, 1},
	}
	stepCount := len(seq) - 1
	stepDir := 2 // Set to 1 or 2 for clockwise, -1 or -2 for counter-clockwise

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)
	defer signal.Stop(quit)

	// Start main loop
	ticker := time.NewTicker(time.Duration(*stepDelay) * time.Millisecond)
	defer ticker.Stop()

	var stepCounter int
	for {
		select {
		case <-ticker.C:
			// set pins to appropriate values for given position in the sequence
			for i, pin := range stepPins {
				if seq[stepCounter][i] != 0 {
					fmt.Printf("Enable pin %d, step %d\n", i, stepCounter)
					if err := pin.Write(embd.High); err != nil {
						panic(err)
					}
				} else {
					if err := pin.Write(embd.Low); err != nil {
						panic(err)
					}
				}
			}
			stepCounter += stepDir

			// If we reach the end of the sequence start again
			if stepCounter >= stepCount {
				stepCounter = 0
			} else if stepCounter < 0 {
				stepCounter = stepCount
			}

		case <-quit:
			return
		}
	}
}
