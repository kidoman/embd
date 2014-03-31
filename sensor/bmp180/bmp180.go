// Package bmp180 allows interfacing with Bosch BMP180 barometric pressure sensor. This sensor
// has the ability to provided compensated temperature and pressure readings.
package bmp180

import (
	"math"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/kidoman/embd"
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

// BMP180 represents a Bosch BMP180 barometric sensor.
type BMP180 struct {
	Bus  embd.I2CBus
	Poll int

	oss uint

	ac1, ac2, ac3      int16
	ac4, ac5, ac6      uint16
	b1, b2, mb, mc, md int16
	b5                 int32
	calibrated         bool
	cmu                sync.RWMutex

	temps     chan uint16
	pressures chan int32
	altitudes chan float64
	quit      chan struct{}
}

// New returns a handle to a BMP180 sensor.
func New(bus embd.I2CBus) *BMP180 {
	return &BMP180{Bus: bus, Poll: pollDelay}
}

func (d *BMP180) calibrate() error {
	d.cmu.RLock()
	if d.calibrated {
		d.cmu.RUnlock()
		return nil
	}
	d.cmu.RUnlock()

	d.cmu.Lock()
	defer d.cmu.Unlock()

	readInt16 := func(reg byte) (int16, error) {
		v, err := d.Bus.ReadWordFromReg(address, reg)
		if err != nil {
			return 0, err
		}
		return int16(v), nil
	}

	readUInt16 := func(reg byte) (uint16, error) {
		v, err := d.Bus.ReadWordFromReg(address, reg)
		if err != nil {
			return 0, err
		}
		return uint16(v), nil
	}

	var err error
	d.ac1, err = readInt16(calAc1)
	if err != nil {
		return err
	}
	d.ac2, err = readInt16(calAc2)
	if err != nil {
		return err
	}
	d.ac3, err = readInt16(calAc3)
	if err != nil {
		return err
	}
	d.ac4, err = readUInt16(calAc4)
	if err != nil {
		return err
	}
	d.ac5, err = readUInt16(calAc5)
	if err != nil {
		return err
	}
	d.ac6, err = readUInt16(calAc6)
	if err != nil {
		return err
	}
	d.b1, err = readInt16(calB1)
	if err != nil {
		return err
	}
	d.b2, err = readInt16(calB2)
	if err != nil {
		return err
	}
	d.mb, err = readInt16(calMB)
	if err != nil {
		return err
	}
	d.mc, err = readInt16(calMC)
	if err != nil {
		return err
	}
	d.md, err = readInt16(calMD)
	if err != nil {
		return err
	}

	d.calibrated = true

	if glog.V(1) {
		glog.Info("bmp180: calibration data retrieved")
		glog.Infof("bmp180: param AC1 = %v", d.ac1)
		glog.Infof("bmp180: param AC2 = %v", d.ac2)
		glog.Infof("bmp180: param AC3 = %v", d.ac3)
		glog.Infof("bmp180: param AC4 = %v", d.ac4)
		glog.Infof("bmp180: param AC5 = %v", d.ac5)
		glog.Infof("bmp180: param AC6 = %v", d.ac6)
		glog.Infof("bmp180: param B1 = %v", d.b1)
		glog.Infof("bmp180: param B2 = %v", d.b2)
		glog.Infof("bmp180: param MB = %v", d.mb)
		glog.Infof("bmp180: param MC = %v", d.mc)
		glog.Infof("bmp180: param MD = %v", d.md)
	}

	return nil
}

func (d *BMP180) readUncompensatedTemp() (uint16, error) {
	if err := d.Bus.WriteByteToReg(address, control, readTempCmd); err != nil {
		return 0, err
	}
	time.Sleep(tempReadDelay)
	temp, err := d.Bus.ReadWordFromReg(address, tempData)
	if err != nil {
		return 0, err
	}
	return temp, nil
}

func (d *BMP180) calcTemp(utemp uint16) uint16 {
	x1 := ((int(utemp) - int(d.ac6)) * int(d.ac5)) >> 15
	x2 := (int(d.mc) << 11) / (x1 + int(d.md))

	d.cmu.Lock()
	d.b5 = int32(x1 + x2)
	d.cmu.Unlock()

	return uint16((d.b5 + 8) >> 4)
}

func (d *BMP180) measureTemp() (uint16, error) {
	if err := d.calibrate(); err != nil {
		return 0, err
	}

	utemp, err := d.readUncompensatedTemp()
	if err != nil {
		return 0, err
	}
	glog.V(1).Infof("bcm085: uncompensated temp: %v", utemp)
	temp := d.calcTemp(utemp)
	glog.V(1).Infof("bcm085: compensated temp %v", temp)
	return temp, nil
}

// Temperature returns the current temperature reading.
func (d *BMP180) Temperature() (float64, error) {
	select {
	case t := <-d.temps:
		temp := float64(t) / 10
		return temp, nil
	default:
		glog.V(1).Infof("bcm085: no temps available... measuring")
		t, err := d.measureTemp()
		if err != nil {
			return 0, err
		}
		temp := float64(t) / 10
		return temp, nil
	}
}

func (d *BMP180) readUncompensatedPressure() (uint32, error) {
	if err := d.Bus.WriteByteToReg(address, control, byte(readPressureCmd+(d.oss<<6))); err != nil {
		return 0, err
	}
	time.Sleep(time.Duration(2+(3<<d.oss)) * time.Millisecond)

	data := make([]byte, 3)
	if err := d.Bus.ReadFromReg(address, pressureData, data); err != nil {
		return 0, err
	}

	pressure := ((uint32(data[0]) << 16) | (uint32(data[1]) << 8) | uint32(data[2])) >> (8 - d.oss)

	return pressure, nil
}

func (d *BMP180) calcPressure(upressure uint32) int32 {
	var x1, x2, x3 int32

	l := func(s string, v interface{}) {
		glog.V(1).Infof("bcm085: %v = %v", s, v)
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

	var p int32
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

	return p
}

func (d *BMP180) calcAltitude(pressure int32) float64 {
	return 44330 * (1 - math.Pow(float64(pressure)/p0, 0.190295))
}

func (d *BMP180) measurePressureAndAltitude() (int32, float64, error) {
	if err := d.calibrate(); err != nil {
		return 0, 0, err
	}

	upressure, err := d.readUncompensatedPressure()
	if err != nil {
		return 0, 0, err
	}
	glog.V(1).Infof("bcm085: uncompensated pressure: %v", upressure)
	pressure := d.calcPressure(upressure)
	glog.V(1).Infof("bcm085: compensated pressure %v", pressure)
	altitude := d.calcAltitude(pressure)
	glog.V(1).Infof("bcm085: calculated altitude %v", altitude)
	return pressure, altitude, nil
}

// Pressure returns the current pressure reading.
func (d *BMP180) Pressure() (int, error) {
	if err := d.calibrate(); err != nil {
		return 0, err
	}

	select {
	case p := <-d.pressures:
		return int(p), nil
	default:
		glog.V(1).Infof("bcm085: no pressures available... measuring")
		p, _, err := d.measurePressureAndAltitude()
		if err != nil {
			return 0, err
		}
		return int(p), nil
	}
}

// Altitude returns the current altitude reading.
func (d *BMP180) Altitude() (float64, error) {
	if err := d.calibrate(); err != nil {
		return 0, err
	}

	select {
	case altitude := <-d.altitudes:
		return altitude, nil
	default:
		glog.V(1).Info("bcm085: no altitudes available... measuring")
		_, altitude, err := d.measurePressureAndAltitude()
		if err != nil {
			return 0, err
		}
		return altitude, nil
	}
}

// Run starts the sensor data acquisition loop.
func (d *BMP180) Run() {
	go func() {
		d.quit = make(chan struct{})
		timer := time.Tick(time.Duration(d.Poll) * time.Millisecond)

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
func (d *BMP180) Close() {
	if d.quit != nil {
		d.quit <- struct{}{}
	}
}
