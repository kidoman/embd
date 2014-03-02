// Package pca9685 allows interfacing with the pca9685 16-channel, 12-bit PWM Controller through I2C protocol.
package pca9685

import (
	"log"
	"math"
	"sync"
	"time"
	"github.com/kidoman/embd"
	"github.com/kidoman/embd/util"
)

const (
	clockFreq        = 25000000
	pwmControlPoints = 4096

	mode1RegAddr    = 0x00
	preScaleRegAddr = 0xFE

	pwm0OnLowReg = 0x6

	minAnalogValue = 0
	maxAnalogValue = 255

	defaultFreq = 490
)

// PCA9685 represents a PCA9685 PWM generator.
type PCA9685 struct {
	Bus  embd.I2CBus
	Addr byte
	Freq int

	initialized bool
	mu          sync.RWMutex

	Debug bool
}

// New creates a new PCA9685 interface.
func New(bus embd.I2CBus, addr byte) *PCA9685 {
	return &PCA9685{
		Bus:  bus,
		Addr: addr,
	}
}

func (d *PCA9685) mode1Reg() (byte, error) {
	return d.Bus.ReadByteFromReg(d.Addr, mode1RegAddr)
}

func (d *PCA9685) setup() (err error) {
	d.mu.RLock()
	if d.initialized {
		d.mu.RUnlock()
		return
	}
	d.mu.RUnlock()

	d.mu.Lock()
	defer d.mu.Unlock()

	mode1Reg, err := d.mode1Reg()
	if err != nil {
		return
	}
	if d.Debug {
		log.Printf("pca9685: read MODE1 Reg [regAddr: %#02x] Value: [%v]", mode1RegAddr, mode1Reg)
	}

	if err = d.sleep(); err != nil {
		return
	}

	if d.Freq == 0 {
		d.Freq = defaultFreq
	}
	preScaleValue := byte(math.Floor(float64(clockFreq/(pwmControlPoints*d.Freq))+float64(0.5)) - 1)
	if d.Debug {
		log.Printf("pca9685: calculated prescale value = %#02x", preScaleValue)
	}
	if err = d.Bus.WriteByteToReg(d.Addr, preScaleRegAddr, byte(preScaleValue)); err != nil {
		return
	}
	if d.Debug {
		log.Printf("pca9685: prescale value [%#02x] written to PRE_SCALE Reg [regAddr: %#02x]", preScaleValue, preScaleRegAddr)
	}

	if err = d.wake(); err != nil {
		return
	}

	newmode := ((mode1Reg | 0x01) & 0xDF)
	if err = d.Bus.WriteByteToReg(d.Addr, mode1RegAddr, newmode); err != nil {
		return
	}
	if d.Debug {
		log.Printf("pca9685: new mode [%#02x] [disabling register auto increment] written to MODE1 Reg [regAddr: %#02x]", newmode, mode1RegAddr)
	}

	d.initialized = true
	if d.Debug {
		log.Printf("pca9685: driver initialized with pwm freq: %v", d.Freq)
	}

	return
}

// SetPwm sets the ON and OFF time registers for pwm signal shaping.
// channel: 0-15
// onTime/offTime: 0-4095
func (d *PCA9685) SetPwm(channel, onTime, offTime int) (err error) {
	if err = d.setup(); err != nil {
		return
	}

	onTimeLowReg := byte(pwm0OnLowReg + (4 * channel))

	onTimeLow := byte(onTime & 0xFF)
	onTimeHigh := byte(onTime >> 8)
	offTimeLow := byte(offTime & 0xFF)
	offTimeHigh := byte(offTime >> 8)

	if err = d.Bus.WriteByteToReg(d.Addr, onTimeLowReg, onTimeLow); err != nil {
		return
	}
	if d.Debug {
		log.Printf("pca9685: writing on-time low [%#02x] to CHAN%v_ON_L reg [reg: %#02x]", onTimeLow, channel, onTimeLowReg)
	}

	onTimeHighReg := onTimeLowReg + 1
	if err = d.Bus.WriteByteToReg(d.Addr, onTimeHighReg, onTimeHigh); err != nil {
		return
	}
	if d.Debug {
		log.Printf("pca9685: writing on-time high [%#02x] to CHAN%v_ON_H reg [reg: %#02x]", onTimeHigh, channel, onTimeHighReg)
	}

	offTimeLowReg := onTimeHighReg + 1
	if err = d.Bus.WriteByteToReg(d.Addr, offTimeLowReg, offTimeLow); err != nil {
		return
	}
	if d.Debug {
		log.Printf("pca9685: writing off-time low [%#02x] to CHAN%v_OFF_L reg [reg: %#02x]", offTimeLow, channel, offTimeLowReg)
	}

	offTimeHighReg := offTimeLowReg + 1
	if err = d.Bus.WriteByteToReg(d.Addr, offTimeHighReg, offTimeHigh); err != nil {
		return
	}
	if d.Debug {
		log.Printf("pca9685: writing off-time high [%#02x] to CHAN%v_OFF_H reg [reg: %#02x]", offTimeHigh, channel, offTimeHighReg)
	}

	return
}

