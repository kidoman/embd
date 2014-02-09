// Package lsm303 allows interfacing with the LSM303 magnetometer.
package lsm303

import (
	"log"
	"math"
	"sync"
	"time"

	"github.com/kidoman/embd/i2c"
)

const (
	magAddress = 0x1E

	magConfigRegA = 0x00

	MagHz75         = 0x00 // ODR = 0.75 Hz
	Mag1Hz5         = 0x04 // ODR = 1.5 Hz
	Mag3Hz          = 0x08 // ODR = 3 Hz
	Mag7Hz5         = 0x0C // ODR = 7.5 Hz
	Mag15Hz         = 0x10 // ODR = 15 Hz
	Mag30Hz         = 0x14 // ODR = 30 Hz
	Mag75Hz         = 0x18 // ODR = 75 Hz
	MagNormal       = 0x00 // Normal mode
	MagPositiveBias = 0x01 // Positive bias mode
	MagNegativeBias = 0x02 // Negative bias mode

	MagCRADefault = Mag15Hz | MagNormal // 15 Hz and normal mode is the default

	magModeReg = 0x02

	MagContinuous = 0x00 // Continuous conversion mode
	MagSingle     = 0x01 // Single conversion mode
	MagSleep      = 0x03 // Sleep mode

	MagMRDefault = MagContinuous // Continuous conversion is the default

	magDataSignal = 0x02
	magData       = 0x03

	pollDelay = 250
)

// A LSM303 implements access to a LSM303 sensor.
type LSM303 interface {
	// SetPollDelay sets the delay between runs of the data acquisition loop.
	SetPollDelay(delay int)

	// Heading returns the current heading [0, 360).
	Heading() (heading float64, err error)

	// Run starts the sensor data acquisition loop.
	Run() error
	// Close closes the sensor data acquisition loop and puts the LSM303 into sleep mode.
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

// Default instance of the LSM303 sensor.
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

	if err = d.bus.WriteByteToReg(magAddress, magConfigRegA, MagCRADefault); err != nil {
		return
	}
	if err = d.bus.WriteByteToReg(magAddress, magModeReg, MagMRDefault); err != nil {
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

// Heading returns the current heading [0, 360).
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

// Run starts the sensor data acquisition loop.
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
	err = d.bus.WriteByteToReg(magAddress, magModeReg, MagSleep)
	return
}

// SetPollDelay sets the delay between runs of the data acquisition loop.
func SetPollDelay(delay int) {
	Default.SetPollDelay(delay)
}

// Heading returns the current heading [0, 360).
func Heading() (heading float64, err error) {
	return Default.Heading()
}

// Run starts the sensor data acquisition loop.
func Run() (err error) {
	return Default.Run()
}

// Close closes the sensor data acquisition loop and put the LSM303 into sleep mode.
func Close() (err error) {
	return Default.Close()
}
