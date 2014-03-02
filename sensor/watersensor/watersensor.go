// Package watersensor allows interfacing with the water sensor
package watersensor

import (
	"sync"

	"github.com/golang/glog"
	"github.com/kidoman/embd/gpio"
)

type WaterSensor struct {
	Pin gpio.DigitalPin

	initialized bool
	mu          sync.RWMutex

	Debug bool
}

// New creates a new WaterSensor struct
func New(pin gpio.DigitalPin) *WaterSensor {
	return &WaterSensor{Pin: pin}
}

func (d *WaterSensor) setup() error {
	d.mu.RLock()
	if d.initialized {
		d.mu.RUnlock()
		return nil
	}
	d.mu.RUnlock()

	d.mu.Lock()
	defer d.mu.Unlock()

	if err := d.Pin.SetDirection(gpio.In); err != nil {
		return err
	}
	d.initialized = true

	return nil
}

// IsWet determines if there is water present on the sensor
func (d *WaterSensor) IsWet() (bool, error) {
	if err := d.setup(); err != nil {
		return false, err
	}

	if d.Debug {
		glog.Infof("watersensor: reading")
	}

	value, err := d.Pin.Read()
	if err != nil {
		return false, err
	}
	if value == gpio.High {
		return true, nil
	} else {
		return false, nil
	}
}
