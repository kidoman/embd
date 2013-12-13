// Package us020 allows interfacing with the US020 ultrasonic range finder.
package us020

import (
	"sync"
	"time"

	"github.com/kid0m4n/go-rpi/sensor/bmp085"
	"github.com/stianeikeland/go-rpio"
)

const (
	pulseDelay = 30000 * time.Nanosecond
)

// A US020 implements access to a US020 ultrasonic range finder.
type US020 interface {
	// Distance computes the distance of the bot from the closest obstruction.
	Distance() (float64, error)
}

type us020 struct {
	echoPinNumber, triggerPinNumber int

	echoPin    rpio.Pin
	triggerPin rpio.Pin

	speedSound float64

	initialized bool
	mu          *sync.RWMutex
}

// New creates a new US020 interface. The bus variable controls
// the I2C bus used to communicate with the device.
func New(e, t int) US020 {
	return &us020{echoPinNumber: e, triggerPinNumber: t, mu: new(sync.RWMutex)}
}

func (d *us020) setup() (err error) {
	d.mu.RLock()
	if d.initialized {
		d.mu.RUnlock()
		return
	}
	d.mu.RUnlock()

	d.mu.Lock()
	defer d.mu.Unlock()

	if err = rpio.Open(); err != nil {
		return
	}

	d.echoPin = rpio.Pin(d.echoPinNumber)       // ECHO port on the US020
	d.triggerPin = rpio.Pin(d.triggerPinNumber) // TRIGGER port on the US020

	temp, err := bmp085.Temperature()
	if err != nil {
		d.speedSound = 340
	} else {
		d.speedSound = 331.4 + 0.606*temp
	}

	d.initialized = true

	return nil
}

// Distance computes the distance of the bot from the closest obstruction.
func (d *us020) Distance() (distance float64, err error) {
	if err = d.setup(); err != nil {
		return
	}

	// Generate a TRIGGER pulse
	d.triggerPin.High()
	time.Sleep(pulseDelay)
	d.triggerPin.Low()

	// Wait until ECHO goes high
	for d.echoPin.Read() == rpio.Low {
	}

	startTime := time.Now() // Record time when ECHO goes high

	// Wait until ECHO goes low
	for d.echoPin.Read() == rpio.High {
	}

	duration := time.Since(startTime) // Calculate time lapsed for ECHO to transition from high to low

	// Calculate the distance based on the time computed
	distance = float64(duration.Nanoseconds()) / 10000000 * (d.speedSound / 2)

	return
}
