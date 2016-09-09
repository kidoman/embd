/* doc
Package embd provides a hardware abstraction layer for doing embedded programming
on supported platforms like the Raspberry Pi and BeagleBone Black. Most of the examples below
will work without change (i.e. the same binary) on all supported platforms. How cool is that?

Although samples are all present in the samples folder, we will show a few choice examples here.

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
