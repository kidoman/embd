package rpi

import (
	"github.com/kidoman/embd/host/generic/linux/gpio"
)

var rev1Pins = gpio.PinMap{
	&gpio.PinDesc{0, []string{"P1_3", "GPIO_0", "SDA", "I2C0_SDA"}, gpio.Normal | gpio.I2C},
	&gpio.PinDesc{1, []string{"P1_5", "GPIO_1", "SCL", "I2C0_SCL"}, gpio.Normal | gpio.I2C},
	&gpio.PinDesc{4, []string{"P1_7", "GPIO_4", "GPCLK0"}, gpio.Normal},
	&gpio.PinDesc{14, []string{"P1_8", "GPIO_14", "TXD", "UART0_TXD"}, gpio.Normal | gpio.UART},
	&gpio.PinDesc{15, []string{"P1_10", "GPIO_15", "RXD", "UART0_RXD"}, gpio.Normal | gpio.UART},
	&gpio.PinDesc{17, []string{"P1_11", "GPIO_17"}, gpio.Normal},
	&gpio.PinDesc{18, []string{"P1_12", "GPIO_18", "PCM_CLK"}, gpio.Normal},
	&gpio.PinDesc{21, []string{"P1_13", "GPIO_21"}, gpio.Normal},
	&gpio.PinDesc{22, []string{"P1_15", "GPIO_22"}, gpio.Normal},
	&gpio.PinDesc{23, []string{"P1_16", "GPIO_23"}, gpio.Normal},
	&gpio.PinDesc{24, []string{"P1_18", "GPIO_24"}, gpio.Normal},
	&gpio.PinDesc{10, []string{"P1_19", "GPIO_10", "MOSI", "SPI0_MOSI"}, gpio.Normal | gpio.SPI},
	&gpio.PinDesc{9, []string{"P1_21", "GPIO_9", "MISO", "SPI0_MISO"}, gpio.Normal | gpio.SPI},
	&gpio.PinDesc{25, []string{"P1_22", "GPIO_25"}, gpio.Normal},
	&gpio.PinDesc{11, []string{"P1_23", "GPIO_11", "SCLK", "SPI0_SCLK"}, gpio.Normal | gpio.SPI},
	&gpio.PinDesc{8, []string{"P1_24", "GPIO_8", "CE0", "SPI0_CE0_N"}, gpio.Normal | gpio.SPI},
	&gpio.PinDesc{7, []string{"P1_26", "GPIO_7", "CE1", "SPI0_CE1_N"}, gpio.Normal | gpio.SPI},
}

var rev2Pins = gpio.PinMap{
	&gpio.PinDesc{2, []string{"P1_3", "GPIO_2", "SDA", "I2C1_SDA"}, gpio.Normal | gpio.I2C},
	&gpio.PinDesc{3, []string{"P1_5", "GPIO_3", "SCL", "I2C1_SCL"}, gpio.Normal | gpio.I2C},
	&gpio.PinDesc{4, []string{"P1_7", "GPIO_4", "GPCLK0"}, gpio.Normal},
	&gpio.PinDesc{14, []string{"P1_8", "GPIO_14", "TXD", "UART0_TXD"}, gpio.Normal | gpio.UART},
	&gpio.PinDesc{15, []string{"P1_10", "GPIO_15", "RXD", "UART0_RXD"}, gpio.Normal | gpio.UART},
	&gpio.PinDesc{17, []string{"P1_11", "GPIO_17"}, gpio.Normal},
	&gpio.PinDesc{18, []string{"P1_12", "GPIO_18", "PCM_CLK"}, gpio.Normal},
	&gpio.PinDesc{27, []string{"P1_13", "GPIO_27"}, gpio.Normal},
	&gpio.PinDesc{22, []string{"P1_15", "GPIO_22"}, gpio.Normal},
	&gpio.PinDesc{23, []string{"P1_16", "GPIO_23"}, gpio.Normal},
	&gpio.PinDesc{24, []string{"P1_18", "GPIO_24"}, gpio.Normal},
	&gpio.PinDesc{10, []string{"P1_19", "GPIO_10", "MOSI", "SPI0_MOSI"}, gpio.Normal | gpio.SPI},
	&gpio.PinDesc{9, []string{"P1_21", "GPIO_9", "MISO", "SPI0_MISO"}, gpio.Normal | gpio.SPI},
	&gpio.PinDesc{25, []string{"P1_22", "GPIO_25"}, gpio.Normal},
	&gpio.PinDesc{11, []string{"P1_23", "GPIO_11", "SCLK", "SPI0_SCLK"}, gpio.Normal | gpio.SPI},
	&gpio.PinDesc{8, []string{"P1_24", "GPIO_8", "CE0", "SPI0_CE0_N"}, gpio.Normal | gpio.SPI},
	&gpio.PinDesc{7, []string{"P1_26", "GPIO_7", "CE1", "SPI0_CE1_N"}, gpio.Normal | gpio.SPI},
}
