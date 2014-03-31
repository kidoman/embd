// Package BH1750FVI allows interfacing with the BH1750FVI ambient light sensor through I2C.
package bh1750fvi

import (
	"sync"
	"time"

	"github.com/kidoman/embd"
)

//accuracy = sensorValue/actualValue] (min = 0.96, typ = 1.2, max = 1.44
const (
	High  = "H"
	High2 = "H2"

	measurementAcuuracy = 1.2
	defReadReg          = 0x00

	sensorI2cAddr = 0x23

	highResOpCode      = 0x10
	highResMode2OpCode = 0x11

	pollDelay = 150
)

// BH1750FVI represents a BH1750FVI ambient light sensor.
type BH1750FVI struct {
	Bus  embd.I2CBus
	Poll int

	mu sync.RWMutex

	lightingReadings chan float64
	quit             chan bool

	i2cAddr       byte
	operationCode byte
}

// New returns a BH1750FVI sensor at the specific resolution mode.
func New(mode string, bus embd.I2CBus) *BH1750FVI {
	switch mode {
	case High:
		return &BH1750FVI{Bus: bus, i2cAddr: sensorI2cAddr, operationCode: highResOpCode, Poll: pollDelay}
	case High2:
		return &BH1750FVI{Bus: bus, i2cAddr: sensorI2cAddr, operationCode: highResMode2OpCode, Poll: pollDelay}
	default:
		return &BH1750FVI{Bus: bus, i2cAddr: sensorI2cAddr, operationCode: highResOpCode, Poll: pollDelay}
	}
}

// NewHighMode returns a BH1750FVI sensor on high resolution mode (1lx resolution)
func NewHighMode(bus embd.I2CBus) *BH1750FVI {
	return New(High, bus)
}

// NewHighMode returns a BH1750FVI sensor on high resolution mode2 (0.5lx resolution)
func NewHigh2Mode(bus embd.I2CBus) *BH1750FVI {
	return New(High2, bus)
}

func (d *BH1750FVI) measureLighting() (float64, error) {
	if err := d.Bus.WriteByte(d.i2cAddr, d.operationCode); err != nil {
		return 0, err
	}
	time.Sleep(180 * time.Millisecond)

	reading, err := d.Bus.ReadWordFromReg(d.i2cAddr, defReadReg)
	if err != nil {
		return 0, err
	}

	return float64(int16(reading)) / measurementAcuuracy, nil
}

// Lighting returns the ambient lighting in lx.
func (d *BH1750FVI) Lighting() (float64, error) {
	select {
	case lighting := <-d.lightingReadings:
		return lighting, nil
	default:
		return d.measureLighting()
	}
}

// Run starts continuous sensor data acquisition loop.
func (d *BH1750FVI) Run() {
	go func() {
		d.quit = make(chan bool)

		timer := time.Tick(time.Duration(d.Poll) * time.Millisecond)

		var lighting float64

		for {
			select {
			case d.lightingReadings <- lighting:
			case <-timer:
				l, err := d.measureLighting()
				if err == nil {
					lighting = l
				}
				if err == nil && d.lightingReadings == nil {
					d.lightingReadings = make(chan float64)
				}
			case <-d.quit:
				d.lightingReadings = nil
				return
			}
		}
	}()
	return
}

// Close.
func (d *BH1750FVI) Close() {
	if d.quit != nil {
		d.quit <- true
	}
	return
}
