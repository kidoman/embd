// Host descriptor data structures.

package embd

import (
	"errors"
	"fmt"
)

// Descriptor represents a host descriptor.
type Descriptor struct {
	GPIODriver func() GPIODriver
	I2CDriver  func() I2CDriver
	LEDDriver  func() LEDDriver
}

// The Describer type is a Descriptor provider.
type Describer func(rev int) *Descriptor

// Describers is a global list of registered host Describers.
var Describers = map[Host]Describer{}

// DescribeHost returns the detected host descriptor.
func DescribeHost() (*Descriptor, error) {
	host, rev, err := DetectHost()
	if err != nil {
		return nil, err
	}

	describer, ok := Describers[host]
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
