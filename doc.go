/*
Package rpi provides modules which will help gophers deal with various sensors.

Use the default i2c bus to read/write data:

	value, err := i2c.ReadInt(0x1E, 0x03)
	...
	value := make([]byte, 6)
	err := i2c.ReadFromReg(0x77, 0xF6, value)
	...
	err := i2c.WriteToReg(0x1E, 0x02, 0x00)

Read data from the BMP085 sensor:

	temp, err := bmp085.Temperature()
	...
	altitude, err := bmp085.Altitude()

Find out the heading from the LSM303 magnetometer:

	heading, err := lsm303.Heading()
*/
package rpi
