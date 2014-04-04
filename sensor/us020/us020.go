// Package us020 allows interfacing with the US020 ultrasonic range finder.
package us020

import (
	"sync"
	"time"

	"github.com/golang/glog"
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
}

// New creates a new US020 interface. The bus variable controls
// the I2C bus used to communicate with the device.
func New(echoPin, triggerPin embd.DigitalPin, thermometer Thermometer) *US020 {
	return &US020{EchoPin: echoPin, TriggerPin: triggerPin, Thermometer: thermometer}
}

func (d *US020) setup() error {
	d.mu.RLock()
	if d.initialized {
		d.mu.RUnlock()
		return nil
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

		glog.V(1).Infof("us020: read a temperature of %v, so speed of sound = %v", temp, d.speedSound)
	} else {
		d.speedSound = 340
	}

	d.initialized = true

	return nil
}

// Distance computes the distance of the bot from the closest obstruction.
func (d *US020) Distance() (float64, error) {
	if err := d.setup(); err != nil {
		return 0, err
	}

	glog.V(2).Infof("us020: trigerring pulse")

	// Generate a TRIGGER pulse
	d.TriggerPin.Write(embd.High)
	time.Sleep(pulseDelay)
	d.TriggerPin.Write(embd.Low)

	glog.V(2).Infof("us020: waiting for echo to go high")

	duration, err := d.EchoPin.TimePulse(embd.High)
	if err != nil {
		return 0, err
	}

	// Calculate the distance based on the time computed
	distance := float64(duration.Nanoseconds()) / 10000000 * (d.speedSound / 2)

	return distance, nil
}

// Close.
func (d *US020) Close() error {
	return d.EchoPin.SetDirection(embd.Out)
}
