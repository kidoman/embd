// Package mcp4725 allows interfacing with the MCP4725 DAC.
package mcp4725

import (
	"log"
	"sync"

	"github.com/kid0m4n/go-rpi/i2c"
)

const (
	dacReg     = 0x40
	programReg = 0x60
	powerDown  = 0x46

	genReset = 0x06
	powerUp  = 0x09
)

// MCP4725 represents a MCP4725 DAC.
type MCP4725 struct {
	// Bus to communicate over.
	Bus i2c.Bus
	// Addr of the sensor.
	Addr byte

	initialized bool
	mu          sync.RWMutex

	// Debug turns on additional debug output.
	Debug bool
}

// New creates a new MCP4725 sensor.
func New(bus i2c.Bus, addr byte) *MCP4725 {
	return &MCP4725{
		Bus:  bus,
		Addr: addr,
	}
}

func (d *MCP4725) setup() (err error) {
	d.mu.RLock()
	if d.initialized {
		d.mu.RUnlock()
		return
	}
	d.mu.RUnlock()

	d.mu.Lock()
	defer d.mu.Unlock()

	if d.Debug {
		log.Print("mcp4725: general call reset")
	}
	if err = d.Bus.WriteByteToReg(d.Addr, 0x00, powerUp); err != nil {
		return
	}
	if err = d.Bus.WriteByteToReg(d.Addr, 0x00, genReset); err != nil {
		return
	}
	d.initialized = true
	return
}

func (d *MCP4725) setVoltage(voltage int, reg byte) (err error) {
	if err = d.setup(); err != nil {
		return
	}
	if voltage > 4095 {
		voltage = 4095
	}
	if voltage < 0 {
		voltage = 0
	}
	if d.Debug {
		log.Printf("mcp4725: setting voltage to %04d", voltage)
	}
	if err = d.Bus.WriteWordToReg(d.Addr, reg, uint16(voltage<<4)); err != nil {
		return
	}
	return
}

// SetVoltage sets the output voltage.
func (d *MCP4725) SetVoltage(voltage int) (err error) {
	return d.setVoltage(voltage, dacReg)
}

// SetPersistedVoltage sets the voltage and programs the EEPROM so
// that the voltage is restored on reboot.
func (d *MCP4725) SetPersistedVoltage(voltage int) (err error) {
	return d.setVoltage(voltage, programReg)
}

// Close puts the DAC into power down mode.
func (d *MCP4725) Close() (err error) {
	if d.Debug {
		log.Print("mcp4725: powering down")
	}
	if err = d.Bus.WriteWordToReg(d.Addr, powerDown, 0); err != nil {
		return
	}
	return
}
