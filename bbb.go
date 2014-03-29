// BeagleBone Black support.
// The following features are supported on Linux kernel 3.8+
//
//	GPIO (digital (rw), analog (ro), pwm)
//	I2C
//	LED

package embd

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var bbbPins = PinMap{
	&PinDesc{ID: "P8_07", Aliases: []string{"66", "GPIO_66", "Caps: TIMER4"}, Caps: CapDigital | CapGPMC, DigitalLogical: 66},
	&PinDesc{ID: "P8_08", Aliases: []string{"67", "GPIO_67", "TIMER7"}, Caps: CapDigital | CapGPMC, DigitalLogical: 67},
	&PinDesc{ID: "P8_09", Aliases: []string{"69", "GPIO_69", "TIMER5"}, Caps: CapDigital | CapGPMC, DigitalLogical: 69},
	&PinDesc{ID: "P8_10", Aliases: []string{"68", "GPIO_68", "TIMER6"}, Caps: CapDigital | CapGPMC, DigitalLogical: 68},
	&PinDesc{ID: "P8_11", Aliases: []string{"45", "GPIO_45"}, Caps: CapDigital | CapGPMC, DigitalLogical: 45},
	&PinDesc{ID: "P8_12", Aliases: []string{"44", "GPIO_44"}, Caps: CapDigital | CapGPMC, DigitalLogical: 44},
	&PinDesc{ID: "P8_13", Aliases: []string{"23", "GPIO_23", "EHRPWM2B"}, Caps: CapDigital | CapGPMC, DigitalLogical: 23},
	&PinDesc{ID: "P8_14", Aliases: []string{"26", "GPIO_26"}, Caps: CapDigital | CapGPMC, DigitalLogical: 26},
	&PinDesc{ID: "P8_15", Aliases: []string{"47", "GPIO_47"}, Caps: CapDigital | CapGPMC, DigitalLogical: 47},
	&PinDesc{ID: "P8_16", Aliases: []string{"46", "GPIO_46"}, Caps: CapDigital | CapGPMC, DigitalLogical: 46},
	&PinDesc{ID: "P8_17", Aliases: []string{"27", "GPIO_27"}, Caps: CapDigital | CapGPMC, DigitalLogical: 27},
	&PinDesc{ID: "P8_18", Aliases: []string{"65", "GPIO_65"}, Caps: CapDigital | CapGPMC, DigitalLogical: 65},
	&PinDesc{ID: "P8_19", Aliases: []string{"22", "GPIO_22", "EHRPWM2A"}, Caps: CapDigital | CapGPMC, DigitalLogical: 22},
	&PinDesc{ID: "P8_26", Aliases: []string{"61", "GPIO_61"}, Caps: CapDigital | CapGPMC, DigitalLogical: 61},
	&PinDesc{ID: "P8_27", Aliases: []string{"86", "GPIO_86"}, Caps: CapDigital | CapLCD, DigitalLogical: 86},
	&PinDesc{ID: "P8_28", Aliases: []string{"88", "GPIO_88"}, Caps: CapDigital | CapLCD, DigitalLogical: 88},
	&PinDesc{ID: "P8_29", Aliases: []string{"87", "GPIO_87"}, Caps: CapDigital | CapLCD, DigitalLogical: 87},
	&PinDesc{ID: "P8_30", Aliases: []string{"89", "GPIO_89"}, Caps: CapDigital | CapLCD, DigitalLogical: 89},
	&PinDesc{ID: "P8_31", Aliases: []string{"10", "GPIO_10", "UART5_CTSN"}, Caps: CapDigital | CapLCD, DigitalLogical: 10},
	&PinDesc{ID: "P8_32", Aliases: []string{"11", "GPIO_11", "UART5_RTSN"}, Caps: CapDigital | CapLCD, DigitalLogical: 11},
	&PinDesc{ID: "P8_33", Aliases: []string{"9", "GPIO_9 ", "UART4_RTSN"}, Caps: CapDigital | CapLCD, DigitalLogical: 9},
	&PinDesc{ID: "P8_34", Aliases: []string{"81", "GPIO_81", "UART3_RTSN"}, Caps: CapDigital | CapLCD, DigitalLogical: 81},
	&PinDesc{ID: "P8_35", Aliases: []string{"8", "GPIO_8 ", "UART4_CTSN"}, Caps: CapDigital | CapLCD, DigitalLogical: 8},
	&PinDesc{ID: "P8_36", Aliases: []string{"80", "GPIO_80", "UART3_CTSN"}, Caps: CapDigital | CapLCD, DigitalLogical: 80},
	&PinDesc{ID: "P8_37", Aliases: []string{"78", "GPIO_78", "UART5_TXD"}, Caps: CapDigital | CapLCD, DigitalLogical: 78},
	&PinDesc{ID: "P8_38", Aliases: []string{"79", "GPIO_79", "UART5_RXD"}, Caps: CapDigital | CapLCD, DigitalLogical: 79},
	&PinDesc{ID: "P8_39", Aliases: []string{"76", "GPIO_76"}, Caps: CapDigital | CapLCD, DigitalLogical: 76},
	&PinDesc{ID: "P8_40", Aliases: []string{"77", "GPIO_77"}, Caps: CapDigital | CapLCD, DigitalLogical: 77},
	&PinDesc{ID: "P8_41", Aliases: []string{"74", "GPIO_74"}, Caps: CapDigital | CapLCD, DigitalLogical: 74},
	&PinDesc{ID: "P8_42", Aliases: []string{"75", "GPIO_75"}, Caps: CapDigital | CapLCD, DigitalLogical: 75},
	&PinDesc{ID: "P8_43", Aliases: []string{"72", "GPIO_72"}, Caps: CapDigital | CapLCD, DigitalLogical: 72},
	&PinDesc{ID: "P8_44", Aliases: []string{"73", "GPIO_73"}, Caps: CapDigital | CapLCD, DigitalLogical: 73},
	&PinDesc{ID: "P8_45", Aliases: []string{"70", "GPIO_70"}, Caps: CapDigital | CapLCD, DigitalLogical: 70},
	&PinDesc{ID: "P8_46", Aliases: []string{"71", "GPIO_71"}, Caps: CapDigital | CapLCD, DigitalLogical: 71},

	&PinDesc{ID: "P9_11", Aliases: []string{"30", "GPIO_30", "UART4_RXD"}, Caps: CapDigital | CapUART, DigitalLogical: 30},
	&PinDesc{ID: "P9_12", Aliases: []string{"60", "GPIO_60", "GPIO1_28"}, Caps: CapDigital, DigitalLogical: 60},
	&PinDesc{ID: "P9_13", Aliases: []string{"31", "GPIO_31", "UART4_TXD"}, Caps: CapDigital | CapUART, DigitalLogical: 31},
	&PinDesc{ID: "P9_14", Aliases: []string{"50", "GPIO_50", "EHRPWM1A"}, Caps: CapDigital | CapPWM, DigitalLogical: 50},
	&PinDesc{ID: "P9_15", Aliases: []string{"48", "GPIO_48", "GPIO1_16"}, Caps: CapDigital, DigitalLogical: 48},
	&PinDesc{ID: "P9_16", Aliases: []string{"51", "GPIO_51", "EHRPWM1B"}, Caps: CapDigital | CapPWM, DigitalLogical: 51},
	&PinDesc{ID: "P9_17", Aliases: []string{"5", "GPIO_5", "I2C1_SCL"}, Caps: CapDigital | CapI2C, DigitalLogical: 5},
	&PinDesc{ID: "P9_18", Aliases: []string{"4", "GPIO_4", "I2C1_SDA"}, Caps: CapDigital | CapI2C, DigitalLogical: 4},
	&PinDesc{ID: "P9_19", Aliases: []string{"13", "GPIO_13", "I2C2_SCL"}, Caps: CapDigital | CapI2C, DigitalLogical: 13},
	&PinDesc{ID: "P9_20", Aliases: []string{"12", "GPIO_12", "I2C2_SDA"}, Caps: CapDigital | CapI2C, DigitalLogical: 12},
	&PinDesc{ID: "P9_21", Aliases: []string{"3", "GPIO_3", "UART2_TXD"}, Caps: CapDigital | CapUART, DigitalLogical: 3},
	&PinDesc{ID: "P9_22", Aliases: []string{"2", "GPIO_2", "UART2_RXD"}, Caps: CapDigital | CapUART, DigitalLogical: 2},
	&PinDesc{ID: "P9_23", Aliases: []string{"49", "GPIO_49", "GPIO1_17"}, Caps: CapDigital, DigitalLogical: 49},
	&PinDesc{ID: "P9_24", Aliases: []string{"15", "GPIO_15", "UART1_TXD"}, Caps: CapDigital | CapUART, DigitalLogical: 15},
	&PinDesc{ID: "P9_25", Aliases: []string{"117", "GPIO_117", "GPIO3_21"}, Caps: CapDigital, DigitalLogical: 117},
	&PinDesc{ID: "P9_26", Aliases: []string{"14", "GPIO_14", "UART1_RXD"}, Caps: CapDigital | CapUART, DigitalLogical: 14},
	&PinDesc{ID: "P9_27", Aliases: []string{"115", "GPIO_115", "GPIO3_19"}, Caps: CapDigital, DigitalLogical: 115},
	&PinDesc{ID: "P9_28", Aliases: []string{"113", "GPIO_113", "SPI1_CS0"}, Caps: CapDigital | CapSPI, DigitalLogical: 113},
	&PinDesc{ID: "P9_29", Aliases: []string{"111", "GPIO_111", "SPI1_D0"}, Caps: CapDigital | CapSPI, DigitalLogical: 111},
	&PinDesc{ID: "P9_30", Aliases: []string{"112", "GPIO_112", "SPI1_D1"}, Caps: CapDigital | CapSPI, DigitalLogical: 112},
	&PinDesc{ID: "P9_31", Aliases: []string{"110", "GPIO_110", "SPI1_SCLK"}, Caps: CapDigital | CapSPI, DigitalLogical: 110},
	&PinDesc{ID: "P9_32", Aliases: []string{"VADC"}},
	&PinDesc{ID: "P9_33", Aliases: []string{"4", "AIN4"}, Caps: CapAnalog, AnalogLogical: 4},
	&PinDesc{ID: "P9_34", Aliases: []string{"AGND"}},
	&PinDesc{ID: "P9_35", Aliases: []string{"6", "AIN6"}, Caps: CapAnalog, AnalogLogical: 6},
	&PinDesc{ID: "P9_36", Aliases: []string{"5", "AIN5"}, Caps: CapAnalog, AnalogLogical: 5},
	&PinDesc{ID: "P9_37", Aliases: []string{"2", "AIN2"}, Caps: CapAnalog, AnalogLogical: 2},
	&PinDesc{ID: "P9_38", Aliases: []string{"3", "AIN3"}, Caps: CapAnalog, AnalogLogical: 3},
	&PinDesc{ID: "P9_39", Aliases: []string{"0", "AIN0"}, Caps: CapAnalog, AnalogLogical: 0},
	&PinDesc{ID: "P9_40", Aliases: []string{"1", "AIN1"}, Caps: CapAnalog, AnalogLogical: 1},
}

