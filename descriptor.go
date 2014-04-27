// Host descriptor data structures.

package embd

import (
	"errors"
	"fmt"

	"github.com/golang/glog"
)

// Descriptor represents a host descriptor.
type Descriptor struct {
	GPIODriver func() GPIODriver
	I2CDriver  func() I2CDriver
	LEDDriver  func() LEDDriver
	SPIDriver  func() SPIDriver
}

// The Describer type is a Descriptor provider.
type Describer func(rev int) *Descriptor

// Describers is a global list of registered host Describers.
var describers = make(map[Host]Describer)

// Register makes a host describer available by the provided host key.
// If Register is called twice with the same host or if describer is nil,
// it panics.
func Register(host Host, describer Describer) {
	if describer == nil {
		panic("embd: describer is nil")
	}
	if _, dup := describers[host]; dup {
		panic("embd: describer already registered")
	}
	describers[host] = describer

	glog.V(1).Infof("embd: host %v is registered", host)
}

var hostOverride Host
var hostRevOverride int
var hostOverriden bool

// SetHost overrides the host and revision no.
func SetHost(host Host, rev int) {
	hostOverride = host
	hostRevOverride = rev

	hostOverriden = true
}

// DescribeHost returns the detected host descriptor.
// Can be overriden by calling SetHost though.
func DescribeHost() (*Descriptor, error) {
	var host Host
	var rev int

	if hostOverriden {
		host, rev = hostOverride, hostRevOverride
	} else {
		var err error
		host, rev, err = DetectHost()
		if err != nil {
			return nil, err
		}
	}

	describer, ok := describers[host]
	if !ok {
		return nil, fmt.Errorf("host: invalid host %q", host)
	}

	return describer(rev), nil
}

// ErrFeatureNotSupported is returned when the host does not support a
// particular feature.
var ErrFeatureNotSupported = errors.New("embd: requested feature is not supported")

// ErrFeatureNotImplemented is returned when a particular feature is supported
// by the host but not implemented yet.
var ErrFeatureNotImplemented = errors.New("embd: requested feature is not implemented")
