// PWM support on the BBB.

package bbb

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/golang/glog"
	"github.com/kidoman/embd"
	"github.com/kidoman/embd/util"
)

const (
	// PWMDefaultPolarity represents the default polarity (Positve or 1) for pwm.
	PWMDefaultPolarity = embd.Positive

	// PWMDefaultDuty represents the default duty (0ns) for pwm.
	PWMDefaultDuty = 0

	// PWMDefaultPeriod represents the default period (500000ns) for pwm. Equals 2000 Hz.
	PWMDefaultPeriod = 500000

	// PWMMaxPulseWidth represents the max period (1000000000ns) supported by pwm. Equals 1 Hz.
	PWMMaxPulseWidth = 1000000000
)

type pwmPin struct {
	n string

	drv embd.GPIODriver

	period   int
	polarity embd.Polarity

	dutyf     *os.File
	periodf   *os.File
	polarityf *os.File

	initialized bool
}

func newPWMPin(pd *embd.PinDesc, drv embd.GPIODriver) embd.PWMPin {
	return &pwmPin{n: pd.ID, drv: drv}
}

func (p *pwmPin) N() string {
	return p.n
}

func (p *pwmPin) id() string {
	return "bone_pwm_" + p.n
}

func (p *pwmPin) init() error {
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
	if p.periodf, err = p.periodFile(basePath); err != nil {
		return err
	}
	if p.dutyf, err = p.dutyFile(basePath); err != nil {
		return err
	}
	if p.polarityf, err = p.polarityFile(basePath); err != nil {
		return err
	}

	p.initialized = true

	if err := p.reset(); err != nil {
		return err
	}

	return nil
}

func (p *pwmPin) ensurePWMEnabled() error {
	return ensureFeatureEnabled("am33xx_pwm")
}

func (p *pwmPin) ensurePinEnabled() error {
	return ensureFeatureEnabled(p.id())
}

func (p *pwmPin) ensurePinDisabled() error {
	return ensureFeatureDisabled(p.id())
}

func (p *pwmPin) basePath() (string, error) {
	pattern := "/sys/devices/ocp.*/pwm_test_" + p.n + ".*"
	return embd.FindFirstMatchingFile(pattern)
}

func (p *pwmPin) openFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_WRONLY, os.ModeExclusive)
}

func (p *pwmPin) ensurePeriodFileExists(basePath string, d time.Duration) error {
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

func (p *pwmPin) periodFilePath(basePath string) string {
	return path.Join(basePath, "period")
}

func (p *pwmPin) periodFile(basePath string) (*os.File, error) {
	return p.openFile(p.periodFilePath(basePath))
}

func (p *pwmPin) dutyFile(basePath string) (*os.File, error) {
	return p.openFile(path.Join(basePath, "duty"))
}

func (p *pwmPin) polarityFile(basePath string) (*os.File, error) {
	return p.openFile(path.Join(basePath, "polarity"))
}

func (p *pwmPin) SetPeriod(ns int) error {
	if err := p.init(); err != nil {
		return err
	}

	if ns > PWMMaxPulseWidth {
		return fmt.Errorf("embd: pwm period for %v is out of bounds (must be =< %vns)", p.n, PWMMaxPulseWidth)
	}

	_, err := p.periodf.WriteString(strconv.Itoa(ns))
	if err != nil {
		return err
	}

	p.period = ns

	return nil
}

func (p *pwmPin) SetDuty(ns int) error {
	if err := p.init(); err != nil {
		return err
	}

	if ns > PWMMaxPulseWidth {
		return fmt.Errorf("embd: pwm duty %v for pin %v is out of bounds (must be =< %vns)", p.n, PWMMaxPulseWidth)
	}

	if ns > p.period {
		return fmt.Errorf("embd: pwm duty %v for pin %v is greater than the period %v", ns, p.n, p.period)
	}

	_, err := p.dutyf.WriteString(strconv.Itoa(ns))
	if err != nil {
		return err
	}

	return nil
}

func (p *pwmPin) SetMicroseconds(us int) error {
	if err := p.init(); err != nil {
		return err
	}

	if p.period != 20000000 {
		glog.Warningf("embd: pwm pin %v has freq %v hz. recommended 50 hz for servo mode", 1000000000/p.period)
	}
	duty := us * 1000 // in nanoseconds
	if duty > p.period {
		return fmt.Errorf("embd: calculated pwm duty %vns for pin %v (servo mode) is greater than the period %vns", duty, p.n, p.period)
	}
	return p.SetDuty(duty)
}

func (p *pwmPin) SetAnalog(value byte) error {
	duty := util.Map(int64(value), 0, 255, 0, int64(p.period))
	return p.SetDuty(int(duty))
}

func (p *pwmPin) SetPolarity(pol embd.Polarity) error {
	if err := p.init(); err != nil {
		return err
	}

	_, err := p.polarityf.WriteString(strconv.Itoa(int(pol)))
	if err != nil {
		return err
	}

	p.polarity = pol

	return nil
}

func (p *pwmPin) reset() error {
	if err := p.SetPolarity(embd.Positive); err != nil {
		return err
	}
	if err := p.SetDuty(PWMDefaultDuty); err != nil {
		return err
	}
	if err := p.SetPeriod(PWMDefaultPeriod); err != nil {
		return err
	}

	return nil
}

func (p *pwmPin) Close() error {
	if err := p.drv.Unregister(p.n); err != nil {
		return err
	}

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
