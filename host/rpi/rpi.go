/*
	Package rpi provides Raspberry Pi (including A+/B+) support.
	The following features are supported on Linux kernel 3.8+

	GPIO (digital (rw))
	I²C
	LED
*/
package rpi

import (
	"github.com/kidoman/embd"
	"github.com/kidoman/embd/host/generic"
)

var spiDeviceMinor = 0

var rev1Pins = embd.PinMap{
	&embd.PinDesc{ID: "P1_3", Aliases: []string{"0", "GPIO_0", "SDA", "I2C0_SDA"}, Caps: embd.CapDigital | embd.CapI2C, DigitalLogical: 0},
	&embd.PinDesc{ID: "P1_5", Aliases: []string{"1", "GPIO_1", "SCL", "I2C0_SCL"}, Caps: embd.CapDigital | embd.CapI2C, DigitalLogical: 1},
	&embd.PinDesc{ID: "P1_7", Aliases: []string{"4", "GPIO_4", "GPCLK0"}, Caps: embd.CapDigital, DigitalLogical: 4},
	&embd.PinDesc{ID: "P1_8", Aliases: []string{"14", "GPIO_14", "TXD", "UART0_TXD"}, Caps: embd.CapDigital | embd.CapUART, DigitalLogical: 14},
	&embd.PinDesc{ID: "P1_10", Aliases: []string{"15", "GPIO_15", "RXD", "UART0_RXD"}, Caps: embd.CapDigital | embd.CapUART, DigitalLogical: 15},
	&embd.PinDesc{ID: "P1_11", Aliases: []string{"17", "GPIO_17"}, Caps: embd.CapDigital, DigitalLogical: 17},
	&embd.PinDesc{ID: "P1_12", Aliases: []string{"18", "GPIO_18", "PCM_CLK"}, Caps: embd.CapDigital, DigitalLogical: 18},
	&embd.PinDesc{ID: "P1_13", Aliases: []string{"21", "GPIO_21"}, Caps: embd.CapDigital, DigitalLogical: 21},
	&embd.PinDesc{ID: "P1_15", Aliases: []string{"22", "GPIO_22"}, Caps: embd.CapDigital, DigitalLogical: 22},
	&embd.PinDesc{ID: "P1_16", Aliases: []string{"23", "GPIO_23"}, Caps: embd.CapDigital, DigitalLogical: 23},
	&embd.PinDesc{ID: "P1_18", Aliases: []string{"24", "GPIO_24"}, Caps: embd.CapDigital, DigitalLogical: 24},
	&embd.PinDesc{ID: "P1_19", Aliases: []string{"10", "GPIO_10", "MOSI", "SPI0_MOSI"}, Caps: embd.CapDigital | embd.CapSPI, DigitalLogical: 10},
	&embd.PinDesc{ID: "P1_21", Aliases: []string{"9", "GPIO_9", "MISO", "SPI0_MISO"}, Caps: embd.CapDigital | embd.CapSPI, DigitalLogical: 9},
	&embd.PinDesc{ID: "P1_22", Aliases: []string{"25", "GPIO_25"}, Caps: embd.CapDigital, DigitalLogical: 25},
	&embd.PinDesc{ID: "P1_23", Aliases: []string{"11", "GPIO_11", "SCLK", "SPI0_SCLK"}, Caps: embd.CapDigital | embd.CapSPI, DigitalLogical: 11},
	&embd.PinDesc{ID: "P1_24", Aliases: []string{"8", "GPIO_8", "CE0", "SPI0_CE0_N"}, Caps: embd.CapDigital | embd.CapSPI, DigitalLogical: 8},
	&embd.PinDesc{ID: "P1_26", Aliases: []string{"7", "GPIO_7", "CE1", "SPI0_CE1_N"}, Caps: embd.CapDigital | embd.CapSPI, DigitalLogical: 7},
}

