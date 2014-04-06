package embd

import "testing"

func TestDigitalPinClose(t *testing.T) {
	pinMap := PinMap{
		&PinDesc{ID: "P1_1", Aliases: []string{"1"}, Caps: CapDigital},
	}
	driver := newGPIODriver(pinMap, newDigitalPin, nil, nil)
	pin, err := driver.DigitalPin(1)
	if err != nil {
		t.Fatalf("Looking up digital pin 1: got %v", err)
	}
	pin.Close()
	pin2, err := driver.DigitalPin(1)
	if err != nil {
		t.Fatalf("Looking up digital pin 1: got %v", err)
	}
	if pin == pin2 {
		t.Fatal("Looking up closed digital pin 1: but got the old instance")
	}
}
