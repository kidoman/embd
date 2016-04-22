// SPI support.

package embd

import (
	"io"
)

const (
	spiCpha = 0x01
	spiCpol = 0x02

	// SPIMode0 represents the mode0 operation (CPOL=0 CPHA=0) of spi.
	SPIMode0 = (0 | 0)

	// SPIMode1 represents the mode0 operation (CPOL=0 CPHA=1) of spi.
	SPIMode1 = (0 | spiCpha)

	// SPIMode2 represents the mode0 operation (CPOL=1 CPHA=0) of spi.
	SPIMode2 = (spiCpol | 0)

	// SPIMode3 represents the mode0 operation (CPOL=1 CPHA=1) of spi.
	SPIMode3 = (spiCpol | spiCpha)
)

// SPIBus interface allows interaction with the SPI bus.
type SPIBus interface {
	io.Writer

	// TransferAndReceiveData transmits data in a buffer(slice) and receives into it.
	TransferAndReceiveData(dataBuffer []uint8) error

	// ReceiveData receives data of length len into a slice.
	ReceiveData(len int) ([]uint8, error)

	// TransferAndReceiveByte transmits a byte data and receives a byte.
	TransferAndReceiveByte(data byte) (byte, error)

	// ReceiveByte receives a byte data.
	ReceiveByte() (byte, error)

	// Close releases the resources associated with the bus.
	Close() error
}

// SPIDriver interface interacts with the host descriptors to allow us
// control of SPI communication.
type SPIDriver interface {
	// Bus returns a SPIBus interface which allows us to use spi functionalities
	Bus(byte, byte, int, int, int) SPIBus

	// Close cleans up all the initialized SPIbus
	Close() error
}

var spiDriverInitialized bool
var spiDriverInstance SPIDriver

// InitSPI initializes the SPI driver.
func InitSPI() error {
	if spiDriverInitialized {
		return nil
	}

	desc, err := DescribeHost()
	if err != nil {
		return err
	}

	if desc.SPIDriver == nil {
		return ErrFeatureNotSupported
	}

	spiDriverInstance = desc.SPIDriver()
	spiDriverInitialized = true

	return nil
}

// CloseSPI releases resources associated with the SPI driver.
func CloseSPI() error {
	return spiDriverInstance.Close()
}

// NewSPIBus returns a SPIBus.
func NewSPIBus(mode, channel byte, speed, bpw, delay int) SPIBus {
	if err := InitSPI(); err != nil {
		panic(err)
	}

	return spiDriverInstance.Bus(mode, channel, speed, bpw, delay)
}
