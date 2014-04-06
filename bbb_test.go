package embd

import "testing"

func TestBBBAnalogPinClose(t *testing.T) {
	pinMap := PinMap{
		&PinDesc{ID: "P1_1", Aliases: []string{"1"}, Caps: CapAnalog},
	}
	driver := newGPIODriver(pinMap, nil, newBBBAnalogPin, nil)
	pin, err := driver.AnalogPin(1)
	if err != nil {
		t.Fatalf("Looking up analog pin 1: got %v", err)
	}
	pin.Close()
	pin2, err := driver.AnalogPin(1)
	if err != nil {
		t.Fatalf("Looking up analog pin 1: got %v", err)
	}
	if pin == pin2 {
		t.Fatal("Looking up closed analog pin 1: but got the old instance")
	}
}

func TestBBBPWMPinClose(t *testing.T) {
	pinMap := PinMap{
		&PinDesc{ID: "P1_1", Aliases: []string{"1"}, Caps: CapPWM},
	}
	driver := newGPIODriver(pinMap, nil, nil, newBBBPWMPin)
	pin, err := driver.PWMPin(1)
	if err != nil {
		t.Fatalf("Looking up pwm pin 1: got %v", err)
	}
	pin.Close()
	pin2, err := driver.PWMPin(1)
	if err != nil {
		t.Fatalf("Looking up pwm pin 1: got %v", err)
	}
	if pin == pin2 {
		t.Fatal("Looking up closed pwm pin 1: but got the old instance")
	}
}
