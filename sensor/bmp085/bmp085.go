// Package bmp085 allows interfacing with Bosch BMP085 barometric pressure sensor. This sensor
// has the ability to provided compensated temperature and pressure readings.
package bmp085

import (
	"log"
	"math"
	"sync"
	"time"

	"github.com/kid0m4n/go-rpi/i2c"
)

const (
	address = 0x77

	calAc1          = 0xAA
	calAc2          = 0xAC
	calAc3          = 0xAE
	calAc4          = 0xB0
	calAc5          = 0xB2
	calAc6          = 0xB4
	calB1           = 0xB6
	calB2           = 0xB8
	calMB           = 0xBA
	calMC           = 0xBC
	calMD           = 0xBE
	control         = 0xF4
	tempData        = 0xF6
	pressureData    = 0xF6
	readTempCmd     = 0x2E
	readPressureCmd = 0x34

	tempReadDelay = 5 * time.Millisecond

	p0 = 101325

	pollDelay = 250
)

// A BMP085 implements access to the Bosch BMP085 sensor.
type BMP085 interface {
	// SetPollDelay sets the delay between runs of the data acquisition loop.
	SetPollDelay(delay int)

	// Temperature returns the current temperature reading.
	Temperature() (temp float64, err error)
	// Pressure returns the current pressure reading.
	Pressure() (pressure int, err error)
	// Altitude returns the current altitude reading.
	Altitude() (altitude float64, err error)

	// Run starts the sensor data acquisition loop.
	Run() error
	// Close.
	Close()
}

type bmp085 struct {
	bus i2c.Bus
	oss uint

	ac1, ac2, ac3      int16
	ac4, ac5, ac6      uint16
	b1, b2, mb, mc, md int16
	b5                 int32
	calibrated         bool
	cmu                *sync.RWMutex

	poll      int
	temps     chan uint16
	pressures chan int32
	altitudes chan float64
	quit      chan struct{}

	debug bool
}

// Default instance of the BMP085 sensor.
var Default = New(i2c.Default)

// New creates a new BMP085 interface. The bus variable controls
// the I2C bus used to communicate with the device.
func New(bus i2c.Bus) BMP085 {
	return &bmp085{bus: bus, cmu: new(sync.RWMutex), poll: pollDelay}
}

// SetPollDelay sets the delay between runs of the data acquisition loop.
func (d *bmp085) SetPollDelay(delay int) {
	d.poll = delay
}

func (d *bmp085) calibrate() (err error) {
	d.cmu.RLock()
	if d.calibrated {
		d.cmu.RUnlock()
		return
	}
	d.cmu.RUnlock()

	d.cmu.Lock()
	defer d.cmu.Unlock()

	readInt16 := func(reg byte) (value int16, err error) {
		var v uint16
		if v, err = d.bus.ReadWordFromReg(address, reg); err != nil {
			return
		}
		value = int16(v)
		return
	}

	readUInt16 := func(reg byte) (value uint16, err error) {
		var v uint16
		if v, err = d.bus.ReadWordFromReg(address, reg); err != nil {
			return
		}
		value = uint16(v)
		return
	}

	d.ac1, err = readInt16(calAc1)
	if err != nil {
		return
	}
	d.ac2, err = readInt16(calAc2)
	if err != nil {
		return
	}
	d.ac3, err = readInt16(calAc3)
	if err != nil {
		return
	}
	d.ac4, err = readUInt16(calAc4)
	if err != nil {
		return
	}
	d.ac5, err = readUInt16(calAc5)
	if err != nil {
		return
	}
	d.ac6, err = readUInt16(calAc6)
	if err != nil {
		return
	}
	d.b1, err = readInt16(calB1)
	if err != nil {
		return
	}
	d.b2, err = readInt16(calB2)
	if err != nil {
		return
	}
	d.mb, err = readInt16(calMB)
	if err != nil {
		return
	}
	d.mc, err = readInt16(calMC)
	if err != nil {
		return
	}
	d.md, err = readInt16(calMD)
	if err != nil {
		return
	}

	d.calibrated = true

	if d.debug {
		log.Print("bmp085: calibration data retrieved")
		log.Printf("bmp085: param AC1 = %v", d.ac1)
		log.Printf("bmp085: param AC2 = %v", d.ac2)
		log.Printf("bmp085: param AC3 = %v", d.ac3)
		log.Printf("bmp085: param AC4 = %v", d.ac4)
		log.Printf("bmp085: param AC5 = %v", d.ac5)
		log.Printf("bmp085: param AC6 = %v", d.ac6)
		log.Printf("bmp085: param B1 = %v", d.b1)
		log.Printf("bmp085: param B2 = %v", d.b2)
		log.Printf("bmp085: param MB = %v", d.mb)
		log.Printf("bmp085: param MC = %v", d.mc)
		log.Printf("bmp085: param MD = %v", d.md)
	}

	return
}

