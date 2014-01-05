// Package servo allows control of servos using a PWM controller.
package servo

import (
	"log"

	"github.com/kid0m4n/go-rpi/util"
)

// A PWM interface implements access to a pwm controller.
type PWM interface {
	SetPwm(channel int, onTime int, offTime int) error
}

type Servo struct {
	PWM     PWM
	Freq    int
	Channel int

	Minms, Maxms float64

	Debug bool
}

// New creates a new Servo interface.
// pwm: instance of a PWM controller.
// freq: Frequency of pwm signal (typically 50Hz)
// channel: PWM channel of the pwm controller to be used
// minms: Pulse width corresponding to servo position of 0deg
// maxms: Pulse width corresponding to servo position of 180deg
func New(pwm PWM, freq, channel int, minms, maxms float64) *Servo {
	return &Servo{
		PWM:     pwm,
		Freq:    freq,
		Channel: channel,
		Minms:   minms,
		Maxms:   maxms,
	}
}

// SetAngle sets the servo angle.
func (s *Servo) SetAngle(angle int) error {
	us := util.Map(int64(angle), 0, 180, int64(s.Minms*1000), int64(s.Maxms*1000))
	offTime := int(us) * s.Freq * 4096 / 1000000

	if s.Debug {
		log.Printf("servo: given angle %v calculated %v us offTime %v", angle, us, offTime)
	}

	return s.PWM.SetPwm(s.Channel, 0, offTime)
}
