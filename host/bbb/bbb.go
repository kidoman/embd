/*
	Package bbb provides BeagleBone Black support.
	The following features are supported on Linux kernel 3.8+

	GPIO (digital (rw), analog (ro), pwm)
	I²C
	LED
*/
package bbb

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/golang/glog"
	"github.com/kidoman/embd"
	"github.com/kidoman/embd/host/generic"
)

var pins = embd.PinMap{
	&embd.PinDesc{ID: "P8_07", Aliases: []string{"66", "GPIO_66", "Caps: TIMER4"}, Caps: embd.CapDigital | embd.CapGPMC, DigitalLogical: 66},
	&embd.PinDesc{ID: "P8_08", Aliases: []string{"67", "GPIO_67", "TIMER7"}, Caps: embd.CapDigital | embd.CapGPMC, DigitalLogical: 67},
	&embd.PinDesc{ID: "P8_09", Aliases: []string{"69", "GPIO_69", "TIMER5"}, Caps: embd.CapDigital | embd.CapGPMC, DigitalLogical: 69},
	&embd.PinDesc{ID: "P8_10", Aliases: []string{"68", "GPIO_68", "TIMER6"}, Caps: embd.CapDigital | embd.CapGPMC, DigitalLogical: 68},
	&embd.PinDesc{ID: "P8_11", Aliases: []string{"45", "GPIO_45"}, Caps: embd.CapDigital | embd.CapGPMC, DigitalLogical: 45},
	&embd.PinDesc{ID: "P8_12", Aliases: []string{"44", "GPIO_44"}, Caps: embd.CapDigital | embd.CapGPMC, DigitalLogical: 44},
	&embd.PinDesc{ID: "P8_13", Aliases: []string{"23", "GPIO_23", "EHRPWM2B"}, Caps: embd.CapDigital | embd.CapGPMC, DigitalLogical: 23},
	&embd.PinDesc{ID: "P8_14", Aliases: []string{"26", "GPIO_26"}, Caps: embd.CapDigital | embd.CapGPMC, DigitalLogical: 26},
	&embd.PinDesc{ID: "P8_15", Aliases: []string{"47", "GPIO_47"}, Caps: embd.CapDigital | embd.CapGPMC, DigitalLogical: 47},
	&embd.PinDesc{ID: "P8_16", Aliases: []string{"46", "GPIO_46"}, Caps: embd.CapDigital | embd.CapGPMC, DigitalLogical: 46},
	&embd.PinDesc{ID: "P8_17", Aliases: []string{"27", "GPIO_27"}, Caps: embd.CapDigital | embd.CapGPMC, DigitalLogical: 27},
	&embd.PinDesc{ID: "P8_18", Aliases: []string{"65", "GPIO_65"}, Caps: embd.CapDigital | embd.CapGPMC, DigitalLogical: 65},
	&embd.PinDesc{ID: "P8_19", Aliases: []string{"22", "GPIO_22", "EHRPWM2A"}, Caps: embd.CapDigital | embd.CapGPMC, DigitalLogical: 22},
	&embd.PinDesc{ID: "P8_26", Aliases: []string{"61", "GPIO_61"}, Caps: embd.CapDigital | embd.CapGPMC, DigitalLogical: 61},
	&embd.PinDesc{ID: "P8_27", Aliases: []string{"86", "GPIO_86"}, Caps: embd.CapDigital | embd.CapLCD, DigitalLogical: 86},
	&embd.PinDesc{ID: "P8_28", Aliases: []string{"88", "GPIO_88"}, Caps: embd.CapDigital | embd.CapLCD, DigitalLogical: 88},
	&embd.PinDesc{ID: "P8_29", Aliases: []string{"87", "GPIO_87"}, Caps: embd.CapDigital | embd.CapLCD, DigitalLogical: 87},
	&embd.PinDesc{ID: "P8_30", Aliases: []string{"89", "GPIO_89"}, Caps: embd.CapDigital | embd.CapLCD, DigitalLogical: 89},
	&embd.PinDesc{ID: "P8_31", Aliases: []string{"10", "GPIO_10", "UART5_CTSN"}, Caps: embd.CapDigital | embd.CapLCD, DigitalLogical: 10},
	&embd.PinDesc{ID: "P8_32", Aliases: []string{"11", "GPIO_11", "UART5_RTSN"}, Caps: embd.CapDigital | embd.CapLCD, DigitalLogical: 11},
	&embd.PinDesc{ID: "P8_33", Aliases: []string{"9", "GPIO_9 ", "UART4_RTSN"}, Caps: embd.CapDigital | embd.CapLCD, DigitalLogical: 9},
	&embd.PinDesc{ID: "P8_34", Aliases: []string{"81", "GPIO_81", "UART3_RTSN"}, Caps: embd.CapDigital | embd.CapLCD, DigitalLogical: 81},
	&embd.PinDesc{ID: "P8_35", Aliases: []string{"8", "GPIO_8 ", "UART4_CTSN"}, Caps: embd.CapDigital | embd.CapLCD, DigitalLogical: 8},
	&embd.PinDesc{ID: "P8_36", Aliases: []string{"80", "GPIO_80", "UART3_CTSN"}, Caps: embd.CapDigital | embd.CapLCD, DigitalLogical: 80},
	&embd.PinDesc{ID: "P8_37", Aliases: []string{"78", "GPIO_78", "UART5_TXD"}, Caps: embd.CapDigital | embd.CapLCD, DigitalLogical: 78},
	&embd.PinDesc{ID: "P8_38", Aliases: []string{"79", "GPIO_79", "UART5_RXD"}, Caps: embd.CapDigital | embd.CapLCD, DigitalLogical: 79},
	&embd.PinDesc{ID: "P8_39", Aliases: []string{"76", "GPIO_76"}, Caps: embd.CapDigital | embd.CapLCD, DigitalLogical: 76},
	&embd.PinDesc{ID: "P8_40", Aliases: []string{"77", "GPIO_77"}, Caps: embd.CapDigital | embd.CapLCD, DigitalLogical: 77},
	&embd.PinDesc{ID: "P8_41", Aliases: []string{"74", "GPIO_74"}, Caps: embd.CapDigital | embd.CapLCD, DigitalLogical: 74},
	&embd.PinDesc{ID: "P8_42", Aliases: []string{"75", "GPIO_75"}, Caps: embd.CapDigital | embd.CapLCD, DigitalLogical: 75},
	&embd.PinDesc{ID: "P8_43", Aliases: []string{"72", "GPIO_72"}, Caps: embd.CapDigital | embd.CapLCD, DigitalLogical: 72},
	&embd.PinDesc{ID: "P8_44", Aliases: []string{"73", "GPIO_73"}, Caps: embd.CapDigital | embd.CapLCD, DigitalLogical: 73},
	&embd.PinDesc{ID: "P8_45", Aliases: []string{"70", "GPIO_70"}, Caps: embd.CapDigital | embd.CapLCD, DigitalLogical: 70},
	&embd.PinDesc{ID: "P8_46", Aliases: []string{"71", "GPIO_71"}, Caps: embd.CapDigital | embd.CapLCD, DigitalLogical: 71},

	&embd.PinDesc{ID: "P9_11", Aliases: []string{"30", "GPIO_30", "UART4_RXD"}, Caps: embd.CapDigital | embd.CapUART, DigitalLogical: 30},
	&embd.PinDesc{ID: "P9_12", Aliases: []string{"60", "GPIO_60", "GPIO1_28"}, Caps: embd.CapDigital, DigitalLogical: 60},
	&embd.PinDesc{ID: "P9_13", Aliases: []string{"31", "GPIO_31", "UART4_TXD"}, Caps: embd.CapDigital | embd.CapUART, DigitalLogical: 31},
	&embd.PinDesc{ID: "P9_14", Aliases: []string{"50", "GPIO_50", "EHRPWM1A"}, Caps: embd.CapDigital | embd.CapPWM, DigitalLogical: 50},
	&embd.PinDesc{ID: "P9_15", Aliases: []string{"48", "GPIO_48", "GPIO1_16"}, Caps: embd.CapDigital, DigitalLogical: 48},
	&embd.PinDesc{ID: "P9_16", Aliases: []string{"51", "GPIO_51", "EHRPWM1B"}, Caps: embd.CapDigital | embd.CapPWM, DigitalLogical: 51},
	&embd.PinDesc{ID: "P9_17", Aliases: []string{"5", "GPIO_5", "I2C1_SCL"}, Caps: embd.CapDigital | embd.CapI2C, DigitalLogical: 5},
	&embd.PinDesc{ID: "P9_18", Aliases: []string{"4", "GPIO_4", "I2C1_SDA"}, Caps: embd.CapDigital | embd.CapI2C, DigitalLogical: 4},
	&embd.PinDesc{ID: "P9_19", Aliases: []string{"13", "GPIO_13", "I2C2_SCL"}, Caps: embd.CapDigital | embd.CapI2C, DigitalLogical: 13},
	&embd.PinDesc{ID: "P9_20", Aliases: []string{"12", "GPIO_12", "I2C2_SDA"}, Caps: embd.CapDigital | embd.CapI2C, DigitalLogical: 12},
	&embd.PinDesc{ID: "P9_21", Aliases: []string{"3", "GPIO_3", "UART2_TXD"}, Caps: embd.CapDigital | embd.CapUART, DigitalLogical: 3},
	&embd.PinDesc{ID: "P9_22", Aliases: []string{"2", "GPIO_2", "UART2_RXD"}, Caps: embd.CapDigital | embd.CapUART, DigitalLogical: 2},
	&embd.PinDesc{ID: "P9_23", Aliases: []string{"49", "GPIO_49", "GPIO1_17"}, Caps: embd.CapDigital, DigitalLogical: 49},
	&embd.PinDesc{ID: "P9_24", Aliases: []string{"15", "GPIO_15", "UART1_TXD"}, Caps: embd.CapDigital | embd.CapUART, DigitalLogical: 15},
	&embd.PinDesc{ID: "P9_25", Aliases: []string{"117", "GPIO_117", "GPIO3_21"}, Caps: embd.CapDigital, DigitalLogical: 117},
	&embd.PinDesc{ID: "P9_26", Aliases: []string{"14", "GPIO_14", "UART1_RXD"}, Caps: embd.CapDigital | embd.CapUART, DigitalLogical: 14},
	&embd.PinDesc{ID: "P9_27", Aliases: []string{"115", "GPIO_115", "GPIO3_19"}, Caps: embd.CapDigital, DigitalLogical: 115},
	&embd.PinDesc{ID: "P9_28", Aliases: []string{"113", "GPIO_113", "SPI1_CS0"}, Caps: embd.CapDigital | embd.CapSPI, DigitalLogical: 113},
	&embd.PinDesc{ID: "P9_29", Aliases: []string{"111", "GPIO_111", "SPI1_D0"}, Caps: embd.CapDigital | embd.CapSPI, DigitalLogical: 111},
	&embd.PinDesc{ID: "P9_30", Aliases: []string{"112", "GPIO_112", "SPI1_D1"}, Caps: embd.CapDigital | embd.CapSPI, DigitalLogical: 112},
	&embd.PinDesc{ID: "P9_31", Aliases: []string{"110", "GPIO_110", "SPI1_SCLK"}, Caps: embd.CapDigital | embd.CapSPI, DigitalLogical: 110},
	&embd.PinDesc{ID: "P9_32", Aliases: []string{"VADC"}},
	&embd.PinDesc{ID: "P9_33", Aliases: []string{"4", "AIN4"}, Caps: embd.CapAnalog, AnalogLogical: 4},
	&embd.PinDesc{ID: "P9_34", Aliases: []string{"AGND"}},
	&embd.PinDesc{ID: "P9_35", Aliases: []string{"6", "AIN6"}, Caps: embd.CapAnalog, AnalogLogical: 6},
	&embd.PinDesc{ID: "P9_36", Aliases: []string{"5", "AIN5"}, Caps: embd.CapAnalog, AnalogLogical: 5},
	&embd.PinDesc{ID: "P9_37", Aliases: []string{"2", "AIN2"}, Caps: embd.CapAnalog, AnalogLogical: 2},
	&embd.PinDesc{ID: "P9_38", Aliases: []string{"3", "AIN3"}, Caps: embd.CapAnalog, AnalogLogical: 3},
	&embd.PinDesc{ID: "P9_39", Aliases: []string{"0", "AIN0"}, Caps: embd.CapAnalog, AnalogLogical: 0},
	&embd.PinDesc{ID: "P9_40", Aliases: []string{"1", "AIN1"}, Caps: embd.CapAnalog, AnalogLogical: 1},
}

