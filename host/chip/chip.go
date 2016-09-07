// Copyright 2016 by Thorsten von Eicken, see LICENSE file

// Package chip provides NextThing C.H.I.P. support.
// References:
//   http://docs.getchip.com/chip.html#chip-hardware
//   http://www.chip-community.org/index.php/Hardware_Information
//
// The following features are supported on Linux kernel 4.4+
//   GPIO (digital (rw))
//   I²C
//   SPI
// Could add LED support by following https://bbs.nextthing.co/t/pwr-and-stat-leds/748/5

package chip

import (
	"github.com/kidoman/embd"
	"github.com/kidoman/embd/host/generic"
)

var spiDeviceMinor = 32766

var chipPins = embd.PinMap{
	// official GPIO pins (U14 connector) using the pcf8574a
	&embd.PinDesc{"XIO-P0", []string{"1016", "0", "U14-13", "gpio0"}, embd.CapDigital, 1016, 0},
	&embd.PinDesc{"XIO-P1", []string{"1017", "1", "U14-14", "gpio1"}, embd.CapDigital, 1017, 0},
	&embd.PinDesc{"XIO-P2", []string{"1018", "2", "U14-15", "gpio2"}, embd.CapDigital, 1018, 0},
	&embd.PinDesc{"XIO-P3", []string{"1019", "3", "U14-16", "gpio3"}, embd.CapDigital, 1019, 0},
	&embd.PinDesc{"XIO-P4", []string{"1020", "4", "U14-17", "gpio4"}, embd.CapDigital, 1020, 0},
	&embd.PinDesc{"XIO-P5", []string{"1021", "5", "U14-18", "gpio5"}, embd.CapDigital, 1021, 0},
	&embd.PinDesc{"XIO-P6", []string{"1022", "6", "U14-19", "gpio6"}, embd.CapDigital, 1022, 0},
	&embd.PinDesc{"XIO-P7", []string{"1023", "7", "U14-20", "gpio7"}, embd.CapDigital, 1023, 0},

	// pins usable on the U13 connector
	&embd.PinDesc{"TWI1-SDA", []string{"48", "U13-9", "I2C0_SDA"}, embd.CapDigital | embd.CapI2C, 48, 0},
	&embd.PinDesc{"TWI1-SCK", []string{"47", "U13-11", "I2C0_SCK"}, embd.CapDigital | embd.CapI2C, 47, 0},
	&embd.PinDesc{"PWM0", []string{"34", "U13-18"}, embd.CapDigital | embd.CapPWM, 34, 0},
	&embd.PinDesc{"LCD-D2", []string{"98", "U13-17"}, embd.CapDigital, 98, 0},
	&embd.PinDesc{"LCD-D3", []string{"99", "U13-20"}, embd.CapDigital, 99, 0},
	&embd.PinDesc{"LCD-D4", []string{"100", "U13-19"}, embd.CapDigital, 100, 0},
	&embd.PinDesc{"LCD-D5", []string{"101", "U13-22"}, embd.CapDigital, 101, 0},
	&embd.PinDesc{"LCD-D6", []string{"102", "U13-21"}, embd.CapDigital, 102, 0},
	&embd.PinDesc{"LCD-D7", []string{"103", "U13-24"}, embd.CapDigital, 103, 0},
	&embd.PinDesc{"LCD-D10", []string{"106", "U13-23"}, embd.CapDigital, 106, 0},
	&embd.PinDesc{"LCD-D11", []string{"107", "U13-26"}, embd.CapDigital, 107, 0},
	&embd.PinDesc{"LCD-D12", []string{"108", "U13-25"}, embd.CapDigital, 108, 0},
	&embd.PinDesc{"LCD-D13", []string{"109", "U13-28"}, embd.CapDigital, 109, 0},
	&embd.PinDesc{"LCD-D14", []string{"110", "U13-27"}, embd.CapDigital, 110, 0},
	&embd.PinDesc{"LCD-D15", []string{"111", "U13-30"}, embd.CapDigital, 111, 0},
	&embd.PinDesc{"LCD-D18", []string{"114", "U13-29"}, embd.CapDigital, 114, 0},
	&embd.PinDesc{"LCD-D19", []string{"115", "U13-32"}, embd.CapDigital, 115, 0},
	&embd.PinDesc{"LCD-D20", []string{"116", "U13-31"}, embd.CapDigital, 116, 0},
	&embd.PinDesc{"LCD-D21", []string{"117", "U13-34"}, embd.CapDigital, 117, 0},
	&embd.PinDesc{"LCD-D22", []string{"118", "U13-33"}, embd.CapDigital, 118, 0},
	&embd.PinDesc{"LCD-D23", []string{"119", "U13-36"}, embd.CapDigital, 119, 0},
	&embd.PinDesc{"LCD-CLK", []string{"120", "U13-35"}, embd.CapDigital, 120, 0},
	&embd.PinDesc{"LCD-VSYNC", []string{"123", "U13-37"}, embd.CapDigital, 123, 0},
	&embd.PinDesc{"LCD-HSYNC", []string{"122", "U13-38"}, embd.CapDigital, 122, 0},
	&embd.PinDesc{"LCD-DE", []string{"121", "U13-40"}, embd.CapDigital, 121, 0},

	// pins usable on the U14 connector
	&embd.PinDesc{"UART1-TX", []string{"195", "U14-3", "EINT3"}, embd.CapDigital | embd.CapUART, 195, 0},
	&embd.PinDesc{"UART1-RX", []string{"196", "U14-5", "EINT4"}, embd.CapDigital | embd.CapUART, 196, 0},
	&embd.PinDesc{"AP-EINT1", []string{"193", "U14-23", "EINT1"}, embd.CapDigital, 193, 0},
	&embd.PinDesc{"AP-EINT3", []string{"35", "U14-24", "EINT3"}, embd.CapDigital, 35, 0},
	&embd.PinDesc{"TWI2-SDA", []string{"50", "U14-25", "I2C2_SDA"}, embd.CapDigital | embd.CapI2C, 50, 0},
	&embd.PinDesc{"TWI2-SCK", []string{"49", "U14-26", "I2C2_SCK"}, embd.CapDigital | embd.CapI2C, 49, 0},
	&embd.PinDesc{"CSIPCK", []string{"128", "U14-27", "SPI2_SCO", "SPI2_CS0"}, embd.CapDigital | embd.CapSPI, 128, 0},
	&embd.PinDesc{"CSICK", []string{"129", "U14-28", "SPI2_CLK"}, embd.CapDigital | embd.CapSPI, 129, 0},
	&embd.PinDesc{"CSIHSYNC", []string{"130", "U14-29", "SPI2_MOSI"}, embd.CapDigital | embd.CapSPI, 130, 0},
	&embd.PinDesc{"CSIVSYNC", []string{"131", "U14-30", "SPI2_MISO"}, embd.CapDigital | embd.CapSPI, 131, 0},
	&embd.PinDesc{"CSID0", []string{"132", "U14-31"}, embd.CapDigital, 132, 0},
	&embd.PinDesc{"CSID1", []string{"133", "U14-32"}, embd.CapDigital, 133, 0},
	&embd.PinDesc{"CSID2", []string{"134", "U14-33"}, embd.CapDigital, 134, 0},
	&embd.PinDesc{"CSID3", []string{"135", "U14-34"}, embd.CapDigital, 135, 0},
	&embd.PinDesc{"CSID4", []string{"136", "U14-35"}, embd.CapDigital, 136, 0},
	&embd.PinDesc{"CSID5", []string{"137", "U14-36"}, embd.CapDigital, 137, 0},
	&embd.PinDesc{"CSID6", []string{"138", "U14-37", "UART1_TX"}, embd.CapDigital | embd.CapUART, 138, 0},
	&embd.PinDesc{"CSID7", []string{"139", "U14-38", "UART1_RX"}, embd.CapDigital | embd.CapUART, 139, 0},
}

func init() {
	embd.Register(embd.HostCHIP, func(rev int) *embd.Descriptor {
		return &embd.Descriptor{
			GPIODriver: func() embd.GPIODriver {
				return embd.NewGPIODriver(chipPins, generic.NewDigitalPin, nil, nil)
			},
			I2CDriver: func() embd.I2CDriver {
				return embd.NewI2CDriver(generic.NewI2CBus)
			},
			//LEDDriver: func() embd.LEDDriver {
			//	return embd.NewLEDDriver(ledMap, generic.NewLED)
			//},
			SPIDriver: func() embd.SPIDriver {
				return embd.NewSPIDriver(spiDeviceMinor, generic.NewSPIBus, nil)
			},
		}
	})
}
