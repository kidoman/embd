// I2C support.

package embd

// I2CBus interface is used to interact with the I2C bus.
type I2CBus interface {
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

// I2CDriver interface interacts with the host descriptors to allow us
// control of I2C communication.
type I2CDriver interface {
	Bus(l byte) I2CBus

	Close() error
}

var i2cDriverInstance I2CDriver

// InitI2C initializes the I2C driver.
func InitI2C() error {
	desc, err := DescribeHost()
	if err != nil {
		return err
	}

	if desc.I2CDriver == nil {
		return ErrFeatureNotSupported
	}

	i2cDriverInstance = desc.I2CDriver()

	return nil
}

// CloseI2C releases resources associated with the I2C driver.
func CloseI2C() error {
	return i2cDriverInstance.Close()
}

// NewI2CBus returns a I2CBus.
func NewI2CBus(l byte) I2CBus {
	return i2cDriverInstance.Bus(l)
}
