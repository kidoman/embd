package embd

import "testing"

type fakeDigitalPin struct {
	n int
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

func (*fakeDigitalPin) ActiveLow(b bool) error {
	return nil
}

func (*fakeDigitalPin) PullUp() error {
	return nil
}

func (*fakeDigitalPin) PullDown() error {
	return nil
}

func (*fakeDigitalPin) Close() error {
	return nil
}

func newFakeDigitalPin(n int) DigitalPin {
	return &fakeDigitalPin{n}
}

func TestGpioDriverDigitalPin(t *testing.T) {
	var tests = []struct {
		key interface{}
		n   int
	}{
		{1, 1},
	}
	var pinMap = PinMap{
		&PinDesc{ID: "P1_1", Aliases: []string{"1"}, Caps: CapDigital, DigitalLogical: 1},
	}
	driver := newGPIODriver(pinMap, newFakeDigitalPin, nil)
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
	n int
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

func newFakeAnalogPin(n int) AnalogPin {
	return &fakeAnalogPin{n}
}

func TestGpioDriverAnalogPin(t *testing.T) {
	var tests = []struct {
		key interface{}
		n   int
	}{
		{1, 1},
	}
	var pinMap = PinMap{
		&PinDesc{ID: "P1_1", Aliases: []string{"1"}, Caps: CapAnalog, AnalogLogical: 1},
	}
	driver := newGPIODriver(pinMap, nil, newFakeAnalogPin)
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
	var pinMap = PinMap{
		&PinDesc{ID: "P1_1", Aliases: []string{"1"}, Caps: CapDigital, DigitalLogical: 1},
		&PinDesc{ID: "P1_2", Aliases: []string{"1"}, Caps: CapAnalog, AnalogLogical: 1},
	}
	driver := newGPIODriver(pinMap, nil, nil)
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
