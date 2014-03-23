// RaspberryPi support.
// The following features are supported on Linux kernel 3.8+
//
//	GPIO (digital (rw))
//	I2C

package embd

func init() {
	Register(HostRPi, func(rev int) *Descriptor {
		var pins = rpiRev1Pins
		if rev > 1 {
			pins = rpiRev2Pins
		}

		return &Descriptor{
			GPIODriver: func() GPIODriver {
				return newGPIODriver(pins, newDigitalPin, nil)
			},
			I2CDriver: newI2CDriver,
		}
	})
}

var rpiRev1Pins = PinMap{
	&PinDesc{ID: "P1_3", Aliases: []string{"0", "GPIO_0", "SDA", "I2C0_SDA"}, Caps: CapDigital | CapI2C, DigitalLogical: 0},
	&PinDesc{ID: "P1_5", Aliases: []string{"1", "GPIO_1", "SCL", "I2C0_SCL"}, Caps: CapDigital | CapI2C, DigitalLogical: 1},
	&PinDesc{ID: "P1_7", Aliases: []string{"4", "GPIO_4", "GPCLK0"}, Caps: CapDigital, DigitalLogical: 4},
	&PinDesc{ID: "P1_8", Aliases: []string{"14", "GPIO_14", "TXD", "UART0_TXD"}, Caps: CapDigital | CapUART, DigitalLogical: 14},
	&PinDesc{ID: "P1_10", Aliases: []string{"15", "GPIO_15", "RXD", "UART0_RXD"}, Caps: CapDigital | CapUART, DigitalLogical: 15},
	&PinDesc{ID: "P1_11", Aliases: []string{"17", "GPIO_17"}, Caps: CapDigital, DigitalLogical: 17},
	&PinDesc{ID: "P1_12", Aliases: []string{"18", "GPIO_18", "PCM_CLK"}, Caps: CapDigital, DigitalLogical: 18},
	&PinDesc{ID: "P1_13", Aliases: []string{"21", "GPIO_21"}, Caps: CapDigital, DigitalLogical: 21},
	&PinDesc{ID: "P1_15", Aliases: []string{"22", "GPIO_22"}, Caps: CapDigital, DigitalLogical: 22},
	&PinDesc{ID: "P1_16", Aliases: []string{"23", "GPIO_23"}, Caps: CapDigital, DigitalLogical: 23},
	&PinDesc{ID: "P1_18", Aliases: []string{"24", "GPIO_24"}, Caps: CapDigital, DigitalLogical: 24},
	&PinDesc{ID: "P1_19", Aliases: []string{"10", "GPIO_10", "MOSI", "SPI0_MOSI"}, Caps: CapDigital | CapSPI, DigitalLogical: 10},
	&PinDesc{ID: "P1_21", Aliases: []string{"9", "GPIO_9", "MISO", "SPI0_MISO"}, Caps: CapDigital | CapSPI, DigitalLogical: 9},
	&PinDesc{ID: "P1_22", Aliases: []string{"25", "GPIO_25"}, Caps: CapDigital, DigitalLogical: 25},
	&PinDesc{ID: "P1_23", Aliases: []string{"11", "GPIO_11", "SCLK", "SPI0_SCLK"}, Caps: CapDigital | CapSPI, DigitalLogical: 11},
	&PinDesc{ID: "P1_24", Aliases: []string{"8", "GPIO_8", "CE0", "SPI0_CE0_N"}, Caps: CapDigital | CapSPI, DigitalLogical: 8},
	&PinDesc{ID: "P1_26", Aliases: []string{"7", "GPIO_7", "CE1", "SPI0_CE1_N"}, Caps: CapDigital | CapSPI, DigitalLogical: 7},
}

var rpiRev2Pins = PinMap{
	&PinDesc{ID: "P1_3", Aliases: []string{"2", "GPIO_2", "SDA", "I2C1_SDA"}, Caps: CapDigital | CapI2C, DigitalLogical: 2},
	&PinDesc{ID: "P1_5", Aliases: []string{"3", "GPIO_3", "SCL", "I2C1_SCL"}, Caps: CapDigital | CapI2C, DigitalLogical: 3},
	&PinDesc{ID: "P1_7", Aliases: []string{"4", "GPIO_4", "GPCLK0"}, Caps: CapDigital, DigitalLogical: 4},
	&PinDesc{ID: "P1_8", Aliases: []string{"14", "GPIO_14", "TXD", "UART0_TXD"}, Caps: CapDigital | CapUART, DigitalLogical: 14},
	&PinDesc{ID: "P1_10", Aliases: []string{"15", "GPIO_15", "RXD", "UART0_RXD"}, Caps: CapDigital | CapUART, DigitalLogical: 15},
	&PinDesc{ID: "P1_11", Aliases: []string{"17", "GPIO_17"}, Caps: CapDigital, DigitalLogical: 17},
	&PinDesc{ID: "P1_12", Aliases: []string{"18", "GPIO_18", "PCM_CLK"}, Caps: CapDigital, DigitalLogical: 18},
	&PinDesc{ID: "P1_13", Aliases: []string{"27", "GPIO_27"}, Caps: CapDigital, DigitalLogical: 27},
	&PinDesc{ID: "P1_15", Aliases: []string{"22", "GPIO_22"}, Caps: CapDigital, DigitalLogical: 22},
	&PinDesc{ID: "P1_16", Aliases: []string{"23", "GPIO_23"}, Caps: CapDigital, DigitalLogical: 23},
	&PinDesc{ID: "P1_18", Aliases: []string{"24", "GPIO_24"}, Caps: CapDigital, DigitalLogical: 24},
	&PinDesc{ID: "P1_19", Aliases: []string{"10", "GPIO_10", "MOSI", "SPI0_MOSI"}, Caps: CapDigital | CapSPI, DigitalLogical: 10},
	&PinDesc{ID: "P1_21", Aliases: []string{"9", "GPIO_9", "MISO", "SPI0_MISO"}, Caps: CapDigital | CapSPI, DigitalLogical: 9},
	&PinDesc{ID: "P1_22", Aliases: []string{"25", "GPIO_25"}, Caps: CapDigital, DigitalLogical: 25},
	&PinDesc{ID: "P1_23", Aliases: []string{"11", "GPIO_11", "SCLK", "SPI0_SCLK"}, Caps: CapDigital | CapSPI, DigitalLogical: 11},
	&PinDesc{ID: "P1_24", Aliases: []string{"8", "GPIO_8", "CE0", "SPI0_CE0_N"}, Caps: CapDigital | CapSPI, DigitalLogical: 8},
	&PinDesc{ID: "P1_26", Aliases: []string{"7", "GPIO_7", "CE1", "SPI0_CE1_N"}, Caps: CapDigital | CapSPI, DigitalLogical: 7},
}