var bbbLEDMap = LEDMap{
	"beaglebone:green:usr0": []string{"0", "USR0", "usr0"},
	"beaglebone:green:usr1": []string{"1", "USR1", "usr1"},
	"beaglebone:green:usr2": []string{"2", "USR2", "usr2"},
	"beaglebone:green:usr3": []string{"3", "USR3", "usr3"},
}

func bbbEnsureFeatureEnabled(id string) error {
	pattern := "/sys/devices/bone_capemgr.*/slots"
	file, err := findFirstMatchingFile(pattern)
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	str := string(bytes)
	if strings.Contains(str, id) {
		return nil
	}
	slots, err := os.OpenFile(file, os.O_WRONLY, os.ModeExclusive)
	if err != nil {
		return err
	}
	defer slots.Close()
	_, err = slots.WriteString(id)
	return err
}

// This cannot be currently used to disable things like the
// analog and pwm modules. Removing them from slots file can
// potentially cause a kernel panic and unsettle things. So the
// recommended thing to do is to simply reboot.
func bbbEnsureFeatureDisabled(id string) error {
	pattern := "/sys/devices/bone_capemgr.*/slots"
	file, err := findFirstMatchingFile(pattern)
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

type bbbAnalogPin struct {
	n int

	val *os.File

	initialized bool
}

func newBBBAnalogPin(n int) AnalogPin {
	return &bbbAnalogPin{n: n}
}

func (p *bbbAnalogPin) N() int {
	return p.n
}

func (p *bbbAnalogPin) init() error {
	if p.initialized {
		return nil
	}

	var err error
	if err = p.ensureEnabled(); err != nil {
		return err
	}
	if p.val, err = p.valueFile(); err != nil {
		return err
	}

	p.initialized = true

	return nil
}

func (p *bbbAnalogPin) ensureEnabled() error {
	return bbbEnsureFeatureEnabled("cape-bone-iio")
}

func (p *bbbAnalogPin) valueFilePath() (string, error) {
	pattern := fmt.Sprintf("/sys/devices/ocp.*/helper.*/AIN%v", p.n)
	return findFirstMatchingFile(pattern)
}

func (p *bbbAnalogPin) openFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDONLY, os.ModeExclusive)
}

