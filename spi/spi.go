package spi

import (
	"fmt"
	"log"
	"os"
	"sync"
	"syscall"
	"unsafe"
)

const (
	SPI_CPHA = 0x01
	SPI_CPOL = 0x02

	SPI_MODE_0 = (0 | 0)
	SPI_MODE_1 = (0 | SPI_CPHA)
	SPI_MODE_2 = (SPI_CPOL | 0)
	SPI_MODE_3 = (SPI_CPOL | SPI_CPHA)

	SPI_IOC_WR_MODE          = 0x40016B01
	SPI_IOC_WR_BITS_PER_WORD = 0x40016B03
	SPI_IOC_WR_MAX_SPEED_HZ  = 0x40046B04

	SPI_IOC_RD_MODE          = 0x80016B01
	SPI_IOC_RD_BITS_PER_WORD = 0x80016B03
	SPI_IOC_RD_MAX_SPEED_HZ  = 0x80046B04

	SPI_IOC_MESSAGE_0   = 1073769216 //0x40006B00
	SPI_IOC_INCREMENTER = 2097152    //0x200000
)

type SpiBus interface {
	TransferAndRecieveByteData(byte) (uint8, error)
}

type spiBus struct {
	file *os.File
	mode byte
	mu   sync.Mutex
}

type spiMode struct {
	mode uintptr
}

type spiIocTransfer struct {
	tx_buf uint64
	rx_buf uint64

	length        uint32
	speed_hz      uint32
	delay_usecs   uint16
	bits_per_word uint8
}

func spi_ioc_message_n(n uint32) uint32 {
	return (SPI_IOC_MESSAGE_0 + (n * SPI_IOC_INCREMENTER))
}

func NewSpiBus() (SpiBus, error) {
	var b *spiBus
	var err error
	b = new(spiBus)

	fmt.Println("Opening SPI device.")
	b.file, err = os.OpenFile("/dev/spidev0.1", os.O_EXCL, os.ModeExclusive)

	if err != nil {
		fmt.Print("spi: Could not open SPI device due to ")
		fmt.Println(err.Error())
		return nil, err
	} else {
		fmt.Println("spi: Successfully initialized spiBus")
		fmt.Println("spi: File Fd = ", b.file.Fd())
	}

	var mode uint8
	mode = 0
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, b.file.Fd(), SPI_IOC_RD_MODE, uintptr(unsafe.Pointer(&mode)))

	if errno != 0 {
		fmt.Println(syscall.Errno(errno))
		err = syscall.Errno(errno)
		return nil, err
	} else {
		fmt.Println("Successfully read the mode")
		fmt.Println(mode)
	}

	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL, b.file.Fd(), SPI_IOC_WR_MODE, uintptr(unsafe.Pointer(&mode)))

	if errno != 0 {
		fmt.Println(syscall.Errno(errno))
		err = syscall.Errno(errno)
		return nil, err
	} else {
		fmt.Println("Successfully read the mode")
		fmt.Println(mode)
	}

	var speed_max uint32
	speed_max = 5000000
	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL, b.file.Fd(), SPI_IOC_RD_MAX_SPEED_HZ, uintptr(unsafe.Pointer(&speed_max)))

	if errno != 0 {
		fmt.Println(syscall.Errno(errno))
		err = syscall.Errno(errno)
		return nil, err
	} else {
		fmt.Println("Successfully read the speed")
		fmt.Println(speed_max)
	}

	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL, b.file.Fd(), SPI_IOC_WR_MAX_SPEED_HZ, uintptr(unsafe.Pointer(&speed_max)))

	if errno != 0 {
		fmt.Println(syscall.Errno(errno))
		err = syscall.Errno(errno)
		return nil, err
	} else {
		fmt.Println("Successfully read the speed")
		fmt.Println(speed_max)
	}

	var bps uint32
	bps = 8
	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL, b.file.Fd(), SPI_IOC_RD_BITS_PER_WORD, uintptr(unsafe.Pointer(&bps)))

	if errno != 0 {
		fmt.Println(syscall.Errno(errno))
		err = syscall.Errno(errno)
		return nil, err
	} else {
		fmt.Println("Successfully read the bps")
		fmt.Println(bps)
	}

	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL, b.file.Fd(), SPI_IOC_WR_BITS_PER_WORD, uintptr(unsafe.Pointer(&bps)))

	if errno != 0 {
		fmt.Println(syscall.Errno(errno))
		err = syscall.Errno(errno)
		return nil, err
	} else {
		fmt.Println("Successfully wrote the bps")
		fmt.Println(bps)
	}

	var bpw uint32
	bpw = 8
	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL, b.file.Fd(), SPI_IOC_RD_BITS_PER_WORD, uintptr(unsafe.Pointer(&bpw)))

	if errno != 0 {
		fmt.Println(syscall.Errno(errno))
		err = syscall.Errno(errno)
		return nil, err
	} else {
		fmt.Println("Successfully read the bpw")
		fmt.Println(bpw)
	}

	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL, b.file.Fd(), SPI_IOC_WR_BITS_PER_WORD, uintptr(unsafe.Pointer(&bpw)))

	if errno != 0 {
		fmt.Println(syscall.Errno(errno))
		err = syscall.Errno(errno)
		return nil, err
	} else {
		fmt.Println("Successfully read the bpw")
		fmt.Println(bpw)
	}

	return b, err
}

func (b *spiBus) TransferAndRecieveByteData(tx_data uint8) (rx_data uint8, err error) {
	var data spiIocTransfer

	data.delay_usecs = 0
	data.length = 8
	data.speed_hz = 5000000
	data.bits_per_word = 8
	data.tx_buf = uint64(uintptr(unsafe.Pointer(&tx_data)))
	data.rx_buf = uint64(uintptr(unsafe.Pointer(&rx_data)))

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, b.file.Fd(), uintptr(spi_ioc_message_n(1)), uintptr(unsafe.Pointer(&data)))
	if errno != 0 {
		err = syscall.Errno(errno)
		return 0, nil
	}

	log.Println("Successfully read the data")
	return rx_data, nil
}
