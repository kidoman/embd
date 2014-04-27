package embd

const (
	spi_cpha = 0x01
	spi_cpol = 0x02

	SPI_MODE_0 = (0 | 0)
	SPI_MODE_1 = (0 | spi_cpha)
	SPI_MODE_2 = (spi_cpol | 0)
	SPI_MODE_3 = (spi_cpol | spi_cpha)
)

type SPIBus interface {
	TransferAndRecieveData(dataBuffer []uint8) error

	ReceiveData(len int) ([]uint8, error)

	TransferAndReceiveByte(data byte) (byte, error)

	ReceiveByte() (byte, error)

	Close() error
}

type SPIDriver interface {
	Bus(byte, byte, int, int, int) SPIBus

	Close() error
}

var spiDriverInitialized bool
var spiDriverInstance SPIDriver

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

func CloseSPI() error {
	return spiDriverInstance.Close()
}

func NewSPIBus(mode, channel byte, speed, bpw, delay int) SPIBus {
	if err := InitSPI(); err != nil {
		panic(err)
	}

	return spiDriverInstance.Bus(mode, channel, speed, bpw, delay)
}
