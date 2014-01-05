// Package pca9685 allows interfacing with the pca9685 16-channel, 12-bit PWM Controller through I2C protocol.
package pca9685

import (
	"log"
	"math"
	"sync"
	"time"

	"github.com/kid0m4n/go-rpi/i2c"
)

const (
	clockFreq        = 25000000
	pwmControlPoints = 4096

	mode1RegAddr    = 0x00
	preScaleRegAddr = 0xFE

	pwm0OnLowReg = 0x6
)

// A PCA9685 interface implements access to the controller.
type PCA9685 interface {
	// SetPwm sets the ON and OFF time registers for pwm signal shaping.
	// n: channel 0-15
	// onTime/offTime: 0-4095
	SetPwm(n int, onTime int, offTime int) error

	// Wake allows the controller to exit sleep mode and resume with PWM generation.
	Wake() error

	// Sleep puts the controller in sleep mode. Does not change the pwm control registers.
	Sleep() error

	// Close stops the controller and resets mode and pwm controller registers.
	Close() error

	// SetDebug is used to enable logging (debug mode).
	SetDebug(status bool)
}

type pca9685 struct {
	bus  i2c.Bus
	addr byte
	freq int

	initialized bool
	mu          sync.RWMutex

	debug bool
}

// The pca9685 controller supports pwm frequencies from 40Hz to 1000Hz

// New creates a new PCA9685 interface.
func New(bus i2c.Bus, addr byte, freq int) PCA9685 {
	return &pca9685{bus: bus, addr: addr, freq: freq}
}

// SetDebug is used to enable logging (debug mode).
func (d *pca9685) SetDebug(status bool) {
	d.debug = status
}

func (d *pca9685) mode1Reg() (byte, error) {
	return d.bus.ReadByteFromReg(d.addr, mode1RegAddr)
}

func (d *pca9685) setup() (err error) {
	d.mu.RLock()
	if d.initialized {
		d.mu.RUnlock()
		return
	}
	d.mu.RUnlock()

	d.mu.Lock()
	defer d.mu.Unlock()

	preScaleValue := byte(math.Floor(float64(clockFreq/(pwmControlPoints*d.freq))+float64(0.5)) - 1)
	if d.debug {
		log.Printf("pca9685: calculated Prescale value = %#02x", preScaleValue)
	}

	mode1Reg, err := d.mode1Reg()
	if err != nil {
		return
	}
	if d.debug {
		log.Printf("pca9685: read MODE1 Reg [regAddr: %#02x] Value: [%v]", mode1RegAddr, mode1Reg)
	}

	if err = d.sleep(); err != nil {
		return
	}

	if err = d.bus.WriteByteToReg(d.addr, preScaleRegAddr, byte(preScaleValue)); err != nil {
		return
	}
	if d.debug {
		log.Printf("pca9685: prescale value [%#02x] written to PRE_SCALE Reg [regAddr: %#02x]", preScaleValue, preScaleRegAddr)
	}

	if err = d.wake(); err != nil {
		return
	}

	newmode := ((mode1Reg | 0x01) & 0xDF)
	if err = d.bus.WriteByteToReg(d.addr, mode1RegAddr, newmode); err != nil {
		return
	}
	if d.debug {
		log.Printf("pca9685: new mode [%#02x] [disabling register auto increment] written to MODE1 Reg [regAddr: %#02x]", newmode, mode1RegAddr)
	}

	d.initialized = true
	if d.debug {
		log.Printf("pca9685: driver initialized with pwm freq: %v", d.freq)
	}

	return
}

