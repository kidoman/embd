package generic

import (
	"fmt"
	"os"
	"sync"
	"syscall"
	"unsafe"

	"github.com/golang/glog"
	"github.com/kidoman/embd"
)

const (
	spiIOCWrMode        = 0x40016B01
	spiIOCWrBitsPerWord = 0x40016B03
	spiIOCWrMaxSpeedHz  = 0x40046B04

	spiIOCRdMode        = 0x80016B01
	spiIOCRdBitsPerWord = 0x80016B03
	spiIOCRdMaxSpeedHz  = 0x80046B04

	spiIOCMessage0    = 1073769216 //0x40006B00
	spiIOCIncrementor = 2097152    //0x200000

	defaultDelayms  = 0
	defaultSPIBPW   = 8
	defaultSPISpeed = 1000000
)

type spiIOCTransfer struct {
	txBuf uint64
	rxBuf uint64

	length      uint32
	speedHz     uint32
	delayus     uint16
	bitsPerWord uint8
	csChange    uint8
	pad         uint32
}

type spiBus struct {
	file *os.File

	spiDevMinor int

	channel byte
	mode    byte
	speed   int
	bpw     int
	delayms int

	mu sync.Mutex

	spiTransferData spiIOCTransfer
	initialized     bool

	initializer func() error
}

func spiIOCMessageN(n uint32) uint32 {
	return (spiIOCMessage0 + (n * spiIOCIncrementor))
}

func NewSPIBus(spiDevMinor int, mode, channel byte, speed, bpw, delay int, i func() error) embd.SPIBus {
	return &spiBus{
		spiDevMinor: spiDevMinor,
		mode:        mode,
		channel:     channel,
		speed:       speed,
		bpw:         bpw,
		delayms:     delay,
		initializer: i,
	}
}

func (b *spiBus) init() error {
	if b.initialized {
		return nil
	}

	if b.initializer != nil {
		if err := b.initializer(); err != nil {
			return err
		}
	}

	var err error
	if b.file, err = os.OpenFile(fmt.Sprintf("/dev/spidev%v.%v", b.spiDevMinor, b.channel), os.O_RDWR, os.ModeExclusive); err != nil {
		return err
	}
	glog.V(3).Infof("spi: sucessfully opened file /dev/spidev%v.%v", b.spiDevMinor, b.channel)

	if err = b.setMode(); err != nil {
		return err
	}

	b.spiTransferData = spiIOCTransfer{}

	if err = b.setSpeed(); err != nil {
		return err
	}

	if err = b.setBPW(); err != nil {
		return err
	}

	b.setDelay()

	glog.V(2).Infof("spi: bus %v initialized", b.channel)
	glog.V(3).Infof("spi: bus %v initialized with spiIOCTransfer as %v", b.channel, b.spiTransferData)

	b.initialized = true
	return nil
}

func (b *spiBus) setMode() error {
	var mode = uint8(b.mode)
	glog.V(3).Infof("spi: setting spi mode to %v", mode)

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, b.file.Fd(), spiIOCWrMode, uintptr(unsafe.Pointer(&mode)))
	if errno != 0 {
		err := syscall.Errno(errno)
		glog.V(3).Infof("spi: failed to set mode due to %v", err.Error())
		return err
	}
	glog.V(3).Infof("spi: mode set to %v", mode)
	return nil
}

func (b *spiBus) setSpeed() error {
	var speed uint32 = defaultSPISpeed
	if b.speed > 0 {
		speed = uint32(b.speed)
	}

	glog.V(3).Infof("spi: setting spi speedMax to %v", speed)
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, b.file.Fd(), spiIOCWrMaxSpeedHz, uintptr(unsafe.Pointer(&speed)))
	if errno != 0 {
		err := syscall.Errno(errno)
		glog.V(3).Infof("spi: failed to set speedMax due to %v", err.Error())
		return err
	}
	glog.V(3).Infof("spi: speedMax set to %v", speed)
	b.spiTransferData.speedHz = speed

	return nil
}

func (b *spiBus) setBPW() error {
	var bpw uint8 = defaultSPIBPW
	if b.bpw > 0 {
		bpw = uint8(b.bpw)
	}

	glog.V(3).Infof("spi: setting spi bpw to %v", bpw)
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, b.file.Fd(), spiIOCWrBitsPerWord, uintptr(unsafe.Pointer(&bpw)))
	if errno != 0 {
		err := syscall.Errno(errno)
		glog.V(3).Infof("spi: failed to set bpw due to %v", err.Error())
		return err
	}
	glog.V(3).Infof("spi: bpw set to %v", bpw)
	b.spiTransferData.bitsPerWord = uint8(bpw)
	return nil
}

func (b *spiBus) setDelay() {
	var delay uint16 = defaultDelayms
	if b.delayms > 0 {
		delay = uint16(b.delayms)
	}

	glog.V(3).Infof("spi: delayms set to %v", delay)
	b.spiTransferData.delayus = delay
}

func (b *spiBus) TransferAndReceiveData(dataBuffer []uint8) error {
	if err := b.init(); err != nil {
		return err
	}

	len := len(dataBuffer)
	dataCarrier := b.spiTransferData

	dataCarrier.length = uint32(len)
	dataCarrier.txBuf = uint64(uintptr(unsafe.Pointer(&dataBuffer[0])))
	dataCarrier.rxBuf = uint64(uintptr(unsafe.Pointer(&dataBuffer[0])))

	glog.V(3).Infof("spi: sending dataBuffer %v with carrier %v", dataBuffer, dataCarrier)
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, b.file.Fd(), uintptr(spiIOCMessageN(1)), uintptr(unsafe.Pointer(&dataCarrier)))
	if errno != 0 {
		err := syscall.Errno(errno)
		glog.V(3).Infof("spi: failed to read due to %v", err.Error())
		return err
	}
	glog.V(3).Infof("spi: read into dataBuffer %v", dataBuffer)
	return nil
}

func (b *spiBus) ReceiveData(len int) ([]uint8, error) {
	if err := b.init(); err != nil {
		return nil, err
	}

	data := make([]uint8, len)
	if err := b.TransferAndReceiveData(data); err != nil {
		return nil, err
	}
	return data, nil
}

func (b *spiBus) TransferAndReceiveByte(data byte) (byte, error) {
	if err := b.init(); err != nil {
		return 0, err
	}

	d := [1]uint8{uint8(data)}
	if err := b.TransferAndReceiveData(d[:]); err != nil {
		return 0, err
	}
	return d[0], nil
}

func (b *spiBus) ReceiveByte() (byte, error) {
	if err := b.init(); err != nil {
		return 0, err
	}

	var d [1]uint8
	if err := b.TransferAndReceiveData(d[:]); err != nil {
		return 0, err
	}
	return byte(d[0]), nil
}

func (b *spiBus) Write(data []byte) (n int, err error) {
	if err := b.init(); err != nil {
		return 0, err
	}
	return b.file.Write(data)
}

func (b *spiBus) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.initialized {
		return nil
	}

	return b.file.Close()
}
