// Package tmp006 allows interfacing with the TMP006 thermopile.
package tmp006

import (
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/golang/glog"
	"github.com/kidoman/embd"
)

const (
	b0   = -0.0000294
	b1   = -0.00000057
	b2   = 0.00000000463
	c2   = 13.4
	tref = 298.15
	a2   = -0.00001678
	a1   = 0.00175
	s0   = 6.4

	vObjReg    = 0x00
	tempAmbReg = 0x01
	configReg  = 0x02
	manIdReg   = 0xFE
	manId      = 0x5449
	devIdReg   = 0xFF
	devId      = 0x0067

	reset  = 0x8000
	modeOn = 0x7000
	drdyEn = 0x0100

	configRegDefault = modeOn | drdyEn
)

type SampleRate struct {
	enabler      uint16
	samples      int
	timeRequired float64
}

var (
	SR1  = &SampleRate{0x0000, 1, 0.25} // 1 sample, 0.25 second between measurements.
	SR2  = &SampleRate{0x0200, 2, 0.5}  // 2 samples, 0.5 second between measurements.
	SR4  = &SampleRate{0x0400, 4, 1}    // 4 samples, 1 second between measurements.
	SR8  = &SampleRate{0x0600, 8, 2}    // 8 samples, 2 seconds between measurements.
	SR16 = &SampleRate{0x0800, 16, 4}   // 16 samples, 4 seconds between measurements.
)

// TMP006 represents a TMP006 thermopile sensor.
type TMP006 struct {
	// Bus to communicate over.
	Bus embd.I2CBus
	// Addr of the sensor.
	Addr byte
	// SampleRate specifies the sampling rate for the sensor.
	SampleRate *SampleRate

	initialized bool
	mu          sync.RWMutex

	rawDieTemps chan float64
	objTemps    chan float64
	closing     chan chan struct{}
}

// New creates a new TMP006 sensor.
func New(bus embd.I2CBus, addr byte) *TMP006 {
	return &TMP006{
		Bus:  bus,
		Addr: addr,
	}
}

func (d *TMP006) validate() error {
	if d.Bus == nil {
		return errors.New("tmp006: bus is nil")
	}
	if d.Addr == 0x00 {
		return fmt.Errorf("tmp006: %#x is not a valid address", d.Addr)
	}
	return nil
}

// Close puts the device into low power mode.
func (d *TMP006) Close() error {
	if err := d.setup(); err != nil {
		return err
	}
	if d.closing != nil {
		waitc := make(chan struct{})
		d.closing <- waitc
		<-waitc
	}
	glog.V(1).Infof("tmp006: resetting")
	if err := d.Bus.WriteWordToReg(d.Addr, configReg, reset); err != nil {
		return err
	}
	return nil
}

// Present checks if the device is present at the given address.
func (d *TMP006) Present() (bool, error) {
	if err := d.validate(); err != nil {
		return false, err
	}
	mid, err := d.Bus.ReadWordFromReg(d.Addr, manIdReg)
	if err != nil {
		return false, err
	}
	glog.V(1).Infof("tmp006: got manufacturer id %#04x", mid)
	if mid != manId {
		return false, fmt.Errorf("tmp006: not found at %#02x, manufacturer id mismatch", d.Addr)
	}
	did, err := d.Bus.ReadWordFromReg(d.Addr, devIdReg)
	if err != nil {
		return false, err
	}
	glog.V(1).Infof("tmp006: got device id %#04x", did)
	if did != devId {
		return false, fmt.Errorf("tmp006: not found at %#02x, device id mismatch", d.Addr)
	}

	return true, nil
}

func (d *TMP006) setup() error {
	d.mu.RLock()
	if d.initialized {
		d.mu.RUnlock()
		return nil
	}
	d.mu.RUnlock()

	d.mu.Lock()
	defer d.mu.Unlock()

	if err := d.validate(); err != nil {
		return err
	}
	if d.SampleRate == nil {
		glog.V(1).Infof("tmp006: sample rate = nil, using SR16")
		d.SampleRate = SR16
	}
	glog.V(1).Infof("tmp006: configuring with %#04x", configRegDefault|d.SampleRate.enabler)
	if err := d.Bus.WriteWordToReg(d.Addr, configReg, configRegDefault|d.SampleRate.enabler); err != nil {
		return err
	}

	d.initialized = true

	return nil
}