// SetPwm sets the ON and OFF time registers for pwm signal shaping.
// channel: 0-15
// onTime/offTime: 0-4095
func (d *pca9685) SetPwm(n, onTime, offTime int) (err error) {
	if err = d.setup(); err != nil {
		return
	}

	onTimeLowReg := byte(pwm0OnLowReg + (4 * n))

	onTimeLow := byte(onTime & 0xFF)
	onTimeHigh := byte(onTime >> 8)
	offTimeLow := byte(offTime & 0xFF)
	offTimeHigh := byte(offTime >> 8)

	if err = d.bus.WriteByteToReg(d.addr, onTimeLowReg, onTimeLow); err != nil {
		return
	}
	if d.debug {
		log.Printf("pca9685: writing On-Time Low [%#02x] to CHAN%v_ON_L reg [RegAddr: %#02x]", onTimeLow, n, onTimeLowReg)
	}

	onTimeHighReg := onTimeLowReg + 1
	if err = d.bus.WriteByteToReg(d.addr, onTimeHighReg, onTimeHighReg); err != nil {
		return
	}
	if d.debug {
		log.Printf("pca9685: writing On-Time High [%#02x] to CHAN%v_ON_H reg [RegAddr: %#02x]", onTimeHigh, n, onTimeHighReg)
	}

	offTimeLowReg := onTimeHighReg + 1
	if err = d.bus.WriteByteToReg(d.addr, offTimeLowReg, offTimeLow); err != nil {
		return
	}
	if d.debug {
		log.Printf("pca9685: writing Off-Time Low [%#02x] to CHAN%v_OFF_L reg [RegAddr: %#02x]", offTimeLow, n, offTimeLowReg)
	}

	offTimeHighReg := offTimeLowReg + 1
	if err = d.bus.WriteByteToReg(d.addr, offTimeHighReg, offTimeHigh); err != nil {
		return
	}
	if d.debug {
		log.Printf("pca9685: writing Off-Time High [%#02x] to CHAN%v_OFF_H reg [RegAddr: %#02x]", offTimeHigh, n, offTimeHighReg)
	}

	return
}

// Close stops the controller and resets mode and pwm controller registers.
func (d *pca9685) Close() (err error) {
	if err = d.setup(); err != nil {
		return
	}

	if err = d.sleep(); err != nil {
		return
	}

	if d.debug {
		log.Println("pca9685: reset request received")
	}

	if err = d.bus.WriteByteToReg(d.addr, mode1RegAddr, 0x00); err != nil {
		return
	}

	if d.debug {
		log.Printf("pca9685: cleaning up all PWM control registers")
	}

	for regAddr := 0x0; regAddr <= 0x45; regAddr++ {
		if err = d.bus.WriteByteToReg(d.addr, byte(regAddr), 0x00); err != nil {
			return
		}
	}

	if d.debug {
		log.Printf("pca9685: done Cleaning up all PWM control registers")
	}

	if d.debug {
		log.Println("pca9685: controller reset")
	}

	return
}

func (d *pca9685) sleep() (err error) {
	if d.debug {
		log.Println("pca9685: sleep request received")
	}

	mode1Reg, err := d.mode1Reg()
	if err != nil {
		return
	}
	sleepmode := (mode1Reg & 0x7F) | 0x10
	if err = d.bus.WriteByteToReg(d.addr, mode1RegAddr, sleepmode); err != nil {
		return
	}
	if d.debug {
		log.Printf("pca9685: sleep mode [%#02x] written to MODE1 Reg [regAddr: %#02x]", sleepmode, mode1RegAddr)
	}

	if d.debug {
		log.Println("pca9685: controller set to Sleep mode")
	}

	return
}

// Sleep puts the controller in sleep mode. Does not change the pwm control registers.
func (d *pca9685) Sleep() (err error) {
	if err = d.setup(); err != nil {
		return
	}

	return d.sleep()
}

func (d *pca9685) wake() (err error) {
	if d.debug {
		log.Println("pca9685: wake request received")
	}

	mode1Reg, err := d.mode1Reg()
	if err != nil {
		return
	}
	wakeMode := mode1Reg & 0xEF
	if (mode1Reg & 0x80) == 0x80 {
		if err = d.bus.WriteByteToReg(d.addr, mode1RegAddr, wakeMode); err != nil {
			return
		}
		if d.debug {
			log.Printf("pca9685: wake mode [%#02x] written to MODE1 Reg [regAddr: %#02x]", wakeMode, mode1RegAddr)
		}

		time.Sleep(500 * time.Microsecond)
	}

	restartOpCode := wakeMode | 0x80
	if err = d.bus.WriteByteToReg(d.addr, mode1RegAddr, restartOpCode); err != nil {
		return
	}
	if d.debug {
		log.Printf("pca9685: restart mode [%#02x] written to MODE1 Reg [regAddr: %#02x]", restartOpCode, mode1RegAddr)
	}

	return
}

// Wake allows the controller to exit sleep mode and resume with PWM generation.
func (d *pca9685) Wake() (err error) {
	if err = d.setup(); err != nil {
		return
	}

	return d.wake()
}
