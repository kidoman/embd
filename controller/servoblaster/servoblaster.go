// Package servoblaster allows interfacing with the software servoblaster driver.
//
// More details on ServoBlaster at: https://github.com/richardghirst/PiBits/tree/master/ServoBlaster
package servoblaster

import (
	"fmt"
	"log"
	"os"
)

// ServoBlaster represents a software RPi PWM/PCM based servo control module.
type ServoBlaster struct {
	initialized bool
	fd          *os.File

	// Debug level.
	Debug bool
}

// New creates a new ServoBlaster instance.
func New() *ServoBlaster {
	return &ServoBlaster{}
}

func (d *ServoBlaster) setup() (err error) {
	if d.initialized {
		return
	}
	if d.fd, err = os.OpenFile("/dev/servoblaster", os.O_WRONLY, os.ModeExclusive); err != nil {
		return
	}
	d.initialized = true
	return
}

// SetMicroseconds sends a command to the PWM driver to generate a us wide pulse.
func (d *ServoBlaster) SetMicroseconds(channel, us int) (err error) {
	if err = d.setup(); err != nil {
		return
	}
	cmd := fmt.Sprintf("%v=%vus\n", channel, us)
	if d.Debug {
		log.Printf("servoblaster: sending command %q", cmd)
	}
	_, err = d.fd.WriteString(cmd)
	return
}

// Close closes the open driver handle.
func (d *ServoBlaster) Close() (err error) {
	if d.fd != nil {
		err = d.fd.Close()
	}
	return
}
