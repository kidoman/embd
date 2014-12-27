package hd44780

import (
	"testing"

	"github.com/kidoman/embd"
)

const (
	rows = 20
	cols = 4
)

func TestNewGPIOCharacterDisplay_initPins(t *testing.T) {
	var pins []*mockDigitalPin
	for i := 0; i < 7; i++ {
		pins = append(pins, newMockDigitalPin())
	}
	NewGPIOCharacterDisplay(
		pins[0],
		pins[1],
		pins[2],
		pins[3],
		pins[4],
		pins[5],
		pins[6],
		Negative,
		cols,
		rows,
	)
	for idx, pin := range pins {
		if pin.direction != embd.Out {
			t.Errorf("Pin %d not set to direction Out(%d), set to %d", idx, embd.Out, pin.direction)
		}
	}
}

func TestDefaultModes(t *testing.T) {
	displayGPIO, _ := NewGPIOCharacterDisplay(
		newMockDigitalPin(),
		newMockDigitalPin(),
		newMockDigitalPin(),
		newMockDigitalPin(),
		newMockDigitalPin(),
		newMockDigitalPin(),
		newMockDigitalPin(),
		Negative,
		cols,
		rows,
	)
	displayI2C, _ := NewI2CCharacterDisplay(
		newMockI2CBus(),
		testAddr,
		MJKDZPinMap,
		cols,
		rows,
	)

	for idx, display := range []*CharacterDisplay{displayGPIO, displayI2C} {
		if display.EightBitModeEnabled() {
			t.Errorf("Display %d: Expected display to be initialized in 4-bit mode", idx)
		}
		if display.TwoLineEnabled() {
			t.Errorf("Display %d: Expected display to be initialized in one-line mode", idx)
		}
		if display.Dots5x10Enabled() {
			t.Errorf("Display %d: Expected display to be initialized in 5x8-dots mode", idx)
		}
		if !display.EntryIncrementEnabled() {
			t.Errorf("Display %d: Expected display to be initialized in entry increment mode", idx)
		}
		if display.EntryShiftEnabled() {
			t.Errorf("Display %d: Expected display to be initialized in entry shift off mode", idx)
		}
		if !display.DisplayEnabled() {
			t.Errorf("Display %d: Expected display to be initialized in display on mode", idx)
		}
		if display.CursorEnabled() {
			t.Errorf("Display %d: Expected display to be initialized in cursor off mode", idx)
		}
		if display.BlinkEnabled() {
			t.Errorf("Display %d: Expected display to be initialized in blink off mode", idx)
		}
	}
}
