# embd [![Build Status](https://travis-ci.org/kidoman/embd.svg?branch=master)](https://travis-ci.org/kidoman/embd) [![GoDoc](http://godoc.org/github.com/kidoman/embd?status.png)](http://godoc.org/github.com/kidoman/embd)

A superheroic hardware abstraction layer for doing embedded programming on supported platforms like the Raspberry Pi and BeagleBone Black.

Development sponsored by [**ThoughtWorks**](http://www.thoughtworks.com/)

## Platforms supported

* [RaspberryPi](http://www.raspberrypi.org/)
* [BeagleBone Black](http://beagleboard.org/Products/BeagleBone%20Black)
* [Intel Galileo](http://www.intel.com/content/www/us/en/do-it-yourself/galileo-maker-quark-board.html) **coming soon**
* [Radxa](http://radxa.com/) **coming soon**
* [Cubietruck](http://www.cubietruck.com/) **coming soon**
* Bring Your Own **coming soon**

## How to use

Package embd provides a superheroic hardware abstraction layer for doing embedded programming
on supported platforms like the Raspberry Pi and BeagleBone Black. Most of the examples below
will work without change (i.e. the same binary) on all supported platforms. How cool is that?

Although samples are all present in the [samples](https://github.com/kidoman/embd/tree/master/samples) folder,
we will show a few choice examples here.

Use the **LED** driver to toggle LEDs on the BBB:

	import "github.com/kidoman/embd"
	...
	embd.InitLED()
	defer embd.CloseLED()
	...
	led, err := embd.NewLED("USR3")
	...
	led.Toggle()

Even shorter while prototyping:

	import "github.com/kidoman/embd"
	...
	embd.InitLED()
	defer embd.CloseLED()
	...
	embd.ToggleLED(3)

NB: **3** == **USR3** for all intents and purposes. The driver is smart enough to figure all this out.

Control **GPIO** pins on the RaspberryPi / BeagleBone Black:

	import "github.com/kidoman/embd"
	...
	embd.InitGPIO()
	defer embd.CloseGPIO()
	...
	embd.SetDirection(10, embd.Out)
	embd.DigitalWrite(10, embd.High)

Could also do:

	import "github.com/kidoman/embd"
	...
	embd.InitGPIO()
	defer embd.CloseGPIO()
	...
	pin, err := embd.NewDigitalPin(10)
	...
	pin.SetDirection(embd.Out)
	pin.Write(embd.High)

Or read data from the **Bosch BMP085** barometric sensor:

	import "github.com/kidoman/embd"
	import "github.com/kidoman/embd/sensor/bmp085"
	...
	bus := embd.NewI2CBus(1)
	...
	baro := bmp085.New(bus)
	...
	temp, err := baro.Temperature()
	altitude, err := baro.Altitude()

Even find out the heading from the **LSM303** magnetometer:

	import "github.com/kidoman/embd"
	import "github.com/kidoman/embd/sensor/lsm303"
	...
	bus := embd.NewI2CBus(1)
	...
	mag := lsm303.New(bus)
	...
	heading, err := mag.Heading()

The above two examples depend on **I2C** and therefore will work without change on almost all
platforms.

## Protocols supported

* **Digital GPIO** [Documentation](http://godoc.org/github.com/kidoman/embd#DigitalPin)
* **Analog GPIO** [Documentation](http://godoc.org/github.com/kidoman/embd#AnalogPin)
* **I2C** [Documentation](http://godoc.org/github.com/kidoman/embd#I2CBus)
* **LED** [Documentation](http://godoc.org/github.com/kidoman/embd#LED)

## Sensors supported

* **TMP006** Thermopile sensor [Documentation](http://godoc.org/github.com/kidoman/embd/sensor/tmp006), [Datasheet](http://www.adafruit.com/datasheets/tmp006.pdf), [Userguide](http://www.adafruit.com/datasheets/tmp006ug.pdf)

* **BMP085** Barometric pressure sensor [Documentation](http://godoc.org/github.com/kidoman/embd/sensor/bmp085), [Datasheet](https://www.sparkfun.com/datasheets/Components/General/BST-BMP085-DS000-05.pdf)

* **BMP180** Barometric pressure sensor [Documentation](http://godoc.org/github.com/kidoman/embd/sensor/bmp180), [Datasheet](http://www.adafruit.com/datasheets/BST-BMP180-DS000-09.pdf)

* **LSM303** Accelerometer and magnetometer [Documentation](http://godoc.org/github.com/kidoman/embd/sensor/lsm303), [Datasheet](https://www.sparkfun.com/datasheets/Sensors/Magneto/LSM303%20Datasheet.pdf)

* **L3GD20** Gyroscope [Documentation](http://godoc.org/github.com/kidoman/embd/sensor/l3gd20), [Datasheet](http://www.adafruit.com/datasheets/L3GD20.pdf)

* **US020** Ultrasonic proximity sensor [Documentation](http://godoc.org/github.com/kidoman/embd/sensor/us020), [Product Page](http://www.digibay.in/sensor/object-detection-and-proximity?product_id=239)

* **BH1750FVI** Luminosity sensor [Documentation](http://godoc.org/github.com/kidoman/embd/sensor/us020), [Datasheet](http://www.elechouse.com/elechouse/images/product/Digital%20light%20Sensor/bh1750fvi-e.pdf)

## Interfaces

* **Keypad(4x3)** [Product Page](http://www.adafruit.com/products/419#Learn)

## Controllers

* **PCA9685** 16-channel, 12-bit PWM Controller with I2C protocol [Documentation](http://godoc.org/github.com/kidoman/embd/controller/pca9685), [Datasheet](http://www.adafruit.com/datasheets/PCA9685.pdf), [Product Page](http://www.adafruit.com/products/815)

* **MCP4725** 12-bit DAC [Documentation](http://godoc.org/github.com/kidoman/embd/controller/mcp4725), [Datasheet](http://www.adafruit.com/datasheets/mcp4725.pdf), [Product Page](http://www.adafruit.com/products/935)

* **ServoBlaster** RPi PWM/PCM based PWM controller [Documentation](http://godoc.org/github.com/kidoman/embd/controller/servoblaster), [Product Page](https://github.com/richardghirst/PiBits/tree/master/ServoBlaster)
