package embd

import (
	"errors"
	"fmt"
)

type Descriptor struct {
	GPIO      func() GPIO
	I2C       func() I2C
	LEDDriver func() LEDDriver
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

var ErrFeatureNotSupport = errors.New("embd: feature is not supported")
