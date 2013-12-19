// Package BH1750FVI allows interfacing with the BH1750FVI ambient light sensor through I2C protocol.
package bh1750Fvi

import (
	"log"
	"sync"
	"time"

	"github.com/kid0m4n/go-rpi/i2c"
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

// A BH1750VI interface implements access to the sensor.
type BH1750VI interface {
	// Run starts continuous sensor data acquisition loop.
	Run() error

	// Lighting returns the ambient lighting in lx.
	Lighting() (lighting float64, err error)

	// Close.
	Close()

	// SetPollDelay sets the delay between run of data acquisition loop.
	SetPollDelay(delay int)
}

type bh1750vi struct {
	bus i2c.Bus
	mu  *sync.RWMutex

	lightingReadings chan float64
	quit             chan bool

	i2cAddr       byte
	operationCode byte

	poll int
}

// Default instance for BH1750FVI sensor
var Default = New(High, i2c.Default)

// Supports three modes:
// "H" -> High resolution mode (1lx), takes 120ms (recommended).
// "H2" -> High resolution mode 2 (0.5lx), takes 120ms (only use for low light).

// New creates a new BH1750FVI interface according to the mode passed.
func New(mode string, bus i2c.Bus) BH1750VI {
	switch mode {
	case High:
		return &bh1750vi{bus: bus, i2cAddr: sensorI2cAddr, operationCode: highResOpCode, mu: new(sync.RWMutex)}
	case High2:
		return &bh1750vi{bus: bus, i2cAddr: sensorI2cAddr, operationCode: highResMode2OpCode, mu: new(sync.RWMutex)}
	default:
		return &bh1750vi{bus: bus, i2cAddr: sensorI2cAddr, operationCode: highResOpCode, mu: new(sync.RWMutex)}
	}
}

// NewHighMode returns a BH1750FVI inteface on high resolution mode (1lx resolution)
func NewHighMode(bus i2c.Bus) BH1750VI {
	return New(High, bus)
}

// NewHighMode returns a BH1750FVI inteface on high resolution mode2 (0.5lx resolution)
func NewHigh2Mode(bus i2c.Bus) BH1750VI {
	return New(High2, bus)
}

func (d *bh1750vi) measureLighting() (lighting float64, err error) {
	err = d.bus.WriteByte(d.i2cAddr, d.operationCode)
	if err != nil {
		log.Print("bh1750vi: Failed to initialize sensor")
		return
	}
	time.Sleep(180 * time.Millisecond)

	var sensorReading int
	if sensorReading, err = d.bus.ReadInt(d.i2cAddr, defReadReg); err != nil {
		return
	}

	lighting = float64(sensorReading) / measurementAcuuracy
	return
}

// Lighting returns the ambient lighting in lx.
func (d *bh1750vi) Lighting() (lighting float64, err error) {
	select {
	case lighting = <-d.lightingReadings:
		return
	default:
		return d.measureLighting()
	}
}

// Run starts continuous sensor data acquisition loop.
func (d *bh1750vi) Run() (err error) {
	go func() {
		d.quit = make(chan bool)

		timer := time.Tick(time.Duration(d.poll) * time.Millisecond)

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
func (d *bh1750vi) Close() {
	if d.quit != nil {
		d.quit <- true
	}
	return
}

// SetPollDelay sets the delay between run of data acquisition loop.
func (d *bh1750vi) SetPollDelay(delay int) {
	d.poll = delay
}

// SetPollDelay sets the delay between run of data acquisition loop.
func SetPollDelay(delay int) {
	Default.SetPollDelay(delay)
}

// Lighting returns the ambient lighting in lx.
func Lighting() (lighting float64, err error) {
	return Default.Lighting()
}

// Run starts continuous sensor data acquisition loop.
func Run() (err error) {
	return Default.Run()
}

// Close.
func Close() {
	Default.Close()
}
