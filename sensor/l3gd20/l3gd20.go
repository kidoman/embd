// Package l3gd20 allows interacting with L3GD20 gyroscoping sensor.
package l3gd20

import (
	"fmt"
	"log"
	"math"
	"sync"
	"time"

	"github.com/kid0m4n/go-rpi/i2c"
)

const (
	address = 0x6B
	id      = 0xD4

	dpsToRps = 0.017453293

	whoAmI    = 0x0F
	ctrlReg1  = 0x20
	ctrlReg2  = 0x21
	ctrlReg3  = 0x22
	ctrlReg4  = 0x23
	ctrlReg5  = 0x24
	tempData  = 0x26
	statusReg = 0x27

	xlReg = 0x28
	xhReg = 0x29
	ylReg = 0x2A
	yhReg = 0x2B
	zlReg = 0x2C
	zhReg = 0x2D

	xEnabled  = 0x01
	xDisabled = 0x00
	yEnabled  = 0x02
	yDisabled = 0x00
	zEnabled  = 0x04
	zDisabled = 0x00

	powerOn   = 0x08
	powerDown = 0x00

	ctrlReg1Default  = xEnabled | yEnabled | zEnabled | powerOn
	ctrlReg1Finished = xDisabled | yDisabled | zDisabled | powerDown

	zyxAvailable = 0x08

	pollDelay = 100
)

// Range represents a L3GD20 range setting.
type Range struct {
	sensitivity float64

	value byte
}

// The three range settings supported by L3GD20.
var (
	R250DPS  = &Range{sensitivity: 0.00875, value: 0x00}
	R500DPS  = &Range{sensitivity: 0.0175, value: 0x10}
	R2000DPS = &Range{sensitivity: 0.070, value: 0x20}
)

type axis struct {
	name string

	lowReg, highReg byte

	availableMask byte
}

func (a *axis) regs() (byte, byte) {
	return a.lowReg, a.highReg
}

func (a axis) String() string {
	return a.name
}

var (
	ax = &axis{name: "X", lowReg: xlReg, highReg: xhReg, availableMask: 0x01}
	ay = &axis{name: "Y", lowReg: ylReg, highReg: yhReg, availableMask: 0x02}
	az = &axis{name: "Z", lowReg: zlReg, highReg: zhReg, availableMask: 0x04}
)

// A L3GD20 implements access to the L3GD20 sensor.
type L3GD20 interface {
	// Orientation returns the current orientation reading.
	Orientation() (x, y, z float64, err error)
	// Temperature returns the current temperature reading.
	Temperature() (temp int, err error)

	// Close.
	Close() error
}

type axisCalibration struct {
	min, max, mean float64
}

func (ac axisCalibration) adjust(value float64) float64 {
	if value >= ac.min && value <= ac.max {
		return 0
	}
	return value - ac.mean
}

func (ac axisCalibration) String() string {
	return fmt.Sprintf("%v, %v, %v", ac.min, ac.max, ac.mean)
}

type data struct {
	x, y, z float64
}

type l3gd20 struct {
	bus i2c.Bus
	rng *Range

	poll int

	initialized bool
	mu          sync.RWMutex

	xac, yac, zac axisCalibration

	orientations chan data
	quit         chan struct{}

	debug bool
}

// New creates a new L3GD20 interface. The bus variable controls
// the I2C bus used to communicate with the device.
func New(bus i2c.Bus, rng *Range) L3GD20 {
	return &l3gd20{
		bus:   bus,
		rng:   rng,
		poll:  pollDelay,
		debug: false,
	}
}

type values []float64

func (vs values) min() float64 {
	value := math.MaxFloat64
	for _, v := range vs {
		value = math.Min(value, v)
	}
	return value
}

func (vs values) max() float64 {
	value := -math.MaxFloat64
	for _, v := range vs {
		value = math.Max(value, v)
	}
	return value
}

func (vs values) mean() float64 {
	sum := 0.0
	for _, v := range vs {
		sum += v
	}
	return sum / float64(len(vs))
}

