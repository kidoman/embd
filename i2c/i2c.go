// Package i2c enables gophers i2c speaking ability.
package i2c

import "github.com/kidoman/embd/host"

type Bus interface {
	// ReadByte reads a byte from the given address.
	ReadByte(addr byte) (value byte, err error)
	// WriteByte writes a byte to the given address.
	WriteByte(addr, value byte) error
	// WriteBytes writes a slice bytes to the given address.
	WriteBytes(addr byte, value []byte) error

	// ReadFromReg reads n (len(value)) bytes from the given address and register.
	ReadFromReg(addr, reg byte, value []byte) error
	// ReadByteFromReg reads a byte from the given address and register.
	ReadByteFromReg(addr, reg byte) (value byte, err error)
	// ReadU16FromReg reads a unsigned 16 bit integer from the given address and register.
	ReadWordFromReg(addr, reg byte) (value uint16, err error)

	// WriteToReg writes len(value) bytes to the given address and register.
	WriteToReg(addr, reg byte, value []byte) error
	// WriteByteToReg writes a byte to the given address and register.
	WriteByteToReg(addr, reg, value byte) error
	// WriteU16ToReg
	WriteWordToReg(addr, reg byte, value uint16) error
}

type I2C interface {
	Bus(l byte) Bus

	Close() error
}

var instance I2C

func Open() error {
	desc, err := host.Describe()
	if err != nil {
		return err
	}

	instance = desc.I2C().(I2C)

	return nil
}

func Close() error {
	return instance.Close()
}

func NewBus(l byte) Bus {
	return instance.Bus(l)
}
