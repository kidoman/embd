package embd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type LEDMap map[string][]string

type ledDriver struct {
	ledMap LEDMap

	initializedLEDs map[string]LED
}

func newLEDDriver(ledMap LEDMap) LEDDriver {
	return &ledDriver{
		ledMap:          ledMap,
		initializedLEDs: map[string]LED{},
	}
}

func (d *ledDriver) lookup(k interface{}) (string, error) {
	var ks string
	switch key := k.(type) {
	case int:
		ks = strconv.Itoa(key)
	case string:
		ks = key
	case fmt.Stringer:
		ks = key.String()
	default:
		return "", errors.New("led: invalid key type")
	}

	for id := range d.ledMap {
		for _, alias := range d.ledMap[id] {
			if alias == ks {
				return id, nil
			}
		}
	}

	return "", fmt.Errorf("led: no match found for %q", k)
}

func (d *ledDriver) LED(k interface{}) (LED, error) {
	id, err := d.lookup(k)
	if err != nil {
		return nil, err
	}

	led := newLED(id)
	d.initializedLEDs[id] = led

	return led, nil
}

func (d *ledDriver) Close() error {
	for _, led := range d.initializedLEDs {
		if err := led.Close(); err != nil {
			return err
		}
	}

	return nil
}

type led struct {
	id string

	brightness *os.File

	initialized bool
}

func newLED(id string) LED {
	return &led{id: id}
}

func (l *led) init() error {
	if l.initialized {
		return nil
	}

	var err error
	if l.brightness, err = l.brightnessFile(); err != nil {
		return err
	}

	l.initialized = true

	return nil
}

func (l *led) brightnessFilePath() string {
	return fmt.Sprintf("/sys/class/leds/%v/brightness", l.id)
}

func (l *led) openFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDWR, os.ModeExclusive)
}

func (l *led) brightnessFile() (*os.File, error) {
	return l.openFile(l.brightnessFilePath())
}

func (l *led) On() error {
	if err := l.init(); err != nil {
		return err
	}

	_, err := l.brightness.WriteString("1")
	return err
}

func (l *led) Off() error {
	if err := l.init(); err != nil {
		return err
	}

	_, err := l.brightness.WriteString("0")
	return err
}

func (l *led) isOn() (bool, error) {
	l.brightness.Seek(0, 0)
	bytes, err := ioutil.ReadAll(l.brightness)
	if err != nil {
		return false, err
	}
	str := string(bytes)
	str = strings.TrimSpace(str)
	if str == "1" {
		return true, nil
	}
	return false, nil
}

func (l *led) Toggle() error {
	if err := l.init(); err != nil {
		return err
	}

	state, err := l.isOn()
	if err != nil {
		return err
	}

	if state {
		return l.Off()
	}
	return l.On()
}

func (l *led) Close() error {
	if !l.initialized {
		return nil
	}

	if err := l.brightness.Close(); err != nil {
		return err
	}

	l.initialized = false

	return nil
}