func (d *l3gd20) calibrate(a *axis) (ac axisCalibration, err error) {
	if d.debug {
		log.Printf("l3gd20: calibrating %v axis", a)
	}

	values := make(values, 0)
	for i := 0; i < 20; i++ {
	again:
		var available bool
		if available, err = d.axisStatus(a); err != nil {
			return
		}
		if !available {
			time.Sleep(100 * time.Microsecond)
			goto again
		}
		var value float64
		if value, err = d.readOrientation(a); err != nil {
			return
		}
		values = append(values, value)
	}
	ac.min, ac.max, ac.mean = values.min(), values.max(), values.mean()

	if d.debug {
		log.Printf("l3gd20: %v axis calibration (%v)", a, ac)
	}

	return
}

func (d *l3gd20) setup() (err error) {
	d.mu.RLock()
	if d.initialized {
		d.mu.RUnlock()
		return
	}
	d.mu.RUnlock()

	d.mu.Lock()
	defer d.mu.Unlock()

	if err = d.bus.WriteToReg(address, ctrlReg1, ctrlReg1Default); err != nil {
		return
	}
	if err = d.bus.WriteToReg(address, ctrlReg4, d.rng.value); err != nil {
		return
	}

	// Calibrate
	if d.xac, err = d.calibrate(ax); err != nil {
		return
	}
	if d.yac, err = d.calibrate(ay); err != nil {
		return
	}
	if d.zac, err = d.calibrate(az); err != nil {
		return
	}

	d.initialized = true

	return
}

func (d *l3gd20) SetPollDelay(delay int) {
	d.poll = delay
}

func (d *l3gd20) axisStatus(a *axis) (available bool, err error) {
	var data byte
	if data, err = d.bus.ReadByteFromReg(address, statusReg); err != nil {
		return
	}

	if data&zyxAvailable == 0 {
		return
	}

	available = data&a.availableMask != 0

	return
}

func (d *l3gd20) readOrientation(a *axis) (value float64, err error) {
	rl, rh := a.regs()
	var l, h byte
	if l, err = d.bus.ReadByteFromReg(address, rl); err != nil {
		return
	}
	if h, err = d.bus.ReadByteFromReg(address, rh); err != nil {
		return
	}

	value = float64(int16(h)<<8 | int16(l))

	sensitivity := d.rng.sensitivity
	value *= sensitivity

	return
}

func (d *l3gd20) calibratedOrientation(a *axis) (value float64, err error) {
	if value, err = d.readOrientation(a); err != nil {
		return
	}

	switch a {
	case ax:
		value = d.xac.adjust(value)
	case ay:
		value = d.yac.adjust(value)
	case az:
		value = d.zac.adjust(value)
	}

	return
}

func (d *l3gd20) measureOrientation() (x, y, z float64, err error) {
	if err = d.setup(); err != nil {
		return
	}

	if x, err = d.calibratedOrientation(ax); err != nil {
		return
	}
	if y, err = d.calibratedOrientation(ay); err != nil {
		return
	}
	if z, err = d.calibratedOrientation(az); err != nil {
		return
	}

	return
}

func (d *l3gd20) Orientation() (x, y, z float64, err error) {
	select {
	case data := <-d.orientations:
		x, y, z = data.x, data.y, data.z
		return
	default:
		if d.debug {
			log.Printf("l3gd20: no orientation available... measuring")
		}
		return d.measureOrientation()
	}

	panic("cannot reach here")
}

func (d *l3gd20) Temperature() (temp int, err error) {
	if err = d.setup(); err != nil {
		return
	}

	var data byte
	if data, err = d.bus.ReadByteFromReg(address, tempData); err != nil {
		return
	}

	temp = int(int8(data))

	return
}

func (d *l3gd20) Run() (err error) {
	go func() {
		d.quit = make(chan struct{})

		timer := time.Tick(time.Duration(d.poll) * time.Millisecond)

		var dt data

		for {
			select {
			case <-timer:
				var err error
				dt.x, dt.y, dt.z, err = d.measureOrientation()
				if err == nil && d.orientations == nil {
					d.orientations = make(chan data)
				}
			case d.orientations <- dt:
			case <-d.quit:
				d.orientations = nil
				return
			}

		}
	}()

	return
}

func (d *l3gd20) Close() (err error) {
	if d.quit != nil {
		d.quit <- struct{}{}
	}
	return d.bus.WriteToReg(address, ctrlReg1, ctrlReg1Finished)
}
