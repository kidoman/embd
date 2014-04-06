// Generic I²C driver.

package embd

import "sync"

type i2cBusFactory func(byte) I2CBus

type i2cDriver struct {
	busMap     map[byte]I2CBus
	busMapLock sync.Mutex

	ibf i2cBusFactory
}

// NewI2CDriver returns a I2CDriver interface which allows control
// over the I²C subsystem.
func NewI2CDriver(ibf i2cBusFactory) I2CDriver {
	return &i2cDriver{
		busMap: make(map[byte]I2CBus),
		ibf:    ibf,
	}
}

func (i *i2cDriver) Bus(l byte) I2CBus {
	i.busMapLock.Lock()
	defer i.busMapLock.Unlock()

	if b, ok := i.busMap[l]; ok {
		return b
	}

	b := i.ibf(l)
	i.busMap[l] = b
	return b
}

func (i *i2cDriver) Close() error {
	for _, b := range i.busMap {
		b.Close()
	}

	return nil
}