var ledMap = embd.LEDMap{
	"beaglebone:green:usr0": []string{"0", "USR0", "usr0"},
	"beaglebone:green:usr1": []string{"1", "USR1", "usr1"},
	"beaglebone:green:usr2": []string{"2", "USR2", "usr2"},
	"beaglebone:green:usr3": []string{"3", "USR3", "usr3"},
}

var spiDeviceMinor int = 1

func ensureFeatureEnabled(id string) error {
	glog.V(3).Infof("bbb: enabling feature %v", id)
	pattern := "/sys/devices/bone_capemgr.*/slots"
	file, err := embd.FindFirstMatchingFile(pattern)
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	str := string(bytes)
	if strings.Contains(str, id) {
		glog.V(3).Infof("bbb: feature %v already enabled", id)
		return nil
	}
	slots, err := os.OpenFile(file, os.O_WRONLY, os.ModeExclusive)
	if err != nil {
		return err
	}
	defer slots.Close()
	glog.V(3).Infof("bbb: writing %v to slots file", id)
	_, err = slots.WriteString(id)
	return err
}

// This cannot be currently used to disable things like the
// analog and pwm modules. Removing them from slots file can
// potentially cause a kernel panic and unsettle things. So the
// recommended thing to do is to simply reboot.
func ensureFeatureDisabled(id string) error {
	pattern := "/sys/devices/bone_capemgr.*/slots"
	file, err := embd.FindFirstMatchingFile(pattern)
	if err != nil {
		return err
	}
	slots, err := os.OpenFile(file, os.O_RDWR, os.ModeExclusive)
	if err != nil {
		return err
	}
	defer slots.Close()
	scanner := bufio.NewScanner(slots)
	for scanner.Scan() {
		text := scanner.Text()
		if !strings.Contains(text, id) {
			continue
		}
		// Extract the id from the line
		idx := strings.Index(text, ":")
		if idx < 0 {
			// Something is off, bail
			continue
		}
		dis := strings.TrimSpace(text[:idx])
		slots.Seek(0, 0)
		_, err = slots.WriteString("-" + dis)
		return err
	}
	// Could not disable the feature
	return fmt.Errorf("embd: could not disable feature %q", id)
}

func spiInitializer() error {
	if err := ensureFeatureEnabled("BB-SPIDEV0"); err != nil {
		return err
	}
	return nil
}

func init() {
	embd.Register(embd.HostBBB, func(rev int) *embd.Descriptor {
		return &embd.Descriptor{
			GPIODriver: func() embd.GPIODriver {
				return embd.NewGPIODriver(pins, generic.NewDigitalPin, newAnalogPin, newPWMPin)
			},
			I2CDriver: func() embd.I2CDriver {
				return embd.NewI2CDriver(generic.NewI2CBus)
			},
			LEDDriver: func() embd.LEDDriver {
				return embd.NewLEDDriver(ledMap, generic.NewLED)
			},
			SPIDriver: func() embd.SPIDriver {
				return embd.NewSPIDriver(spiDeviceMinor, generic.NewSPIBus, spiInitializer)
			},
		}
	})
}
