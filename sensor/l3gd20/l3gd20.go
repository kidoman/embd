// Package l3gd20 allows interacting with L3GD20 gyroscoping sensor.
package l3gd20

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/kidoman/embd"
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

	dr95  = 0x00
	dr190 = 0x40
	dr380 = 0x80
	dr760 = 0xC0

	xEnabled = 0x01
	yEnabled = 0x02
	zEnabled = 0x04

	powerOn  = 0x08
	powerOff = 0x00

	ctrlReg1Default  = powerOn | xEnabled | yEnabled | zEnabled
	ctrlReg1Finished = powerOff | xEnabled | yEnabled | zEnabled

	zyxAvailable = 0x08

	odr       = 95
	mult      = 1.0 / odr
	pollDelay = mult * 1000 * 1000
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

type Orientation struct {
	X, Y, Z float64
}

// L3GD20 represents a L3GD20 3-axis gyroscope.
type L3GD20 struct {
	Bus   embd.I2CBus
	Range *Range

	initialized bool
	mu          sync.RWMutex

	xac, yac, zac axisCalibration

	orientations chan Orientation
	closing      chan chan struct{}
}

// New creates a new L3GD20 interface. The bus variable controls
// the I2C bus used to communicate with the device.
func New(bus embd.I2CBus, Range *Range) *L3GD20 {
	return &L3GD20{
		Bus:   bus,
		Range: Range,
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

func (d *L3GD20) calibrate(a *axis) (axisCalibration, error) {
	glog.V(1).Infof("l3gd20: calibrating %v axis", a)

	values := make(values, 0)
	for i := 0; i < 20; i++ {
	again:
		available, err := d.axisStatus(a)
		if err != nil {
			return axisCalibration{}, err
		}
		if !available {
			time.Sleep(100 * time.Microsecond)
			goto again
		}
		value, err := d.readOrientationDelta(a)
		if err != nil {
			return axisCalibration{}, err
		}
		values = append(values, value)
	}
	ac := axisCalibration{min: values.min(), max: values.max(), mean: values.mean()}

	glog.V(1).Infof("l3gd20: %v axis calibration (%v)", a, ac)

	return ac, nil
}

func (d *L3GD20) setup() error {
	d.mu.RLock()
	if d.initialized {
		d.mu.RUnlock()
		return nil
	}
	d.mu.RUnlock()

	d.mu.Lock()
	defer d.mu.Unlock()

	d.orientations = make(chan Orientation)

	if err := d.Bus.WriteByteToReg(address, ctrlReg1, ctrlReg1Default); err != nil {
		return err
	}
	if err := d.Bus.WriteByteToReg(address, ctrlReg4, d.Range.value); err != nil {
		return err
	}

	// Calibrate
	var err error
	if d.xac, err = d.calibrate(ax); err != nil {
		return err
	}
	if d.yac, err = d.calibrate(ay); err != nil {
		return err
	}
	if d.zac, err = d.calibrate(az); err != nil {
		return err
	}

	d.initialized = true

	return nil
}

func (d *L3GD20) axisStatus(a *axis) (bool, error) {
	data, err := d.Bus.ReadByteFromReg(address, statusReg)
	if err != nil {
		return false, err
	}

	if data&zyxAvailable == 0 {
		return false, nil
	}

	available := data&a.availableMask != 0

	return available, nil
}

func (d *L3GD20) readOrientationDelta(a *axis) (float64, error) {
	rl, rh := a.regs()
	l, err := d.Bus.ReadByteFromReg(address, rl)
	if err != nil {
		return 0, err
	}
	h, err := d.Bus.ReadByteFromReg(address, rh)
	if err != nil {
		return 0, err
	}

	value := float64(int16(h)<<8 | int16(l))
	value *= d.Range.sensitivity

	return value, nil
}

func (d *L3GD20) calibratedOrientationDelta(a *axis) (float64, error) {
	value, err := d.readOrientationDelta(a)
	if err != nil {
		return 0, err
	}

	switch a {
	case ax:
		value = d.xac.adjust(value)
	case ay:
		value = d.yac.adjust(value)
	case az:
		value = d.zac.adjust(value)
	}

	return value, nil
}

func (d *L3GD20) measureOrientationDelta() (float64, float64, float64, error) {
	if err := d.setup(); err != nil {
		return 0, 0, 0, err
	}

	dx, err := d.calibratedOrientationDelta(ax)
	if err != nil {
		return 0, 0, 0, err
	}
	dy, err := d.calibratedOrientationDelta(ay)
	if err != nil {
		return 0, 0, 0, err
	}
	dz, err := d.calibratedOrientationDelta(az)
	if err != nil {
		return 0, 0, 0, err
	}

	return dx, dy, dz, nil
}

// Orientation returns the current orientation reading.
func (d *L3GD20) OrientationDelta() (float64, float64, float64, error) {
	return d.measureOrientationDelta()
}

// Temperature returns the current temperature reading.
func (d *L3GD20) Temperature() (int, error) {
	if err := d.setup(); err != nil {
		return 0, err
	}

	data, err := d.Bus.ReadByteFromReg(address, tempData)
	if err != nil {
		return 0, err
	}

	temp := int(int8(data))

	return temp, nil
}

// Orientations returns a channel which will have the current temperature reading.
func (d *L3GD20) Orientations() (<-chan Orientation, error) {
	if err := d.setup(); err != nil {
		return nil, err
	}

	return d.orientations, nil
}

// Start starts the data acquisition loop.
func (d *L3GD20) Start() error {
	if err := d.setup(); err != nil {
		return err
	}

	d.closing = make(chan chan struct{})

	go func() {
		var x, y, z float64
		var orientations chan Orientation

		timer := time.Tick(time.Duration(math.Floor(pollDelay)) * time.Microsecond)

		for {
			select {
			case <-timer:
				dx, dy, dz, err := d.measureOrientationDelta()
				if err != nil {
					glog.Errorf("l3gd20: %v", err)
				} else {
					x += dx * mult
					y += dy * mult
					z += dz * mult
					orientations = d.orientations
				}
			case orientations <- Orientation{x, y, z}:
				orientations = nil
			case waitc := <-d.closing:
				waitc <- struct{}{}
				close(d.orientations)
				return
			}

		}
	}()

	return nil
}

// Stop the data acquisition loop.
func (d *L3GD20) Stop() error {
	if d.closing != nil {
		waitc := make(chan struct{})
		d.closing <- waitc
		<-waitc
		d.closing = nil
	}
	if err := d.Bus.WriteByteToReg(address, ctrlReg1, ctrlReg1Finished); err != nil {
		return err
	}
	d.initialized = false
	return nil
}

// Close.
func (d *L3GD20) Close() error {
	return d.Stop()
}
