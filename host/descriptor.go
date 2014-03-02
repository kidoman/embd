package host

import "fmt"

type Descriptor struct {
	GPIO func() interface{}
	I2C  func() interface{}
}

type Describer func(rev int) *Descriptor

var Describers = map[Host]Describer{}

func Describe() (*Descriptor, error) {
	host, rev, err := Detect()
	if err != nil {
		return nil, err
	}

	describer, ok := Describers[host]
	if !ok {
		return nil, fmt.Errorf("host: invalid host %q", host)
	}

	return describer(rev), nil
}
