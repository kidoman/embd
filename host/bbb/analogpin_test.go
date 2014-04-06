package bbb

import (
	"testing"

	"github.com/kidoman/embd"
)

func TestAnalogPinClose(t *testing.T) {
	pinMap := embd.PinMap{
		&embd.PinDesc{ID: "P1_1", Aliases: []string{"1"}, Caps: embd.CapAnalog},
	}
	driver := embd.NewGPIODriver(pinMap, nil, newAnalogPin, nil)
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
