package gpio

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/golang/glog"
)

type Direction int
type State int
type Pull int

const (
	Input Direction = iota
	Output
)

const (
	Low State = iota
	High
)

const (
	PullOff Pull = iota
	PullDown
	PullUp
)

const (
	normal = 1 << iota
	i2c
	spi
	uart
)

type gpio struct {
	exporter, unexporter *os.File

	initialized bool

	pinMap          pinMap
	initializedPins map[int]*pin
}

var instance *gpio
var instanceLock sync.Mutex

var Default *gpio

func New() *gpio {
	instanceLock.Lock()
	defer instanceLock.Unlock()

	if instance == nil {
		instance = &gpio{
			pinMap:          rev2Pins,
			initializedPins: map[int]*pin{},
		}
	}

	return instance
}

func (io *gpio) init() (err error) {
	if io.initialized {
		return
	}

	if io.exporter, err = os.OpenFile("/sys/class/gpio/export", os.O_WRONLY, os.ModeExclusive); err != nil {
		return
	}
	if io.unexporter, err = os.OpenFile("/sys/class/gpio/unexport", os.O_WRONLY, os.ModeExclusive); err != nil {
		return
	}

	io.initialized = true

	return
}

func (io *gpio) lookupKey(key interface{}) (*pinDesc, bool) {
	return io.pinMap.lookup(key)
}

func (io *gpio) export(n int) (err error) {
	_, err = io.exporter.WriteString(strconv.Itoa(n))
	return
}

func (io *gpio) unexport(n int) (err error) {
	_, err = io.unexporter.WriteString(strconv.Itoa(n))
	return
}

func (io *gpio) pin(key interface{}) (pin *pin, err error) {
	pd, found := io.lookupKey(key)
	if !found {
		err = fmt.Errorf("gpio: could not find pin matching %q", key)
		return
	}

	n := pd.n

	var ok bool
	if pin, ok = io.initializedPins[n]; ok {
		return
	}

	if pd.caps&normal == 0 {
		err = fmt.Errorf("gpio: sorry, pin %q cannot be used for GPIO", key)
		return
	}

	if pd.caps != normal {
		glog.Infof("gpio: pin %q is not a dedicated GPIO pin. please refer to the system reference manual for more details", key)
	}

	if err = io.export(n); err != nil {
		return
	}

	if pin, err = NewPin(n); err != nil {
		io.unexport(n)
		return
	}

	io.initializedPins[n] = pin

	return
}

func (io *gpio) Pin(key interface{}) (pin *pin, err error) {
	if err = io.init(); err != nil {
		return
	}

	return io.pin(key)
}

func (io *gpio) Mode(key interface{}, dir Direction) (err error) {
	if err = io.init(); err != nil {
		return
	}

	var pin *pin
	if pin, err = io.pin(key); err != nil {
		return
	}

	return pin.Mode(dir)
}

func (io *gpio) Input(key interface{}) error {
	return io.Mode(key, Input)
}

func (io *gpio) Output(key interface{}) error {
	return io.Mode(key, Output)
}

func (io *gpio) Read(key interface{}) (state State, err error) {
	if err = io.init(); err != nil {
		return
	}

	var pin *pin
	if pin, err = io.pin(key); err != nil {
		return
	}

	return pin.Read()
}

func (io *gpio) Write(key interface{}, state State) (err error) {
	if err = io.init(); err != nil {
		return
	}

	var pin *pin
	if pin, err = io.pin(key); err != nil {
		return
	}

	return pin.Write(state)
}

func (io *gpio) Low(key interface{}) error {
	return io.Write(key, Low)
}

func (io *gpio) High(key interface{}) error {
	return io.Write(key, High)
}

func (io *gpio) SetActiveLow(key interface{}, b bool) (err error) {
	if err = io.init(); err != nil {
		return
	}

	var pin *pin
	if pin, err = io.pin(key); err != nil {
		return
	}

	return pin.SetActiveLow(b)
}

func (io *gpio) ActiveLow(key interface{}) error {
	return io.SetActiveLow(key, true)
}

func (io *gpio) ActiveHigh(key interface{}) error {
	return io.SetActiveLow(key, false)
}

func (io *gpio) Close() {
	for n := range io.initializedPins {
		io.unexport(n)
	}

	io.exporter.Close()
	io.unexporter.Close()
}
