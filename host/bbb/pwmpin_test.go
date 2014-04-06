package bbb

import (
	"testing"

	"github.com/kidoman/embd"
)

func TestPWMPinClose(t *testing.T) {
	pinMap := embd.PinMap{
		&embd.PinDesc{ID: "P1_1", Aliases: []string{"1"}, Caps: embd.CapPWM},
	}
	driver := embd.NewGPIODriver(pinMap, nil, nil, newPWMPin)
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
