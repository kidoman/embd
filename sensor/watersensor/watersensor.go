// Package watersensor allows interfacing with the water sensor
package watersensor

import (
        "log"
        "sync"
        "github.com/stianeikeland/go-rpio"
       )

type watersensor struct {
    waterPinNumber int
    waterPin rpio.Pin

    initialized bool
    mu *sync.RWMutex

    debug bool
}

// WaterSensor implements access to a water sensor
type WaterSensor interface {
    // IsWet determines if there is water present on the sensor
    IsWet() (b bool,err error)
}

// New creates a new WaterSensor interface
func New(pinNumber int) WaterSensor {
    return &watersensor{waterPinNumber: pinNumber, mu: new(sync.RWMutex)}
}

func (d *watersensor) Setup() (err error) {
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

    d.waterPin = rpio.Pin(d.waterPinNumber)
    d.waterPin.Input()
    d.initialized = true

    return nil
}

// IsWet determines if there is water present on the sensor
func (d *watersensor) IsWet() (b bool, err error) {
    if err = d.Setup(); err != nil {
        return
    }

    if d.debug {
        log.Print("Getting reading")
    }

    // Read the pin value of the sensor
    if d.waterPin.Read() == rpio.High {
        b=true
    } else {
        b=false
    }

    return
}