func (d *TMP006) measureRawDieTemp() (float64, error) {
	if err := d.setup(); err != nil {
		return 0, err
	}
	raw, err := d.Bus.ReadWordFromReg(d.Addr, tempAmbReg)
	if err != nil {
		return 0, err
	}
	raw >>= 2
	glog.V(2).Infof("tmp006: raw die temp %#04x", raw)

	temp := float64(int16(raw)) * 0.03125

	return temp, nil
}

func (d *TMP006) measureRawVoltage() (int16, error) {
	if err := d.setup(); err != nil {
		return 0, err
	}
	vlt, err := d.Bus.ReadWordFromReg(d.Addr, vObjReg)
	if err != nil {
		return 0, err
	}
	volt := int16(vlt)
	glog.V(2).Infof("tmp006: raw voltage %#04x", volt)
	return volt, nil
}

func (d *TMP006) measureObjTemp() (float64, error) {
	if err := d.setup(); err != nil {
		return 0, err
	}
	tDie, err := d.measureRawDieTemp()
	if err != nil {
		return 0, err
	}
	glog.V(2).Infof("tmp006: tdie = %.2f C", tDie)
	tDie += 273.15 // Convert to K
	vo, err := d.measureRawVoltage()
	if err != nil {
		return 0, err
	}
	vObj := float64(vo)
	vObj *= 156.25 // 156.25 nV per LSB
	vObj /= 1000   // nV -> uV
	glog.V(2).Infof("tmp006: vObj = %.5f uV", vObj)
	vObj /= 1000 // uV -> mV
	vObj /= 1000 // mV -> V

	tdie_tref := tDie - tref
	s := 1 + a1*tdie_tref + a2*tdie_tref*tdie_tref
	s *= s0
	s /= 10000000
	s /= 10000000

	Vos := b0 + b1*tdie_tref + b2*tdie_tref*tdie_tref
	fVobj := (vObj - Vos) + c2*(vObj-Vos)*(vObj-Vos)

	temp := math.Sqrt(math.Sqrt(tDie*tDie*tDie*tDie + fVobj/s))
	temp -= 273.15

	return temp, nil
}

// RawDieTemp returns the current raw die temp reading.
func (d *TMP006) RawDieTemp() (float64, error) {
	select {
	case temp := <-d.rawDieTemps:
		return temp, nil
	default:
		return d.measureRawDieTemp()
	}
}

// RawDieTemps returns a channel to get future raw die temps from.
func (d *TMP006) RawDieTemps() <-chan float64 {
	return d.rawDieTemps
}

// ObjTemp returns the current obj temp reading.
func (d *TMP006) ObjTemp() (float64, error) {
	select {
	case temp := <-d.objTemps:
		return temp, nil
	default:
		return d.measureObjTemp()
	}
}

// ObjTemps returns a channel to fetch obj temps from.
func (d *TMP006) ObjTemps() <-chan float64 {
	return d.objTemps
}

// Start starts the data acquisition loop.
func (d *TMP006) Start() error {
	if err := d.setup(); err != nil {
		return err
	}

	d.rawDieTemps = make(chan float64)
	d.objTemps = make(chan float64)

	go func() {
		var rawDieTemp, objTemp float64
		var rdtAvlb, otAvlb bool
		var err error
		var timer <-chan time.Time
		resetTimer := func() {
			timer = time.After(time.Duration(d.SampleRate.timeRequired*1000) * time.Millisecond)
		}
		resetTimer()

		for {
			var rawDieTemps, objTemps chan float64

			if rdtAvlb {
				rawDieTemps = d.rawDieTemps
			}
			if otAvlb {
				objTemps = d.objTemps
			}

			select {
			case <-timer:
				var rdt float64
				if rdt, err = d.measureRawDieTemp(); err != nil {
					glog.Errorf("tmp006: %v", err)
				} else {
					rawDieTemp = rdt
					rdtAvlb = true
				}
				var ot float64
				if ot, err = d.measureObjTemp(); err != nil {
					glog.Errorf("tmp006: %v", err)
				} else {
					objTemp = ot
					otAvlb = true
				}
				resetTimer()
			case rawDieTemps <- rawDieTemp:
				rdtAvlb = false
			case objTemps <- objTemp:
				otAvlb = false
			case waitc := <-d.closing:
				waitc <- struct{}{}
				close(d.rawDieTemps)
				close(d.objTemps)
				return
			}
		}
	}()

	return nil
}
