// Package watersensor allows interfacing with the water sensor.
package watersensor

import (
	"sync"

	"github.com/golang/glog"
	"github.com/kidoman/embd"
)

// WaterSensor represents a water sensor.
type WaterSensor struct {
	Pin embd.DigitalPin

	initialized bool
	mu          sync.RWMutex
}

// New creates a new WaterSensor struct
func New(pin embd.DigitalPin) *WaterSensor {
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

	if err := d.Pin.SetDirection(embd.In); err != nil {
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

	glog.V(1).Infof("watersensor: reading")

	value, err := d.Pin.Read()
	if err != nil {
		return false, err
	}
	if value == embd.High {
		return true, nil
	} else {
		return false, nil
	}
}
