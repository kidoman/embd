// Package lsm303 allows interfacing with the LSM303 magnetometer.
package lsm303

import (
	"math"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/kidoman/embd"
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

// LSM303 represents a LSM303 magnetometer.
type LSM303 struct {
	Bus  embd.I2CBus
	Poll int

	initialized bool
	mu          sync.RWMutex

	headings chan float64

	quit chan struct{}
}

// New creates a new LSM303 interface. The bus variable controls
// the I2C bus used to communicate with the device.
func New(bus embd.I2CBus) *LSM303 {
	return &LSM303{Bus: bus, Poll: pollDelay}
}

// Initialize the device
func (d *LSM303) setup() error {
	d.mu.RLock()
	if d.initialized {
		d.mu.RUnlock()
		return nil
	}
	d.mu.RUnlock()

	d.mu.Lock()
	defer d.mu.Unlock()

	if err := d.Bus.WriteByteToReg(magAddress, magConfigRegA, MagCRADefault); err != nil {
		return err
	}
	if err := d.Bus.WriteByteToReg(magAddress, magModeReg, MagMRDefault); err != nil {
		return err
	}

	d.initialized = true

	return nil
}

func (d *LSM303) measureHeading() (float64, error) {
	if err := d.setup(); err != nil {
		return 0, err
	}

	if _, err := d.Bus.ReadByteFromReg(magAddress, magDataSignal); err != nil {
		return 0, err
	}

	data := make([]byte, 6)
	if err := d.Bus.ReadFromReg(magAddress, magData, data); err != nil {
		return 0, err
	}

	x := int16(data[0])<<8 | int16(data[1])
	y := int16(data[2])<<8 | int16(data[3])

	heading := math.Atan2(float64(y), float64(x)) / math.Pi * 180
	if heading < 0 {
		heading += 360
	}

	return heading, nil
}

// Heading returns the current heading [0, 360).
func (d *LSM303) Heading() (float64, error) {
	select {
	case heading := <-d.headings:
		return heading, nil
	default:
		glog.V(2).Infof("lsm303: no headings available... measuring")
		return d.measureHeading()
	}
}

// Run starts the sensor data acquisition loop.
func (d *LSM303) Run() error {
	go func() {
		d.quit = make(chan struct{})

		timer := time.Tick(time.Duration(d.Poll) * time.Millisecond)

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

	return nil
}

// Close the sensor data acquisition loop and put the LSM303 into sleep mode.
func (d *LSM303) Close() error {
	if d.quit != nil {
		d.quit <- struct{}{}
	}
	return d.Bus.WriteByteToReg(magAddress, magModeReg, MagSleep)
}
