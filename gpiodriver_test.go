package embd

import (
	"testing"
	"time"
)

type fakeDigitalPin struct {
	id string
	n  int

	drv GPIODriver
}

func (p *fakeDigitalPin) N() int {
	return p.n
}

func (*fakeDigitalPin) SetDirection(dir Direction) error {
	return nil
}

func (*fakeDigitalPin) Read() (int, error) {
	return 0, nil
}

func (*fakeDigitalPin) Write(val int) error {
	return nil
}

func (*fakeDigitalPin) TimePulse(state int) (time.Duration, error) {
	return 0, nil
}

func (*fakeDigitalPin) ActiveLow(b bool) error {
	return nil
}

func (*fakeDigitalPin) PullUp() error {
	return nil
}

func (*fakeDigitalPin) PullDown() error {
	return nil
}

func (p *fakeDigitalPin) Close() error {
	return p.drv.Unregister(p.id)
}

func (p *fakeDigitalPin) Watch(edge Edge, handler func(DigitalPin)) error {
	return nil
}

func (p *fakeDigitalPin) StopWatching() error {
	return nil
}

func newFakeDigitalPin(pd *PinDesc, drv GPIODriver) DigitalPin {
	return &fakeDigitalPin{id: pd.ID, n: pd.DigitalLogical, drv: drv}
}

func TestGpioDriverDigitalPin(t *testing.T) {
	tests := []struct {
		key interface{}
		n   int
	}{
		{1, 1},
	}
	pinMap := PinMap{
		&PinDesc{ID: "P1_1", Aliases: []string{"1"}, Caps: CapDigital, DigitalLogical: 1},
	}
	driver := NewGPIODriver(pinMap, newFakeDigitalPin, nil, nil)
	for _, test := range tests {
		pin, err := driver.DigitalPin(test.key)
		if err != nil {
			t.Errorf("Looking up %v: unexpected error: %v", test.key, err)
			continue
		}
		if pin.N() != test.n {
			t.Errorf("Looking up %v: got %v, want %v", test.key, pin.N(), test.n)
		}
	}
}

type fakeAnalogPin struct {
	id string
	n  int

	drv GPIODriver
}

func (p *fakeAnalogPin) N() int {
	return p.n
}

func (*fakeAnalogPin) Read() (int, error) {
	return 0, nil
}

func (*fakeAnalogPin) Write(val int) error {
	return nil
}

func (*fakeAnalogPin) Close() error {
	return nil
}

func newFakeAnalogPin(pd *PinDesc, drv GPIODriver) AnalogPin {
	return &fakeAnalogPin{id: pd.ID, n: pd.AnalogLogical, drv: drv}
}

func TestGpioDriverAnalogPin(t *testing.T) {
	tests := []struct {
		key interface{}
		n   int
	}{
		{1, 1},
	}
	pinMap := PinMap{
		&PinDesc{ID: "P1_1", Aliases: []string{"1"}, Caps: CapAnalog, AnalogLogical: 1},
	}
	driver := NewGPIODriver(pinMap, nil, newFakeAnalogPin, nil)
	for _, test := range tests {
		pin, err := driver.AnalogPin(test.key)
		if err != nil {
			t.Errorf("Looking up %v: unexpected error: %v", test.key, err)
			continue
		}
		if pin.N() != test.n {
			t.Errorf("Looking up %v: got %v, want %v", test.key, pin.N(), test.n)
		}
	}
}

func TestGpioDriverUnavailablePinType(t *testing.T) {
	pinMap := PinMap{
		&PinDesc{ID: "P1_1", Aliases: []string{"1"}, Caps: CapDigital, DigitalLogical: 1},
		&PinDesc{ID: "P1_2", Aliases: []string{"1"}, Caps: CapAnalog, AnalogLogical: 1},
	}
	driver := NewGPIODriver(pinMap, nil, nil, nil)
	_, err := driver.DigitalPin(1)
	if err == nil {
		t.Fatal("Looking up digital pin 1: did not get error")
	}
	expected := "gpio: digital io not supported on this host"
	if err.Error() != expected {
		t.Fatalf("Looking up digital pin 1: got error %q, expected %q", err, expected)
	}
	_, err = driver.AnalogPin(1)
	if err == nil {
		t.Fatal("Looking up analog pin 1: did not get error")
	}
	expected = "gpio: analog io not supported on this host"
	if err.Error() != expected {
		t.Fatalf("Looking up analog pin 1: got error %q, expected %q", err, expected)
	}
}

func TestGpioPinCaching(t *testing.T) {
	pinMap := PinMap{
		&PinDesc{ID: "P1_1", Aliases: []string{"1"}, Caps: CapDigital},
	}
	driver := NewGPIODriver(pinMap, newFakeDigitalPin, nil, nil)
	pin, err := driver.DigitalPin(1)
	if err != nil {
		t.Fatalf("Looking up digital pin 1: got %v", err)
	}
	// Lookup the same pin again
	pin2, err := driver.DigitalPin(1)
	if err != nil {
		t.Fatalf("Looking up digital pin 1: got %v", err)
	}
	if pin != pin2 {
		t.Fatalf("Looking up digital pin 1 for the second time: got %v, want %v", &pin2, &pin)
	}
	// Looking up a closed pin
	pin.Close()
	pin3, err := driver.DigitalPin(1)
	if err != nil {
		t.Fatalf("Looking up digital pin 1: got %v", err)
		return
	}
	if pin == pin3 {
		t.Fatal("Looking up a closed pin, but got the same old instance")
	}
}
