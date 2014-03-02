package embd

func init() {
	Describers[HostRPi] = func(rev int) *Descriptor {
		var pins = rpiRev1Pins
		if rev > 1 {
			pins = rpiRev2Pins
		}

		return &Descriptor{
			GPIO: func() GPIO {
				return newGPIODriver(pins)
			},
			I2C: newI2CDriver,
		}
	}
}

var rpiRev1Pins = PinMap{
	&PinDesc{0, []string{"P1_3", "GPIO_0", "SDA", "I2C0_SDA"}, CapNormal | CapI2C},
	&PinDesc{1, []string{"P1_5", "GPIO_1", "SCL", "I2C0_SCL"}, CapNormal | CapI2C},
	&PinDesc{4, []string{"P1_7", "GPIO_4", "GPCLK0"}, CapNormal},
	&PinDesc{14, []string{"P1_8", "GPIO_14", "TXD", "UART0_TXD"}, CapNormal | CapUART},
	&PinDesc{15, []string{"P1_10", "GPIO_15", "RXD", "UART0_RXD"}, CapNormal | CapUART},
	&PinDesc{17, []string{"P1_11", "GPIO_17"}, CapNormal},
	&PinDesc{18, []string{"P1_12", "GPIO_18", "PCM_CLK"}, CapNormal},
	&PinDesc{21, []string{"P1_13", "GPIO_21"}, CapNormal},
	&PinDesc{22, []string{"P1_15", "GPIO_22"}, CapNormal},
	&PinDesc{23, []string{"P1_16", "GPIO_23"}, CapNormal},
	&PinDesc{24, []string{"P1_18", "GPIO_24"}, CapNormal},
	&PinDesc{10, []string{"P1_19", "GPIO_10", "MOSI", "SPI0_MOSI"}, CapNormal | CapSPI},
	&PinDesc{9, []string{"P1_21", "GPIO_9", "MISO", "SPI0_MISO"}, CapNormal | CapSPI},
	&PinDesc{25, []string{"P1_22", "GPIO_25"}, CapNormal},
	&PinDesc{11, []string{"P1_23", "GPIO_11", "SCLK", "SPI0_SCLK"}, CapNormal | CapSPI},
	&PinDesc{8, []string{"P1_24", "GPIO_8", "CE0", "SPI0_CE0_N"}, CapNormal | CapSPI},
	&PinDesc{7, []string{"P1_26", "GPIO_7", "CE1", "SPI0_CE1_N"}, CapNormal | CapSPI},
}

var rpiRev2Pins = PinMap{
	&PinDesc{2, []string{"P1_3", "GPIO_2", "SDA", "I2C1_SDA"}, CapNormal | CapI2C},
	&PinDesc{3, []string{"P1_5", "GPIO_3", "SCL", "I2C1_SCL"}, CapNormal | CapI2C},
	&PinDesc{4, []string{"P1_7", "GPIO_4", "GPCLK0"}, CapNormal},
	&PinDesc{14, []string{"P1_8", "GPIO_14", "TXD", "UART0_TXD"}, CapNormal | CapUART},
	&PinDesc{15, []string{"P1_10", "GPIO_15", "RXD", "UART0_RXD"}, CapNormal | CapUART},
	&PinDesc{17, []string{"P1_11", "GPIO_17"}, CapNormal},
	&PinDesc{18, []string{"P1_12", "GPIO_18", "PCM_CLK"}, CapNormal},
	&PinDesc{27, []string{"P1_13", "GPIO_27"}, CapNormal},
	&PinDesc{22, []string{"P1_15", "GPIO_22"}, CapNormal},
	&PinDesc{23, []string{"P1_16", "GPIO_23"}, CapNormal},
	&PinDesc{24, []string{"P1_18", "GPIO_24"}, CapNormal},
	&PinDesc{10, []string{"P1_19", "GPIO_10", "MOSI", "SPI0_MOSI"}, CapNormal | CapSPI},
	&PinDesc{9, []string{"P1_21", "GPIO_9", "MISO", "SPI0_MISO"}, CapNormal | CapSPI},
	&PinDesc{25, []string{"P1_22", "GPIO_25"}, CapNormal},
	&PinDesc{11, []string{"P1_23", "GPIO_11", "SCLK", "SPI0_SCLK"}, CapNormal | CapSPI},
	&PinDesc{8, []string{"P1_24", "GPIO_8", "CE0", "SPI0_CE0_N"}, CapNormal | CapSPI},
	&PinDesc{7, []string{"P1_26", "GPIO_7", "CE1", "SPI0_CE1_N"}, CapNormal | CapSPI},
}
