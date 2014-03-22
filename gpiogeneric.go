package embd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/golang/glog"
)

const (
	CapNormal int = 1 << iota
	CapI2C
	CapUART
	CapSPI
	CapGPMC
	CapLCD
	CapPWM
)

type PinDesc struct {
	N    int
	IDs  []string
	Caps int
}

type PinMap []*PinDesc

func (m PinMap) Lookup(k interface{}) (*PinDesc, bool) {
	switch key := k.(type) {
	case int:
		for i := range m {
			if m[i].N == key {
				return m[i], true
			}
		}
	case string:
		for i := range m {
			for j := range m[i].IDs {
				if m[i].IDs[j] == key {
					return m[i], true
				}
			}
		}
	}

	return nil, false
}

type gpioDriver struct {
	exporter, unexporter *os.File

	initialized bool

	pinMap          PinMap
	initializedPins map[int]*digitalPin
}

func newGPIODriver(pinMap PinMap) *gpioDriver {
	return &gpioDriver{
		pinMap:          pinMap,
		initializedPins: map[int]*digitalPin{},
	}
}

func (io *gpioDriver) init() error {
	if io.initialized {
		return nil
	}

	var err error
	if io.exporter, err = os.OpenFile("/sys/class/gpio/export", os.O_WRONLY, os.ModeExclusive); err != nil {
		return err
	}
	if io.unexporter, err = os.OpenFile("/sys/class/gpio/unexport", os.O_WRONLY, os.ModeExclusive); err != nil {
		return err
	}

	io.initialized = true

	return nil
}

func (io *gpioDriver) lookupKey(key interface{}) (*PinDesc, bool) {
	return io.pinMap.Lookup(key)
}

func (io *gpioDriver) export(n int) error {
	_, err := io.exporter.WriteString(strconv.Itoa(n))
	return err
}

func (io *gpioDriver) unexport(n int) error {
	_, err := io.unexporter.WriteString(strconv.Itoa(n))
	return err
}

func (io *gpioDriver) digitalPin(key interface{}) (*digitalPin, error) {
	pd, found := io.lookupKey(key)
	if !found {
		err := fmt.Errorf("gpio: could not find pin matching %q", key)
		return nil, err
	}

	n := pd.N

	p, ok := io.initializedPins[n]
	if ok {
		return p, nil
	}

	if pd.Caps&CapNormal == 0 {
		err := fmt.Errorf("gpio: sorry, pin %q cannot be used for GPIO", key)
		return nil, err
	}

	if pd.Caps != CapNormal {
		glog.Infof("gpio: pin %q is not a dedicated GPIO pin. please refer to the system reference manual for more details", key)
	}

	if err := io.export(n); err != nil {
		return nil, err
	}

	p, err := newDigitalPin(n)
	if err != nil {
		io.unexport(n)
		return nil, err
	}

	io.initializedPins[n] = p

	return p, nil
}

func (io *gpioDriver) DigitalPin(key interface{}) (DigitalPin, error) {
	if err := io.init(); err != nil {
		return nil, err
	}

	return io.digitalPin(key)
}

func (io *gpioDriver) Close() error {
	for n := range io.initializedPins {
		io.unexport(n)
	}

	io.exporter.Close()
	io.unexporter.Close()

	return nil
}
