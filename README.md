# go-rpi

Use various sensors on the RaspberryPi with Golang (like a ninja!)

## Documentation

[![GoDoc](http://godoc.org/github.com/kid0m4n/go-rpi?status.png)](http://godoc.org/github.com/kid0m4n/go-rpi)

## Protocols supported

* I2C [Documentation](http://godoc.org/github.com/kid0m4n/go-rpi/i2c)

## Sensors supported

* TMP006 Thermopile sensor [Documentation](http://godoc.org/github.com/kid0m4n/go-rpi/sensor/tmp006), [Datasheet](http://www.adafruit.com/datasheets/tmp006.pdf), [Userguide](http://www.adafruit.com/datasheets/tmp006ug.pdf)

* BMP085 Barometric pressure sensor [Documentation](http://godoc.org/github.com/kid0m4n/go-rpi/sensor/bmp085), [Datasheet](https://www.sparkfun.com/datasheets/Components/General/BST-BMP085-DS000-05.pdf)

* BMP180 Barometric pressure sensor [Documentation](http://godoc.org/github.com/kid0m4n/go-rpi/sensor/bmp180), [Datasheet](http://www.adafruit.com/datasheets/BST-BMP180-DS000-09.pdf)

* LSM303 Accelerometer and magnetometer [Documentation](http://godoc.org/github.com/kid0m4n/go-rpi/sensor/lsm303), [Datasheet](https://www.sparkfun.com/datasheets/Sensors/Magneto/LSM303%20Datasheet.pdf)

* L3GD20 Gyroscope [Documentation](http://godoc.org/github.com/kid0m4n/go-rpi/sensor/l3gd20), [Datasheet](http://www.adafruit.com/datasheets/L3GD20.pdf)

* US020 Ultrasonic proximity sensor [Documentation](http://godoc.org/github.com/kid0m4n/go-rpi/sensor/us020), [Product Page](http://www.digibay.in/sensor/object-detection-and-proximity?product_id=239)

* BH1750FVI Luminosity sensor [Documentation](http://godoc.org/github.com/kid0m4n/go-rpi/sensor/us020), [Datasheet](http://www.elechouse.com/elechouse/images/product/Digital%20light%20Sensor/bh1750fvi-e.pdf)

## Interfaces

* Keypad(4x3) [Product Page](http://www.adafruit.com/products/419#Learn)

## Controllers

* PCA9685 16-channel, 12-bit PWM Controller with I2C protocol [Documentation](http://godoc.org/github.com/kid0m4n/go-rpi/controller/pca9685), [Datasheet](http://www.adafruit.com/datasheets/PCA9685.pdf), [Product Page](http://www.adafruit.com/products/815)

* MCP4725 12-bit DAC [Documentation](http://godoc.org/github.com/kid0m4n/go-rpi/controller/mcp4725), [Datasheet](http://www.adafruit.com/datasheets/mcp4725.pdf), [Product Page](http://www.adafruit.com/products/935)
