// Package i2c enables gophers i2c speaking ability.
package i2c

import (
	"fmt"
	"os"
	"reflect"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

const (
	delay = 20

	slaveCmd = 0x0703
	rdrwCmd  = 0x0707

	I2C_M_RD = 0x0001
)

type Bus interface {
	ReadByte(addr byte) (value byte, err error)
	WriteByte(addr, value byte) error
	WriteBytes(addr byte, value []byte) error

	ReadFromReg(addr, reg byte, value []byte) (err error)
	ReadByteFromReg(addr, reg byte) (value byte, err error)
	ReadInt(addr, reg byte) (value int, err error)

	WriteToReg(addr, reg, value byte) error
}

type i2c_msg struct {
	addr  uint16
	flags uint16
	len   uint16
	buf   uintptr
}

type i2c_rdwr_ioctl_data struct {
	msgs uintptr
	nmsg uint32
}

var busMap map[byte]*bus
var busMapLock sync.Mutex
var Default Bus

type bus struct {
	file *os.File
	addr byte
	mu   sync.Mutex
}

func init() {
	busMap = make(map[byte]*bus)
	var err error
	Default, err = NewBus(1)
	if err != nil {
		panic(err)
	}
}

// NewBus creates a new I2C bus interface. The l variable
// controls which bus we connect to.
//
// For the newer RaspberryPi, the number is 1 (earlier model uses 0.)
func NewBus(l byte) (Bus, error) {
	busMapLock.Lock()
	defer busMapLock.Unlock()

	var b *bus

	if b = busMap[l]; b == nil {
		b = new(bus)
		var err error
		if b.file, err = os.OpenFile(fmt.Sprintf("/dev/i2c-%v", l), os.O_RDWR, os.ModeExclusive); err != nil {
			return nil, err
		}
		busMap[l] = b
	}

	return b, nil
}

func (b *bus) setAddress(addr byte) (err error) {
	if addr != b.addr {
		if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, b.file.Fd(), slaveCmd, uintptr(addr)); errno != 0 {
			err = syscall.Errno(errno)
			return
		}

		b.addr = addr
	}

	return
}

// Read a byte from the given address.
func (b *bus) ReadByte(addr byte) (value byte, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if err = b.setAddress(addr); err != nil {
		return
	}

	bytes := make([]byte, 1)
	n, err := b.file.Read(bytes)

	if n != 1 {
		err = fmt.Errorf("i2c: Unexpected number (%v) of bytes read", n)
	}

	value = bytes[0]

	return
}

// Write a byte to the given address.
func (b *bus) WriteByte(addr, value byte) (err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if err = b.setAddress(addr); err != nil {
		return
	}

	n, err := b.file.Write([]byte{value})

	if n != 1 {
		err = fmt.Errorf("i2c: Unexpected number (%v) of bytes written in WriteByte", n)
	}

	return
}

// Write a bunch of bytes ot the given address.
func (b *bus) WriteBytes(addr byte, value []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if err := b.setAddress(addr); err != nil {
		return err
	}

	for i := range value {
		n, err := b.file.Write([]byte{value[i]})

		if n != 1 {
			return fmt.Errorf("i2c: Unexpected number (%v) of bytes written in WriteBytes", n)
		}
		if err != nil {
			return err
		}

		time.Sleep(delay * time.Millisecond)
	}

	return nil
}

// Read a bunch of bytes (len(value)) from the given address and register.
func (b *bus) ReadFromReg(addr, reg byte, value []byte) (err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if err = b.setAddress(addr); err != nil {
		return
	}

	hdrp := (*reflect.SliceHeader)(unsafe.Pointer(&value))

	var messages [2]i2c_msg
	messages[0].addr = uint16(addr)
	messages[0].flags = 0
	messages[0].len = 1
	messages[0].buf = uintptr(unsafe.Pointer(&reg))

	messages[1].addr = uint16(addr)
	messages[1].flags = I2C_M_RD
	messages[1].len = uint16(len(value))
	messages[1].buf = uintptr(unsafe.Pointer(hdrp.Data))

	var packets i2c_rdwr_ioctl_data

	packets.msgs = uintptr(unsafe.Pointer(&messages))
	packets.nmsg = 2

	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, b.file.Fd(), rdrwCmd, uintptr(unsafe.Pointer(&packets))); errno != 0 {
		return syscall.Errno(errno)
	}

	return nil
}

// Read a byte from the given address and register.
func (b *bus) ReadByteFromReg(addr, reg byte) (value byte, err error) {
	buf := make([]byte, 1)
	if err = b.ReadFromReg(addr, reg, buf); err != nil {
		return
	}
	value = buf[0]
	return
}

// Read a int from the given address and register.
func (b *bus) ReadInt(addr, reg byte) (value int, err error) {
	var buf = make([]byte, 2)
	if err = b.ReadFromReg(addr, reg, buf); err != nil {
		return
	}
	value = int((int(buf[0]) << 8) | int(buf[1]))
	return
}

// Write a byte to the given address and register.
func (b *bus) WriteToReg(addr, reg, value byte) (err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if err = b.setAddress(addr); err != nil {
		return
	}

	var outbuf [2]byte
	var messages i2c_msg
	messages.addr = uint16(addr)
	messages.flags = 0
	messages.len = uint16(len(outbuf))
	messages.buf = uintptr(unsafe.Pointer(&outbuf))

	outbuf[0] = reg
	outbuf[1] = value

	var packets i2c_rdwr_ioctl_data

	packets.msgs = uintptr(unsafe.Pointer(&messages))
	packets.nmsg = 1

	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, b.file.Fd(), rdrwCmd, uintptr(unsafe.Pointer(&packets))); errno != 0 {
		err = syscall.Errno(errno)
		return
	}

	return
}

// Read a byte from the given address.
func ReadByte(addr byte) (value byte, err error) {
	return Default.ReadByte(addr)
}

// Write a byte to the given address.
func WriteByte(addr, value byte) (err error) {
	return Default.WriteByte(addr, value)
}

// Write a bunch of bytes ot the given address.
func WriteBytes(addr byte, value []byte) error {
	return Default.WriteBytes(addr, value)
}

// Read a bunch of bytes (len(value)) from the given address and register.
func ReadFromReg(addr, reg byte, value []byte) (err error) {
	return Default.ReadFromReg(addr, reg, value)
}

// Read a byte from the given address and register.
func ReadByteFromReg(addr, reg byte) (value byte, err error) {
	return Default.ReadByteFromReg(addr, reg)
}

// Read a int from the given address and register.
func ReadInt(addr, reg byte) (value int, err error) {
	return Default.ReadInt(addr, reg)
}

// Write a byte to the given address and register.
func WriteToReg(addr, reg, value byte) (err error) {
	return Default.WriteToReg(addr, reg, value)
}
