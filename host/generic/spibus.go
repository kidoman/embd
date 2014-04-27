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
	SPI_IOC_WR_MODE          = 0x40016B01
	SPI_IOC_WR_BITS_PER_WORD = 0x40016B03
	SPI_IOC_WR_MAX_SPEED_HZ  = 0x40046B04

	SPI_IOC_RD_MODE          = 0x80016B01
	SPI_IOC_RD_BITS_PER_WORD = 0x80016B03
	SPI_IOC_RD_MAX_SPEED_HZ  = 0x80046B04

	SPI_IOC_MESSAGE_0   = 1073769216 //0x40006B00
	SPI_IOC_INCREMENTER = 2097152    //0x200000

	DEFAULT_DELAYMS   = uint16(0)
	DEFAULT_SPI_BPW   = uint8(8)
	DEFAULT_SPI_SPEED = uint32(1000000)
)

type spiIocTransfer struct {
	tx_buf uint64
	rx_buf uint64

	length        uint32
	speed_hz      uint32
	delay_usecs   uint16
	bits_per_word uint8
}

type spiBus struct {
	file *os.File

	spiDevMinor byte

	channel byte
	mode    byte
	speed   int
	bpw     int
	delayms int

	mu sync.Mutex

	spiTransferData spiIocTransfer
	initialized     bool
}

func spi_ioc_message_n(n uint32) uint32 {
	return (SPI_IOC_MESSAGE_0 + (n * SPI_IOC_INCREMENTER))
}

func NewSPIBus(spiDevMinor, mode, channel byte, speed, bpw, delay int) embd.SPIBus {
	return &spiBus{
		spiDevMinor: spiDevMinor,
		mode:        mode,
		channel:     channel,
		speed:       speed,
		bpw:         bpw,
		delayms:     delay,
	}
}

func (b *spiBus) init() error {
	if b.initialized {
		return nil
	}

	var err error
	if b.file, err = os.OpenFile(fmt.Sprintf("/dev/spidev%v.%v", b.spiDevMinor, b.channel), os.O_RDWR, os.ModeExclusive); err != nil {
		return err
	}

	err = b.setMode()
	if err != nil {
		return err
	}

	b.spiTransferData = spiIocTransfer{}

	err = b.setSpeed()
	if err != nil {
		return err
	}

	err = b.setBpw()
	if err != nil {
		return err
	}

	b.setDelay()

	glog.V(2).Infof("spi: bus %v initialized", b.channel)

	b.initialized = true
	return nil

}

func (b *spiBus) setMode() error {
	var mode = uint8(b.mode)
	var err error
	glog.V(3).Infof("spi: setting spi mode to %v", mode)

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, b.file.Fd(), SPI_IOC_WR_MODE, uintptr(unsafe.Pointer(&mode)))
	if errno != 0 {
		err = syscall.Errno(errno)
		glog.V(3).Infof("spi: failed to set mode due to %v", err.Error())
		return err
	}
	glog.V(3).Infof("spi: mode set to %v", mode)
	return nil
}

func (b *spiBus) setSpeed() error {
	var speed uint32
	if b.speed > 0 {
		speed = uint32(b.speed)
	} else {
		speed = DEFAULT_SPI_SPEED
	}

	glog.V(3).Infof("spi: setting spi speed_max to %v", speed)
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, b.file.Fd(), SPI_IOC_WR_MAX_SPEED_HZ, uintptr(unsafe.Pointer(&speed)))
	if errno != 0 {
		err := syscall.Errno(errno)
		glog.V(3).Infof("spi: failed to set speed_max due to %v", err.Error())
		return err
	}
	glog.V(3).Infof("spi: speed_max set to %v", speed)
	b.spiTransferData.speed_hz = speed

	return nil
}

func (b *spiBus) setBpw() error {
	var bpw uint8

	if b.bpw > 0 {
		bpw = uint8(b.bpw)
	} else {
		bpw = DEFAULT_SPI_BPW
	}

	glog.V(3).Infof("spi: setting spi bpw to %v", bpw)
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, b.file.Fd(), SPI_IOC_WR_BITS_PER_WORD, uintptr(unsafe.Pointer(&bpw)))
	if errno != 0 {
		err := syscall.Errno(errno)
		glog.V(3).Infof("spi: failed to set bpw due to %v", err.Error())
		return err
	}
	glog.V(3).Infof("spi: bpw set to %v", bpw)
	b.spiTransferData.bits_per_word = uint8(bpw)
	return nil
}

func (b *spiBus) setDelay() {
	var delay uint16

	if b.delayms > 0 {
		delay = uint16(b.delayms)
	} else {
		delay = DEFAULT_DELAYMS
	}
	glog.V(3).Infof("spi: delay_ms set to %v", delay)
	b.spiTransferData.delay_usecs = delay
}

func (b *spiBus) TransferAndRecieveData(dataBuffer []uint8) error {
	len := len(dataBuffer)
	dataCarrier := b.spiTransferData

	dataCarrier.length = uint32(len)
	dataCarrier.tx_buf = uint64(uintptr(unsafe.Pointer(&dataBuffer[0])))
	dataCarrier.rx_buf = uint64(uintptr(unsafe.Pointer(&dataBuffer[0])))

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, b.file.Fd(), uintptr(spi_ioc_message_n(1)), uintptr(unsafe.Pointer(&dataCarrier)))
	if errno != 0 {
		err := syscall.Errno(errno)
		glog.V(3).Infof("spi: failed to read due to %v", err.Error())
		return err
	}
	return nil
}

func (b *spiBus) ReceiveData(len int) ([]uint8, error) {
	data := make([]uint8, len)
	var err error
	err = b.TransferAndRecieveData(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (b *spiBus) TransferAndReceiveByte(data byte) (byte, error) {
	d := make([]uint8, 1)
	d[0] = uint8(data)
	err := b.TransferAndRecieveData(d)
	if err != nil {
		return 0, err
	}
	return d[0], nil
}

func (b *spiBus) ReceiveByte() (byte, error) {
	d := make([]uint8, 1)
	err := b.TransferAndRecieveData(d)
	if err != nil {
		return 0, err
	}
	return byte(d[0]), nil
}

func (b *spiBus) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.initialized {
		return nil
	}

	return b.file.Close()
}
