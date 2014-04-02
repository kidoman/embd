// Package servoblaster allows interfacing with the software servoblaster driver.
//
// More details on ServoBlaster at: https://github.com/richardghirst/PiBits/tree/master/ServoBlaster
package servoblaster

import (
	"fmt"
	"os"

	"github.com/golang/glog"
)

// ServoBlaster represents a software RPi PWM/PCM based servo control module.
type ServoBlaster struct {
	initialized bool
	fd          *os.File
}

// New creates a new ServoBlaster instance.
func New() *ServoBlaster {
	return &ServoBlaster{}
}

func (d *ServoBlaster) setup() error {
	if d.initialized {
		return nil
	}
	var err error
	if d.fd, err = os.OpenFile("/dev/servoblaster", os.O_WRONLY, os.ModeExclusive); err != nil {
		return err
	}
	d.initialized = true
	return nil
}

type pwmChannel struct {
	d *ServoBlaster

	channel int
}

func (p *pwmChannel) SetMicroseconds(us int) error {
	return p.d.setMicroseconds(p.channel, us)
}

func (d *ServoBlaster) Channel(channel int) *pwmChannel {
	return &pwmChannel{d: d, channel: channel}
}

// SetMicroseconds sends a command to the PWM driver to generate a us wide pulse.
func (d *ServoBlaster) setMicroseconds(channel, us int) error {
	if err := d.setup(); err != nil {
		return err
	}
	cmd := fmt.Sprintf("%v=%vus\n", channel, us)
	glog.V(1).Infof("servoblaster: sending command %q", cmd)
	_, err := d.fd.WriteString(cmd)
	return err
}

// Close closes the open driver handle.
func (d *ServoBlaster) Close() error {
	if d.fd != nil {
		return d.fd.Close()
	}
	return nil
}
