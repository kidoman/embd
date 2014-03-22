package embd

import "testing"

func TestGpioDriverDigitalPin(t *testing.T) {
	var tests = []struct {
		key interface{}
		n   int
	}{
		{1, 1},
	}
	var pinMap = PinMap{
		&PinDesc{ID: "P1_1", Aliases: []string{"1"}, Caps: CapNormal, DigitalLogical: 1},
	}
	driver := newGPIODriver(pinMap)
	for _, test := range tests {
		pin, err := driver.digitalPin(test.key)
		if err != nil {
			t.Errorf("Looking up %v: unexpected error: %v", test.key, err)
			continue
		}
		if pin.n != test.n {
			t.Errorf("Looking up %v: got %v, want %v", test.key, pin.n, test.n)
		}
	}
}