func (d *PCA9685) SetMicroseconds(channel, us int) (err error) {
	offTime := us * d.Freq * pwmControlPoints / 1000000
	return d.SetPwm(channel, 0, offTime)
}

func (d *PCA9685) SetAnalog(channel int, value byte) (err error) {
	offTime := util.Map(int64(value), minAnalogValue, maxAnalogValue, 0, pwmControlPoints-1)
	return d.SetPwm(channel, 0, int(offTime))
}

// Close stops the controller and resets mode and pwm controller registers.
func (d *PCA9685) Close() (err error) {
	if err = d.setup(); err != nil {
		return
	}

	if err = d.sleep(); err != nil {
		return
	}

	if d.Debug {
		log.Println("pca9685: reset request received")
	}

	if err = d.Bus.WriteByteToReg(d.Addr, mode1RegAddr, 0x00); err != nil {
		return
	}

	if d.Debug {
		log.Printf("pca9685: cleaning up all PWM control registers")
	}

	for regAddr := 0x06; regAddr <= 0x45; regAddr++ {
		if err = d.Bus.WriteByteToReg(d.Addr, byte(regAddr), 0x00); err != nil {
			return
		}
	}

	if d.Debug {
		log.Printf("pca9685: done Cleaning up all PWM control registers")
	}

	if d.Debug {
		log.Println("pca9685: controller reset")
	}

	return
}

func (d *PCA9685) sleep() (err error) {
	if d.Debug {
		log.Println("pca9685: sleep request received")
	}

	mode1Reg, err := d.mode1Reg()
	if err != nil {
		return
	}
	sleepmode := (mode1Reg & 0x7F) | 0x10
	if err = d.Bus.WriteByteToReg(d.Addr, mode1RegAddr, sleepmode); err != nil {
		return
	}
	if d.Debug {
		log.Printf("pca9685: sleep mode [%#02x] written to MODE1 Reg [regAddr: %#02x]", sleepmode, mode1RegAddr)
	}

	if d.Debug {
		log.Println("pca9685: controller set to Sleep mode")
	}

	return
}

// Sleep puts the controller in sleep mode. Does not change the pwm control registers.
func (d *PCA9685) Sleep() (err error) {
	if err = d.setup(); err != nil {
		return
	}

	return d.sleep()
}

func (d *PCA9685) wake() (err error) {
	if d.Debug {
		log.Println("pca9685: wake request received")
	}

	mode1Reg, err := d.mode1Reg()
	if err != nil {
		return
	}
	wakeMode := mode1Reg & 0xEF
	if (mode1Reg & 0x80) == 0x80 {
		if err = d.Bus.WriteByteToReg(d.Addr, mode1RegAddr, wakeMode); err != nil {
			return
		}
		if d.Debug {
			log.Printf("pca9685: wake mode [%#02x] written to MODE1 Reg [regAddr: %#02x]", wakeMode, mode1RegAddr)
		}

		time.Sleep(500 * time.Microsecond)
	}

	restartOpCode := wakeMode | 0x80
	if err = d.Bus.WriteByteToReg(d.Addr, mode1RegAddr, restartOpCode); err != nil {
		return
	}
	if d.Debug {
		log.Printf("pca9685: restart mode [%#02x] written to MODE1 Reg [regAddr: %#02x]", restartOpCode, mode1RegAddr)
	}

	return
}

// Wake allows the controller to exit sleep mode and resume with PWM generation.
func (d *PCA9685) Wake() (err error) {
	if err = d.setup(); err != nil {
		return
	}

	return d.wake()
}
