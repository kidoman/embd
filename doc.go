/*
Package embd provides a hardware abstraction layer for doing embedded programming
on supported platforms like the Raspberry Pi, BeagleBone Black and CHIP. Most of the examples below
will work without change (i.e. the same binary) on all supported platforms.

== Overall structure

It's best to think of the top-level embd package as a switchboard that doesn't implement anything
on its own but rather relies on sub-packages for hosts drivers and devices and stitches them
together. The exports in the top-level package serve a number of different purposes,
which can be confusing at first:
- it defines a number of driver interfaces, such as the GPIODriver, this is the interface that
the driver for each specific platform must implement and is not something of concern to the
typical user.
- it defines the main low-level hardware interface types: analog pins, digital pins,
interrupt pins, I2Cbuses, SPI buses, PWM pins and LEDs. Each type has a New function to
instantiate one of these pins or buses.
- it defines a number of InitXXX functions that initialize the various drivers, however, these
are called by the coresponding NewXXX functions, so can be ignored.
- it defines a number of top-level convenience functions, such as DigitalWrite, that can be
called as 1-liners instead of first instantiating a DigitalPin and then writing to it

To get started a host driver needs to be registered with the top-level embd package. This is
most easily accomplished by doing an "underscore import" on of the sub-packages of embd/host,
e.g., `import _ "github.com/kidoman/embd/host/chip"`. An `Init()` function in the host driver
registers all the individual drivers with embd.

After getting the host driver the next step might be to instantiate a GPIO pin using
`NewDigitalPin` or an I2CBus using `NewI2CBus`. Such a pin or bus can be used directly but
often it is passed into the initializer of a sensor, controller or other user-level driver
which provides a high-level interface to some device. For example, the New function
for the BMP180 type in the `embd/sensor/bmp180` package takes an I2CBus as argument, which
it will use to reach the sensor.

== Samples

This section shows a few choice samples, more are available in the samples folder.

Use the LED driver to toggle LEDs on the BBB:

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

BBB + PWM:

	import "github.com/kidoman/embd"
	...
	embd.InitGPIO()
	defer embd.CloseGPIO()
	...
	pwm, _ := embd.NewPWMPin("P9_14")
	defer pwm.Close()
	...
	pwm.SetDuty(1000)

Control GPIO pins on the RaspberryPi / BeagleBone Black:

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

Or read data from the Bosch BMP085 barometric sensor:

	import "github.com/kidoman/embd"
	import "github.com/kidoman/embd/sensor/bmp085"
	...
	bus := embd.NewI2CBus(1)
	...
	baro := bmp085.New(bus)
	...
	temp, err := baro.Temperature()
	altitude, err := baro.Altitude()

Even find out the heading from the LSM303 magnetometer:

	import "github.com/kidoman/embd"
	import "github.com/kidoman/embd/sensor/lsm303"
	...
	bus := embd.NewI2CBus(1)
	...
	mag := lsm303.New(bus)
	...
	heading, err := mag.Heading()

The above two examples depend on I2C and therefore will work without change on almost all
platforms.
*/
package embd
