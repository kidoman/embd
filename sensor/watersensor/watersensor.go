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

type WaterSensor interface {
    IsWet() (b bool,err error)
}

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

func (d *watersensor) IsWet() (b bool, err error) {
    if err = d.Setup(); err != nil {
        return
    }

    if d.debug {
        log.Print("Getting reading")
    }

    if d.waterPin.Read() == rpio.High {
        b=true
    } else {
        b=false
    }

    return
}
