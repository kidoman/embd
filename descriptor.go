package embd

import (
	"errors"
	"fmt"
)

type Descriptor struct {
	GPIODriver func() GPIODriver
	I2CDriver  func() I2CDriver
	LEDDriver  func() LEDDriver
}

type Describer func(rev int) *Descriptor

var Describers = map[Host]Describer{}

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

var ErrFeatureNotSupported = errors.New("embd: requested feature is not supported")
var ErrFeatureNotImplemented = errors.New("embd: requested feature is not implemented")
