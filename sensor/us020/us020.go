// Package us020 allows interfacing with the US020 ultrasonic range finder.
package us020

import (
	"log"
	"sync"
	"time"

	"github.com/kidoman/embd"
)

const (
	pulseDelay  = 30000 * time.Nanosecond
	defaultTemp = 25
)

type Thermometer interface {
	Temperature() (float64, error)
}

type nullThermometer struct {
}

func (*nullThermometer) Temperature() (float64, error) {
	return defaultTemp, nil
}

var NullThermometer = &nullThermometer{}

// US020 represents a US020 ultrasonic range finder.
type US020 struct {
	EchoPin, TriggerPin embd.DigitalPin

	Thermometer Thermometer

	speedSound float64

	initialized bool
	mu          sync.RWMutex

	Debug bool
}

// New creates a new US020 interface. The bus variable controls
// the I2C bus used to communicate with the device.
func New(echoPin, triggerPin embd.DigitalPin, thermometer Thermometer) *US020 {
	return &US020{EchoPin: echoPin, TriggerPin: triggerPin, Thermometer: thermometer}
}

func (d *US020) setup() (err error) {
	d.mu.RLock()
	if d.initialized {
		d.mu.RUnlock()
		return
	}
	d.mu.RUnlock()

	d.mu.Lock()
	defer d.mu.Unlock()

	d.TriggerPin.SetDirection(embd.Out)
	d.EchoPin.SetDirection(embd.In)

	if d.Thermometer == nil {
		d.Thermometer = NullThermometer
	}

	if temp, err := d.Thermometer.Temperature(); err == nil {
		d.speedSound = 331.3 + 0.606*temp

		if d.Debug {
			log.Printf("read a temperature of %v, so speed of sound = %v", temp, d.speedSound)
		}
	} else {
		d.speedSound = 340
	}

	d.initialized = true

	return
}

// Distance computes the distance of the bot from the closest obstruction.
func (d *US020) Distance() (distance float64, err error) {
	if err = d.setup(); err != nil {
		return
	}

	if d.Debug {
		log.Print("us020: trigerring pulse")
	}

	// Generate a TRIGGER pulse
	d.TriggerPin.Write(embd.High)
	time.Sleep(pulseDelay)
	d.TriggerPin.Write(embd.Low)

	if d.Debug {
		log.Print("us020: waiting for echo to go high")
	}

	// Wait until ECHO goes high
	for {
		v, err := d.EchoPin.Read()
		if err != nil {
			return 0, err
		}

		if v != embd.Low {
			break
		}
	}

	startTime := time.Now() // Record time when ECHO goes high

	if d.Debug {
		log.Print("us020: waiting for echo to go low")
	}

	// Wait until ECHO goes low
	for {
		v, err := d.EchoPin.Read()
		if err != nil {
			return 0, err
		}

		if v != embd.High {
			break
		}
	}

	duration := time.Since(startTime) // Calculate time lapsed for ECHO to transition from high to low

	// Calculate the distance based on the time computed
	distance = float64(duration.Nanoseconds()) / 10000000 * (d.speedSound / 2)

	return
}

// Close.
func (d *US020) Close() {
	d.EchoPin.SetDirection(embd.Out)
}