func (p *bbbAnalogPin) valueFile() (*os.File, error) {
	path, err := p.valueFilePath()
	if err != nil {
		return nil, err
	}
	return p.openFile(path)
}

func (p *bbbAnalogPin) Read() (int, error) {
	if err := p.init(); err != nil {
		return 0, err
	}

	p.val.Seek(0, 0)
	bytes, err := ioutil.ReadAll(p.val)
	if err != nil {
		return 0, err
	}
	str := string(bytes)
	str = strings.TrimSpace(str)
	return strconv.Atoi(str)
}

func (p *bbbAnalogPin) Close() error {
	if !p.initialized {
		return nil
	}

	if err := p.val.Close(); err != nil {
		return err
	}

	p.initialized = false

	return nil
}

const (
	// BBBPWMDefaultPolarity represents the default polarity (Positve or 1) for pwm.
	BBBPWMDefaultPolarity = Positive

	// BBBPWMDefaultDuty represents the default duty (0ns) for pwm.
	BBBPWMDefaultDuty = 0

	// BBBPWMDefaultPeriod represents the default period (500000ns) for pwm.
	BBBPWMDefaultPeriod = 500000

	// BBBPWMMaxPulseWidth represents the max period (1000000000ns) supported by pwm.
	BBBPWMMaxPulseWidth = 1000000000
)