func (d *bmp085) readUncompensatedTemp() (temp uint16, err error) {
	if err = d.bus.WriteByteToReg(address, control, readTempCmd); err != nil {
		return
	}
	time.Sleep(tempReadDelay)
	if temp, err = d.bus.ReadWordFromReg(address, tempData); err != nil {
		return
	}
	return
}

func (d *bmp085) calcTemp(utemp uint16) uint16 {
	x1 := ((int(utemp) - int(d.ac6)) * int(d.ac5)) >> 15
	x2 := (int(d.mc) << 11) / (x1 + int(d.md))

	d.cmu.Lock()
	d.b5 = int32(x1 + x2)
	d.cmu.Unlock()

	return uint16((d.b5 + 8) >> 4)
}

func (d *bmp085) measureTemp() (temp uint16, err error) {
	if err = d.calibrate(); err != nil {
		return
	}

	var utemp uint16
	if utemp, err = d.readUncompensatedTemp(); err != nil {
		return
	}
	if d.debug {
		log.Printf("bcm085: uncompensated temp: %v", utemp)
	}
	temp = d.calcTemp(utemp)
	if d.debug {
		log.Printf("bcm085: compensated temp %v", temp)
	}
	return
}

// Temperature returns the current temperature reading.
func (d *bmp085) Temperature() (temp float64, err error) {

	select {
	case t := <-d.temps:
		temp = float64(t) / 10
		return
	default:
		if d.debug {
			log.Print("bcm085: no temps available... measuring")
		}
		var t uint16
		t, err = d.measureTemp()
		if err != nil {
			return
		}
		temp = float64(t) / 10
		return
	}
}

func (d *bmp085) readUncompensatedPressure() (pressure uint32, err error) {
	if err = d.bus.WriteByteToReg(address, control, byte(readPressureCmd+(d.oss<<6))); err != nil {
		return
	}
	time.Sleep(time.Duration(2+(3<<d.oss)) * time.Millisecond)

	data := make([]byte, 3)
	if err = d.bus.ReadFromReg(address, pressureData, data); err != nil {
		return
	}

	pressure = ((uint32(data[0]) << 16) | (uint32(data[1]) << 8) | uint32(data[2])) >> (8 - d.oss)

	return
}

func (d *bmp085) calcPressure(upressure uint32) (p int32) {
	var x1, x2, x3 int32

	l := func(s string, v interface{}) {
		if d.debug {
			log.Printf("bcm085: %v = %v", s, v)
		}
	}

	b6 := d.b5 - 4000
	l("b6", b6)

	// Calculate b3
	x1 = (int32(d.b2) * int32(b6*b6) >> 12) >> 11
	x2 = (int32(d.ac2) * b6) >> 11
	x3 = x1 + x2
	b3 := (((int32(d.ac1)*4 + x3) << d.oss) + 2) >> 2

	l("x1", x1)
	l("x2", x2)
	l("x3", x3)
	l("b3", b3)

	// Calculate b4
	x1 = (int32(d.ac3) * b6) >> 13
	x2 = (int32(d.b1) * ((b6 * b6) >> 12)) >> 16
	x3 = ((x1 + x2) + 2) >> 2
	b4 := (uint32(d.ac4) * uint32(x3+32768)) >> 15

	l("x1", x1)
	l("x2", x2)
	l("x3", x3)
	l("b4", b4)

	b7 := (uint32(upressure-uint32(b3)) * (50000 >> d.oss))
	if b7 < 0x80000000 {
		p = int32((b7 << 1) / b4)
	} else {
		p = int32((b7 / b4) << 1)
	}
	l("b7", b7)
	l("p", p)

	x1 = (p >> 8) * (p >> 8)
	x1 = (x1 * 3038) >> 16
	x2 = (-7357 * p) >> 16
	p += (x1 + x2 + 3791) >> 4

	l("x1", x1)
	l("x2", x2)
	l("x3", x3)
	l("p", p)

	return
}

