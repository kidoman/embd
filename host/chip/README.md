# Using embd on CHIP

The CHIP drivers support gpio, I2C, SPI, and pin interrupts. Not supported are PWM or LED.
The names of the pins on chip have multiple aliases. The official CHIP pin names are supported, 
for example XIO-P1 or LCD-D2 and the pin number are also supported, such as U14-14 (same as XIO-P1)
or U13-17. Some of the alternate function names are also supported, like "SPI2_MOSI", and the
linux 4.4 kernel gpio pin numbers as well, e.g., 1017 for XIO-P1. Finally, the official GPIO pins
(XIO-P0 thru XIO-P7) can be addressed as gpio0-gpio7.

A simple demo to blink an LED connected with a small resistor between XIO-P6 and 3.3V is

```
package main
import (
	"time"
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/chip"
)

func main() {
	embd.InitGPIO()
	defer embd.CloseGPIO()

	embd.SetDirection("gpio6", embd.Out)
	on := 0
	for {
		embd.DigitalWrite("gpio6", on)
		on = 1 - on
		time.Sleep(250 * time.Millisecond)
	}
}
```
Run it as root: `sudo ./blinky`

