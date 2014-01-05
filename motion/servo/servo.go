// Package servo allows control of servos using a PWM controller.
package servo

import (
	"log"
)

// A PWM interface implements access to a pwm controller.
type PWM interface {
	SetPwm(n int, onTime int, offTime int) error
}

// A SERVO interface implements access to the servo.
type Servo interface {
	SetAngle(angle int) error

	SetDebug(v bool)
}

type servo struct {
	pwm  PWM
	freq int
	n    int

	minms, maxms float64

	debug bool
}

// New creates a new SERVO interface.
// pwm: instance of a PWM controller.
// freq: Frequency of pwm signal (typically 50Hz)
// n: PWM channel of the pwm controller to be used
// minms: Pulse width corresponding to servo position of 0deg
// maxms: Pulse width corresponding to servo position of 180deg
func New(pwm PWM, freq int, n int, minms, maxms float64) Servo {
	return &servo{
		pwm:   pwm,
		freq:  freq,
		n:     n,
		minms: minms,
		maxms: maxms,
	}
}

// SetDebug is used to enable logging (debug mode).
func (s *servo) SetDebug(v bool) {
	s.debug = v
}

// SetAngle sets the servo angle.
func (s *servo) SetAngle(angle int) error {
	ms := s.minms + float64(angle)*(s.maxms-s.minms)/180
	offTime := int(ms * float64(s.freq) * 4096 / 1000)

	if s.debug {
		log.Printf("servo: given angle %v calculated ms %2.2f offTime %v", angle, ms, offTime)
	}

	return s.pwm.SetPwm(s.n, 0, offTime)
}