func (d *bmp085) calcAltitude(pressure int32) float64 {
	return 44330 * (1 - math.Pow(float64(pressure)/p0, 0.190295))
}

func (d *bmp085) measurePressureAndAltitude() (pressure int32, altitude float64, err error) {
	if err = d.calibrate(); err != nil {
		return
	}

	var upressure uint32
	if upressure, err = d.readUncompensatedPressure(); err != nil {
		return
	}
	if d.debug {
		log.Printf("bcm085: uncompensated pressure: %v", upressure)
	}
	pressure = d.calcPressure(upressure)
	if d.debug {
		log.Printf("bcm085: compensated pressure %v", pressure)
	}
	altitude = d.calcAltitude(pressure)
	if d.debug {
		log.Printf("bcm085: calculated altitude %v", altitude)
	}
	return
}

// Pressure returns the current pressure reading.
func (d *bmp085) Pressure() (pressure int, err error) {
	if err = d.calibrate(); err != nil {
		return
	}

	select {
	case p := <-d.pressures:
		pressure = int(p)
		return
	default:
		if d.debug {
			log.Print("bcm085: no pressures available... measuring")
		}
		var p int32
		p, _, err = d.measurePressureAndAltitude()
		if err != nil {
			return
		}
		pressure = int(p)
		return
	}
}

// Altitude returns the current altitude reading.
func (d *bmp085) Altitude() (altitude float64, err error) {
	if err = d.calibrate(); err != nil {
		return
	}

	select {
	case altitude = <-d.altitudes:
		return
	default:
		if d.debug {
			log.Print("bcm085: no altitudes available... measuring")
		}
		_, altitude, err = d.measurePressureAndAltitude()
		if err != nil {
			return
		}
		return
	}
}

// Run starts the sensor data acquisition loop.
func (d *bmp085) Run() (err error) {
	go func() {
		d.quit = make(chan struct{})
		timer := time.Tick(time.Duration(d.poll) * time.Millisecond)

		var temp uint16
		var pressure int32
		var altitude float64

		for {
			select {
			case <-timer:
				t, err := d.measureTemp()
				if err == nil {
					temp = t
				}
				if err == nil && d.temps == nil {
					d.temps = make(chan uint16)
				}
				p, a, err := d.measurePressureAndAltitude()
				if err == nil {
					pressure = p
					altitude = a
				}
				if err == nil && d.pressures == nil && d.altitudes == nil {
					d.pressures = make(chan int32)
					d.altitudes = make(chan float64)
				}
			case d.temps <- temp:
			case d.pressures <- pressure:
			case d.altitudes <- altitude:
			case <-d.quit:
				d.temps = nil
				d.pressures = nil
				d.altitudes = nil
				return
			}
		}
	}()

	return
}

// Close.
func (d *bmp085) Close() {
	if d.quit != nil {
		d.quit <- struct{}{}
	}
}

// SetPollDelay sets the delay between runs of the data acquisition loop.
func SetPollDelay(delay int) {
	Default.SetPollDelay(delay)
}

// Temperature returns the current temperature reading.
func Temperature() (temp float64, err error) {
	return Default.Temperature()
}

// Pressure returns the current pressure reading.
func Pressure() (pressure int, err error) {
	return Default.Pressure()
}

// Altitude returns the current altitude reading.
func Altitude() (altitude float64, err error) {
	return Default.Altitude()
}

// Run starts the sensor data acquisition loop.
func Run() (err error) {
	return Default.Run()
}

// Close.
func Close() {
	Default.Close()
}