var rev2Pins = embd.PinMap{
	&embd.PinDesc{ID: "P1_3", Aliases: []string{"2", "GPIO_2", "SDA", "I2C1_SDA"}, Caps: embd.CapDigital | embd.CapI2C, DigitalLogical: 2},
	&embd.PinDesc{ID: "P1_5", Aliases: []string{"3", "GPIO_3", "SCL", "I2C1_SCL"}, Caps: embd.CapDigital | embd.CapI2C, DigitalLogical: 3},
	&embd.PinDesc{ID: "P1_7", Aliases: []string{"4", "GPIO_4", "GPCLK0"}, Caps: embd.CapDigital, DigitalLogical: 4},
	&embd.PinDesc{ID: "P1_8", Aliases: []string{"14", "GPIO_14", "TXD", "UART0_TXD"}, Caps: embd.CapDigital | embd.CapUART, DigitalLogical: 14},
	&embd.PinDesc{ID: "P1_10", Aliases: []string{"15", "GPIO_15", "RXD", "UART0_RXD"}, Caps: embd.CapDigital | embd.CapUART, DigitalLogical: 15},
	&embd.PinDesc{ID: "P1_11", Aliases: []string{"17", "GPIO_17"}, Caps: embd.CapDigital, DigitalLogical: 17},
	&embd.PinDesc{ID: "P1_12", Aliases: []string{"18", "GPIO_18", "PCM_CLK"}, Caps: embd.CapDigital, DigitalLogical: 18},
	&embd.PinDesc{ID: "P1_13", Aliases: []string{"27", "GPIO_27"}, Caps: embd.CapDigital, DigitalLogical: 27},
	&embd.PinDesc{ID: "P1_15", Aliases: []string{"22", "GPIO_22"}, Caps: embd.CapDigital, DigitalLogical: 22},
	&embd.PinDesc{ID: "P1_16", Aliases: []string{"23", "GPIO_23"}, Caps: embd.CapDigital, DigitalLogical: 23},
	&embd.PinDesc{ID: "P1_18", Aliases: []string{"24", "GPIO_24"}, Caps: embd.CapDigital, DigitalLogical: 24},
	&embd.PinDesc{ID: "P1_19", Aliases: []string{"10", "GPIO_10", "MOSI", "SPI0_MOSI"}, Caps: embd.CapDigital | embd.CapSPI, DigitalLogical: 10},
	&embd.PinDesc{ID: "P1_21", Aliases: []string{"9", "GPIO_9", "MISO", "SPI0_MISO"}, Caps: embd.CapDigital | embd.CapSPI, DigitalLogical: 9},
	&embd.PinDesc{ID: "P1_22", Aliases: []string{"25", "GPIO_25"}, Caps: embd.CapDigital, DigitalLogical: 25},
	&embd.PinDesc{ID: "P1_23", Aliases: []string{"11", "GPIO_11", "SCLK", "SPI0_SCLK"}, Caps: embd.CapDigital | embd.CapSPI, DigitalLogical: 11},
	&embd.PinDesc{ID: "P1_24", Aliases: []string{"8", "GPIO_8", "CE0", "SPI0_CE0_N"}, Caps: embd.CapDigital | embd.CapSPI, DigitalLogical: 8},
	&embd.PinDesc{ID: "P1_26", Aliases: []string{"7", "GPIO_7", "CE1", "SPI0_CE1_N"}, Caps: embd.CapDigital | embd.CapSPI, DigitalLogical: 7},
}

// This is the same as the Rev 2 for the first 26 pins.
var rev3Pins = append(append(embd.PinMap(nil), rev2Pins...), embd.PinMap{
	&embd.PinDesc{ID: "P1_29", Aliases: []string{"5", "GPIO_5"}, Caps: embd.CapDigital, DigitalLogical: 5},
	&embd.PinDesc{ID: "P1_31", Aliases: []string{"6", "GPIO_6"}, Caps: embd.CapDigital, DigitalLogical: 6},
	&embd.PinDesc{ID: "P1_32", Aliases: []string{"12", "GPIO_12"}, Caps: embd.CapDigital, DigitalLogical: 12},
	&embd.PinDesc{ID: "P1_33", Aliases: []string{"13", "GPIO_13"}, Caps: embd.CapDigital, DigitalLogical: 13},
	&embd.PinDesc{ID: "P1_35", Aliases: []string{"19", "GPIO_19"}, Caps: embd.CapDigital, DigitalLogical: 19},
	&embd.PinDesc{ID: "P1_36", Aliases: []string{"16", "GPIO_16"}, Caps: embd.CapDigital, DigitalLogical: 16},
	&embd.PinDesc{ID: "P1_37", Aliases: []string{"26", "GPIO_26"}, Caps: embd.CapDigital, DigitalLogical: 26},
	&embd.PinDesc{ID: "P1_38", Aliases: []string{"20", "GPIO_20"}, Caps: embd.CapDigital, DigitalLogical: 20},
	&embd.PinDesc{ID: "P1_40", Aliases: []string{"21", "GPIO_21"}, Caps: embd.CapDigital, DigitalLogical: 21},
}...)

var ledMap = embd.LEDMap{
	"led0": []string{"0", "led0", "LED0"},
}

func init() {
	embd.Register(embd.HostRPi, func(rev int) *embd.Descriptor {
		// Refer to http://elinux.org/RPi_HardwareHistory#Board_Revision_History
		// for details.
		pins := rev3Pins
		if rev < 16 {
			pins = rev2Pins
		}
		if rev < 4 {
			pins = rev1Pins
		}

		return &embd.Descriptor{
			GPIODriver: func() embd.GPIODriver {
				return embd.NewGPIODriver(pins, generic.NewDigitalPin, nil, nil)
			},
			I2CDriver: func() embd.I2CDriver {
				return embd.NewI2CDriver(generic.NewI2CBus)
			},
			LEDDriver: func() embd.LEDDriver {
				return embd.NewLEDDriver(ledMap, generic.NewLED)
			},
			SPIDriver: func() embd.SPIDriver {
				return embd.NewSPIDriver(spiDeviceMinor, generic.NewSPIBus, nil)
			},
		}
	})
}