type bbbPWMPin struct {
	n string

	duty     *os.File
	period   *os.File
	polarity *os.File

	initialized bool
}

func newBBBPWMPin(n string) PWMPin {
	return &bbbPWMPin{n: n}
}

func (p *bbbPWMPin) N() string {
	return p.n
}

func (p *bbbPWMPin) id() string {
	return "bone_pwm_" + p.n
}

func (p *bbbPWMPin) init() error {
	if p.initialized {
		return nil
	}

	if err := p.ensurePWMEnabled(); err != nil {
		return err
	}
	if err := p.ensurePinEnabled(); err != nil {
		return err
	}

	basePath, err := p.basePath()
	if err != nil {
		return err
	}
	if err := p.ensurePeriodFileExists(basePath, 500*time.Millisecond); err != nil {
		return err
	}
	if p.period, err = p.periodFile(basePath); err != nil {
		return err
	}
	if p.duty, err = p.dutyFile(basePath); err != nil {
		return err
	}
	if p.polarity, err = p.polarityFile(basePath); err != nil {
		return err
	}

	p.initialized = true

	return nil
}

func (p *bbbPWMPin) ensurePWMEnabled() error {
	return bbbEnsureFeatureEnabled("am33xx_pwm")
}

func (p *bbbPWMPin) ensurePinEnabled() error {
	return bbbEnsureFeatureEnabled(p.id())
}

