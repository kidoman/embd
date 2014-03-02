// Package BH1750FVI allows interfacing with the BH1750FVI ambient light sensor through I2C protocol.
package bh1750fvi

import (
	"log"
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

type BH1750FVI struct {
	Bus  embd.I2CBus
	Poll int

	mu sync.RWMutex

	lightingReadings chan float64
	quit             chan bool

	i2cAddr       byte
	operationCode byte
}

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

// NewHighMode returns a BH1750FVI inteface on high resolution mode (1lx resolution)
func NewHighMode(bus embd.I2CBus) *BH1750FVI {
	return New(High, bus)
}

// NewHighMode returns a BH1750FVI inteface on high resolution mode2 (0.5lx resolution)
func NewHigh2Mode(bus embd.I2CBus) *BH1750FVI {
	return New(High2, bus)
}

func (d *BH1750FVI) measureLighting() (lighting float64, err error) {
	err = d.Bus.WriteByte(d.i2cAddr, d.operationCode)
	if err != nil {
		log.Print("bh1750fvi: Failed to initialize sensor")
		return
	}
	time.Sleep(180 * time.Millisecond)

	var reading uint16
	if reading, err = d.Bus.ReadWordFromReg(d.i2cAddr, defReadReg); err != nil {
		return
	}

	lighting = float64(int16(reading)) / measurementAcuuracy
	return
}

// Lighting returns the ambient lighting in lx.
func (d *BH1750FVI) Lighting() (lighting float64, err error) {
	select {
	case lighting = <-d.lightingReadings:
		return
	default:
		return d.measureLighting()
	}
}

// Run starts continuous sensor data acquisition loop.
func (d *BH1750FVI) Run() (err error) {
	go func() {
		d.quit = make(chan bool)

		timer := time.Tick(time.Duration(d.Poll) * time.Millisecond)

		var lighting float64

		for {
			select {
			case d.lightingReadings <- lighting:
			case <-timer:
				if l, err := d.measureLighting(); err == nil {
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
