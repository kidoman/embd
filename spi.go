package embd

const (
	spiCpha = 0x01
	spiCpol = 0x02

	SpiMode0 = (0 | 0)
	SpiMode1 = (0 | spiCpha)
	SpiMode2 = (spiCpol | 0)
	SpiMode3 = (spiCpol | spiCpha)
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