func (p *bbbPWMPin) ensurePinDisabled() error {
	return bbbEnsureFeatureDisabled(p.id())
}

func (p *bbbPWMPin) basePath() (string, error) {
	pattern := "/sys/devices/ocp.*/pwm_test_" + p.n + ".*"
	return findFirstMatchingFile(pattern)
}

func (p *bbbPWMPin) openFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_WRONLY, os.ModeExclusive)
}

func (p *bbbPWMPin) ensurePeriodFileExists(basePath string, d time.Duration) error {
	path := p.periodFilePath(basePath)
	timeout := time.After(d)

	for {
		select {
		case <-timeout:
			return errors.New("embd: period file not found before timeout")
		default:
			if _, err := os.Stat(path); err == nil {
				return nil
			}
		}

		// We are looping, wait a bit.
		time.Sleep(10 * time.Millisecond)
	}
}

func (p *bbbPWMPin) periodFilePath(basePath string) string {
	return path.Join(basePath, "period")
}

func (p *bbbPWMPin) periodFile(basePath string) (*os.File, error) {
	return p.openFile(p.periodFilePath(basePath))
}

func (p *bbbPWMPin) dutyFile(basePath string) (*os.File, error) {
	return p.openFile(path.Join(basePath, "duty"))
}

func (p *bbbPWMPin) polarityFile(basePath string) (*os.File, error) {
	return p.openFile(path.Join(basePath, "polarity"))
}

func (p *bbbPWMPin) SetPeriod(ns int) error {
	if err := p.init(); err != nil {
		return err
	}

	if ns > BBBPWMMaxPulseWidth {
		return fmt.Errorf("embd: pwm period for %v is out of bounds (must be =< %vns)", p.n, BBBPWMMaxPulseWidth)
	}

	_, err := p.period.WriteString(strconv.Itoa(ns))
	return err
}

func (p *bbbPWMPin) SetDuty(ns int) error {
	if err := p.init(); err != nil {
		return err
	}

	if ns > BBBPWMMaxPulseWidth {
		return fmt.Errorf("embd: pwm duty for %v is out of bounds (must be =< %vns)", p.n, BBBPWMMaxPulseWidth)
	}

	_, err := p.duty.WriteString(strconv.Itoa(ns))
	return err
}

func (p *bbbPWMPin) SetPolarity(pol Polarity) error {
	if err := p.init(); err != nil {
		return err
	}

	_, err := p.polarity.WriteString(strconv.Itoa(int(pol)))
	return err
}

func (p *bbbPWMPin) reset() error {
	if err := p.SetPolarity(Positive); err != nil {
		return err
	}
	if err := p.SetDuty(BBBPWMDefaultDuty); err != nil {
		return err
	}
	if err := p.SetPeriod(BBBPWMDefaultPeriod); err != nil {
		return err
	}

	return nil
}

func (p *bbbPWMPin) Close() error {
	if !p.initialized {
		return nil
	}

	if err := p.reset(); err != nil {
		return err
	}
	if err := p.ensurePinDisabled(); err != nil {
		return err
	}

	p.initialized = false

	return nil
}

func init() {
	Register(HostBBB, func(rev int) *Descriptor {
		return &Descriptor{
			GPIODriver: func() GPIODriver {
				return newGPIODriver(bbbPins, newDigitalPin, newBBBAnalogPin, newBBBPWMPin)
			},
			I2CDriver: newI2CDriver,
			LEDDriver: func() LEDDriver {
				return newLEDDriver(bbbLEDMap)
			},
		}
	})
}
