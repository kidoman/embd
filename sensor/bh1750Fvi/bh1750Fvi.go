// Package BH1750FVI allows interfacing with the BH1750FVI ambient light sensor through I2C protocol
package bh1750Fvi

import (
	"log"
	"sync"
	"time"

	"github.com/kid0m4n/go-rpi/i2c"
)

const (
	measurementAcuuracy = 1.2 // [accuracy = sensorValue/actualValue] (min = 0.96, typ = 1.2, max = 1.44)

	highResolutionReadAddress   = 0x23
	highResolutionOperationCode = 0x10

	lowResolutionReadAddress   = 0x5c
	lowResolutionOperationCode = 0x23

	highResolutionMode2ReadAddress   = 0x23
	highResolutionMode2OperationCode = 0x11

	pollDelay = 150
)

type BH1750VI interface {
	Start() error

	Lighting() (lighting float64, err error)

	Close()

	SetPollDelay(delay int)
}

type bh1750vi struct {
	bus i2c.Bus
	mu  *sync.RWMutex

	lightingReadings chan float64
	quit             chan bool

	ready bool

	readRegisterAddress byte
	operationCode       byte

	continuousMode bool

	poll int
}

var Default = New("H", i2c.Default)

func New(mode string, bus i2c.Bus) BH1750VI {

	/*
		Supports three modes:
		"L" -> Low resolution mode (4lx), takes 16ms
		"H" -> High resolution mode (1lx), takes 120ms (recommended)
		"H2" -> High resolution mode 2 (0.5lx), takes 120ms (only use for low light)
	*/

	switch mode {
	case "L":
		return &bh1750vi{bus: bus, readRegisterAddress: lowResolutionReadAddress, operationCode: lowResolutionOperationCode, continuousMode: false}
	case "H":
		return &bh1750vi{bus: bus, readRegisterAddress: highResolutionReadAddress, operationCode: highResolutionOperationCode, continuousMode: true}
	case "H2":
		return &bh1750vi{bus: bus, readRegisterAddress: highResolutionMode2ReadAddress, operationCode: highResolutionMode2OperationCode, continuousMode: true}
	default:
		return &bh1750vi{bus: bus, readRegisterAddress: highResolutionReadAddress, operationCode: highResolutionOperationCode, continuousMode: true}
	}
}

func (sensor *bh1750vi) setup() (err error) {
	sensor.mu.Lock()

	if sensor.ready && sensor.continuousMode { //for Low resolution mode sensor has to be initialized for every read.
		sensor.mu.Unlock()
		return
	}

	defer sensor.mu.Unlock()

	err = sensor.bus.WriteByte(sensor.readRegisterAddress, sensor.operationCode)
	if err != nil {
		log.Print("bh1750vi: Failed to initialize sensor")
		return
	}

	sensor.ready = true
	return
}

func (sensor *bh1750vi) measureLighting() (lighting float64, err error) {
	if err = sensor.setup(); err != nil {
		return
	}

	sensorData := make([]byte, 2)
	if err = sensor.bus.ReadFromReg(sensor.readRegisterAddress, sensor.operationCode, sensorData); err != nil {
		return
	}

	sensorReading := (int16(sensorData[0] << 8)) | (int16(sensorData[1]))
	lighting = float64(sensorReading) / measurementAcuuracy

	return
}

func (sensor *bh1750vi) Lighting() (lighting float64, err error) {
	select {
	case lighting = <-sensor.lightingReadings:
		return
	default:
		return sensor.measureLighting()

	}
}

func (sensor *bh1750vi) Start() (err error) {
	go func() {
		sensor.quit = make(chan bool)

		timer := time.Tick(time.Duration(sensor.poll) * time.Millisecond)

		var lighting float64

		for {
			select {
			case sensor.lightingReadings <- lighting:
			case <-timer:
				if l, err := sensor.measureLighting(); err == nil {
					lighting = l
				}
				if err == nil && sensor.lightingReadings == nil {
					sensor.lightingReadings = make(chan float64)
				}
			case <-sensor.quit:
				sensor.lightingReadings = nil
				return
			}
		}
	}()
	return
}

func (sensor *bh1750vi) Close() {
	if sensor.quit != nil {
		sensor.quit <- true
	}
	return
}

func (sensor *bh1750vi) SetPollDelay(delay int) {
	sensor.poll = delay
}

func SetPollDelay(delay int) {
	Default.SetPollDelay(delay)
}

func Lighting() (lighting float64, err error) {
	return Default.Lighting()
}

func Start() (err error) {
	return Default.Start()
}

func Close() {
	Default.Close()
}
