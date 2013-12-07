// Package lsm303 allows interfacing with the LSM303 magnetometer.
package lsm303

import (
	"log"
	"math"
	"sync"
	"time"

	"github.com/kid0m4n/go-rpi/i2c"
)

const (
	magAddress = 0x1E

	magConfigRegA = 0x00

	MagHz75         = 0x00
	Mag1Hz5         = 0x04
	Mag3Hz          = 0x08
	Mag7Hz5         = 0x0C
	Mag15Hz         = 0x10
	Mag30Hz         = 0x14
	Mag75Hz         = 0x18
	MagNormal       = 0x00
	MagPositiveBias = 0x01
	MagNegativeBias = 0x02

	MagCRADefault = Mag15Hz | MagNormal

	magModeReg = 0x02

	MagContinuous = 0x00
	MagSleep      = 0x03

	MagMRDefault = MagContinuous

	magDataSignal = 0x02
	magData       = 0x03

	pollDelay = 250
)

type LSM303 interface {
	SetPollDelay(delay int)

	Heading() (heading float64, err error)

	Run() error
	Close() error
}

type lsm303 struct {
	bus i2c.Bus

	initialized bool
	mu          *sync.RWMutex

	headings chan float64

	poll int
	quit chan struct{}

	debug bool
}

var Default = New(i2c.Default)

// New creates a new LSM303 interface. The bus variable controls
// the I2C bus used to communicate with the device.
func New(bus i2c.Bus) LSM303 {
	return &lsm303{bus: bus, mu: new(sync.RWMutex), poll: pollDelay}
}

// Initialize the device
func (d *lsm303) setup() (err error) {
	d.mu.RLock()
	if d.initialized {
		d.mu.RUnlock()
		return
	}
	d.mu.RUnlock()

	d.mu.Lock()
	defer d.mu.Unlock()

	if err = d.bus.WriteToReg(magAddress, magConfigRegA, MagCRADefault); err != nil {
		return
	}
	if err = d.bus.WriteToReg(magAddress, magModeReg, MagMRDefault); err != nil {
		return
	}

	d.initialized = true

	return
}

// SetPollDelay sets the delay between runs of the data acquisition loop.
func (d *lsm303) SetPollDelay(delay int) {
	d.poll = delay
}

func (d *lsm303) measureHeading() (heading float64, err error) {
	if err = d.setup(); err != nil {
		return
	}

	if _, err = d.bus.ReadByteFromReg(magAddress, magDataSignal); err != nil {
		return
	}

	data := make([]byte, 6)
	if err = d.bus.ReadFromReg(magAddress, magData, data); err != nil {
		return
	}

	x := int16(data[0])<<8 | int16(data[1])
	y := int16(data[2])<<8 | int16(data[3])

	heading = math.Atan2(float64(y), float64(x)) / math.Pi * 180
	if heading < 0 {
		heading += 360
	}

	return
}

// Return heading [0, 360).
func (d *lsm303) Heading() (heading float64, err error) {
	select {
	case heading = <-d.headings:
		return
	default:
		if d.debug {
			log.Print("lsm303: no headings available... measuring")
		}
		return d.measureHeading()
	}
}

// Start the sensor data acquisition loop.
func (d *lsm303) Run() (err error) {
	go func() {
		d.quit = make(chan struct{})

		timer := time.Tick(time.Duration(d.poll) * time.Millisecond)

		var heading float64

		for {
			select {
			case <-timer:
				h, err := d.measureHeading()
				if err == nil {
					heading = h
				}
				if err == nil && d.headings == nil {
					d.headings = make(chan float64)
				}
			case d.headings <- heading:
			case <-d.quit:
				d.headings = nil
				return
			}

		}
	}()

	return
}

// Close the sensor data acquisition loop and put the LSM303 into sleep mode.
func (d *lsm303) Close() (err error) {
	if d.quit != nil {
		d.quit <- struct{}{}
	}
	err = d.bus.WriteToReg(magAddress, magModeReg, MagSleep)
	return
}

// Return heading [0, 360).
func Heading() (heading float64, err error) {
	return Default.Heading()
}

// Start the sensor data acquisition loop.
func Run() (err error) {
	return Default.Run()
}

// Close the sensor data acquisition loop and put the LSM303 into sleep mode.
func Close() (err error) {
	return Default.Close()
}
