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

	slaveCmd = 0x0703 // Cmd to set slave address
	rdrwCmd  = 0x0707 // Cmd to read/write data together

	rd = 0x0001
)

// A Bus implements access to the I2C two wire bus.
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

// Default instance of the i2c bus
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
	messages[1].flags = rd
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

func (b *bus) ReadByteFromReg(addr, reg byte) (value byte, err error) {
	buf := make([]byte, 1)
	if err = b.ReadFromReg(addr, reg, buf); err != nil {
		return
	}
	value = buf[0]
	return
}

func (b *bus) ReadWordFromReg(addr, reg byte) (value uint16, err error) {
	buf := make([]byte, 2)
	if err = b.ReadFromReg(addr, reg, buf); err != nil {
		return
	}
	value = uint16((uint16(buf[0]) << 8) | uint16(buf[1]))
	return
}

func (b *bus) WriteToReg(addr, reg byte, value []byte) (err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if err = b.setAddress(addr); err != nil {
		return
	}

	outbuf := append([]byte{reg}, value...)

	hdrp := (*reflect.SliceHeader)(unsafe.Pointer(&outbuf))

	var message i2c_msg
	message.addr = uint16(addr)
	message.flags = 0
	message.len = uint16(len(outbuf))
	message.buf = uintptr(unsafe.Pointer(&hdrp.Data))

	var packets i2c_rdwr_ioctl_data

	packets.msgs = uintptr(unsafe.Pointer(&message))
	packets.nmsg = 1

	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, b.file.Fd(), rdrwCmd, uintptr(unsafe.Pointer(&packets))); errno != 0 {
		err = syscall.Errno(errno)
		return
	}

	return
}

func (b *bus) WriteByteToReg(addr, reg, value byte) (err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if err = b.setAddress(addr); err != nil {
		return
	}

	outbuf := [...]byte{
		reg,
		value,
	}

	var message i2c_msg
	message.addr = uint16(addr)
	message.flags = 0
	message.len = uint16(len(outbuf))
	message.buf = uintptr(unsafe.Pointer(&outbuf))

	var packets i2c_rdwr_ioctl_data

	packets.msgs = uintptr(unsafe.Pointer(&message))
	packets.nmsg = 1

	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, b.file.Fd(), rdrwCmd, uintptr(unsafe.Pointer(&packets))); errno != 0 {
		err = syscall.Errno(errno)
		return
	}

	return
}

func (b *bus) WriteWordToReg(addr, reg byte, value uint16) (err error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if err = b.setAddress(addr); err != nil {
		return
	}

	outbuf := [...]byte{
		reg,
		byte(value >> 8),
		byte(value),
	}

	var messages i2c_msg
	messages.addr = uint16(addr)
	messages.flags = 0
	messages.len = uint16(len(outbuf))
	messages.buf = uintptr(unsafe.Pointer(&outbuf))

	var packets i2c_rdwr_ioctl_data

	packets.msgs = uintptr(unsafe.Pointer(&messages))
	packets.nmsg = 1

	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, b.file.Fd(), rdrwCmd, uintptr(unsafe.Pointer(&packets))); errno != 0 {
		err = syscall.Errno(errno)
		return
	}

	return
}

// ReadByte reads a byte from the given address.
func ReadByte(addr byte) (value byte, err error) {
	return Default.ReadByte(addr)
}

// WriteByte writes a byte to the given address.
func WriteByte(addr, value byte) error {
	return Default.WriteByte(addr, value)
}

// WriteBytes writes a slice bytes to the given address.
func WriteBytes(addr byte, value []byte) error {
	return Default.WriteBytes(addr, value)
}

// ReadFromReg reads n (len(value)) bytes from the given address and register.
func ReadFromReg(addr, reg byte, value []byte) error {
	return Default.ReadFromReg(addr, reg, value)
}

// ReadByteFromReg reads a byte from the given address and register.
func ReadByteFromReg(addr, reg byte) (value byte, err error) {
	return Default.ReadByteFromReg(addr, reg)
}

// ReadU16FromReg reads a unsigned 16 bit integer from the given address and register.
func ReadWordFromReg(addr, reg byte) (value uint16, err error) {
	return Default.ReadWordFromReg(addr, reg)
}

// WriteToReg writes len(value) bytes to the given address and register.
func WriteToReg(addr, reg byte, value []byte) error {
	return Default.WriteToReg(addr, reg, value)
}

// WriteByteToReg writes a byte to the given address and register.
func WriteByteToReg(addr, reg, value byte) error {
	return Default.WriteByteToReg(addr, reg, value)
}

// WriteU16ToReg
func WriteWordToReg(addr, reg byte, value uint16) error {
	return Default.WriteWordToReg(addr, reg, value)
}
