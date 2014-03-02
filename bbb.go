package embd

func init() {
	Describers[HostBBB] = func(rev int) *Descriptor {
		return &Descriptor{
			GPIO: func() GPIO {
				return newGPIODriver(bbbPins)
			},
			I2C: newI2CDriver,
		}
	}
}

var bbbPins = PinMap{
	&PinDesc{66, []string{"P8_07", "GPIO_66", "TIMER4"}, CapNormal | CapGPMC},
	&PinDesc{67, []string{"P8_08", "GPIO_67", "TIMER7"}, CapNormal | CapGPMC},
	&PinDesc{69, []string{"P8_09", "GPIO_69", "TIMER5"}, CapNormal | CapGPMC},
	&PinDesc{68, []string{"P8_10", "GPIO_68", "TIMER6"}, CapNormal | CapGPMC},
	&PinDesc{45, []string{"P8_11", "GPIO_45"}, CapNormal | CapGPMC},
	&PinDesc{44, []string{"P8_12", "GPIO_44"}, CapNormal | CapGPMC},
	&PinDesc{23, []string{"P8_13", "GPIO_23", "EHRPWM2B"}, CapNormal | CapGPMC},
	&PinDesc{26, []string{"P8_14", "GPIO_26"}, CapNormal | CapGPMC},
	&PinDesc{47, []string{"P8_15", "GPIO_47"}, CapNormal | CapGPMC},
	&PinDesc{46, []string{"P8_16", "GPIO_46"}, CapNormal | CapGPMC},
	&PinDesc{27, []string{"P8_17", "GPIO_27"}, CapNormal | CapGPMC},
	&PinDesc{65, []string{"P8_18", "GPIO_65"}, CapNormal | CapGPMC},
	&PinDesc{22, []string{"P8_19", "GPIO_22", "EHRPWM2A"}, CapNormal | CapGPMC},
	&PinDesc{61, []string{"P8_26", "GPIO_61"}, CapNormal | CapGPMC},
	&PinDesc{86, []string{"P8_27", "GPIO_86"}, CapNormal | CapLCD},
	&PinDesc{88, []string{"P8_28", "GPIO_88"}, CapNormal | CapLCD},
	&PinDesc{87, []string{"P8_29", "GPIO_87"}, CapNormal | CapLCD},
	&PinDesc{89, []string{"P8_30", "GPIO_89"}, CapNormal | CapLCD},
	&PinDesc{10, []string{"P8_31", "GPIO_10", "UART5_CTSN"}, CapNormal | CapLCD},
	&PinDesc{11, []string{"P8_32", "GPIO_11", "UART5_RTSN"}, CapNormal | CapLCD},
	&PinDesc{9, []string{"P8_33", "GPIO_9 ", "UART4_RTSN"}, CapNormal | CapLCD},
	&PinDesc{81, []string{"P8_34", "GPIO_81", "UART3_RTSN"}, CapNormal | CapLCD},
	&PinDesc{8, []string{"P8_35", "GPIO_8 ", "UART4_CTSN"}, CapNormal | CapLCD},
	&PinDesc{80, []string{"P8_36", "GPIO_80", "UART3_CTSN"}, CapNormal | CapLCD},
	&PinDesc{78, []string{"P8_37", "GPIO_78", "UART5_TXD"}, CapNormal | CapLCD},
	&PinDesc{79, []string{"P8_38", "GPIO_79", "UART5_RXD"}, CapNormal | CapLCD},
	&PinDesc{76, []string{"P8_39", "GPIO_76"}, CapNormal | CapLCD},
	&PinDesc{77, []string{"P8_40", "GPIO_77"}, CapNormal | CapLCD},
	&PinDesc{74, []string{"P8_41", "GPIO_74"}, CapNormal | CapLCD},
	&PinDesc{75, []string{"P8_42", "GPIO_75"}, CapNormal | CapLCD},
	&PinDesc{72, []string{"P8_43", "GPIO_72"}, CapNormal | CapLCD},
	&PinDesc{73, []string{"P8_44", "GPIO_73"}, CapNormal | CapLCD},
	&PinDesc{70, []string{"P8_45", "GPIO_70"}, CapNormal | CapLCD},
	&PinDesc{71, []string{"P8_46", "GPIO_71"}, CapNormal | CapLCD},

	&PinDesc{30, []string{"P9_11", "GPIO_30", "UART4_RXD"}, CapNormal | CapUART},
	&PinDesc{60, []string{"P9_12", "GPIO_60", "GPIO1_28"}, CapNormal},
	&PinDesc{31, []string{"P9_13", "GPIO_31", "UART4_TXD"}, CapNormal | CapUART},
	&PinDesc{50, []string{"P9_14", "GPIO_50", "EHRPWM1A"}, CapNormal | CapPWM},
	&PinDesc{48, []string{"P9_15", "GPIO_48", "GPIO1_16"}, CapNormal},
	&PinDesc{51, []string{"P9_16", "GPIO_51", "EHRPWM1B"}, CapNormal | CapPWM},
	&PinDesc{5, []string{"P9_17", "GPIO_5", "I2C1_SCL"}, CapNormal | CapI2C},
	&PinDesc{4, []string{"P9_18", "GPIO_4", "I2C1_SDA"}, CapNormal | CapI2C},
	&PinDesc{13, []string{"P9_19", "GPIO_13", "I2C2_SCL"}, CapNormal | CapI2C},
	&PinDesc{12, []string{"P9_20", "GPIO_12", "I2C2_SDA"}, CapNormal | CapI2C},
	&PinDesc{3, []string{"P9_21", "GPIO_3", "UART2_TXD"}, CapNormal | CapUART},
	&PinDesc{2, []string{"P9_22", "GPIO_2", "UART2_RXD"}, CapNormal | CapUART},
	&PinDesc{49, []string{"P9_23", "GPIO_49", "GPIO1_17"}, CapNormal},
	&PinDesc{15, []string{"P9_24", "GPIO_15", "UART1_TXD"}, CapNormal | CapUART},
	&PinDesc{117, []string{"P9_25", "GPIO_117", "GPIO3_21"}, CapNormal},
	&PinDesc{14, []string{"P9_26", "GPIO_14", "UART1_RXD"}, CapNormal | CapUART},
	&PinDesc{115, []string{"P9_27", "GPIO_115", "GPIO3_19"}, CapNormal},
	&PinDesc{113, []string{"P9_28", "GPIO_113", "SPI1_CS0"}, CapNormal | CapSPI},
	&PinDesc{111, []string{"P9_29", "GPIO_111", "SPI1_D0"}, CapNormal | CapSPI},
	&PinDesc{112, []string{"P9_30", "GPIO_112", "SPI1_D1"}, CapNormal | CapSPI},
	&PinDesc{110, []string{"P9_31", "GPIO_110", "SPI1_SCLK"}, CapNormal | CapSPI},
}
