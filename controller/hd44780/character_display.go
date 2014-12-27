package hd44780

import (
	"github.com/golang/glog"
	"github.com/kidoman/embd"
)

// DefaultModes are the default initialization modes for a CharacterDisplay.
var DefaultModes []ModeSetter = []ModeSetter{
	FourBitMode,
	OneLine,
	Dots5x8,
	EntryIncrement,
	EntryShiftOff,
	DisplayOn,
	CursorOff,
	BlinkOff,
}

// CharacterDisplay represents an abstract character display and provides a
// convenience layer on top of the basic HD44780 library.
type CharacterDisplay struct {
	*HD44780
	Cols int
	Rows int
	p    *position
}

type position struct {
	col int
	row int
}

// NewGPIOCharacterDisplay creates a new CharacterDisplay connected by a 4-bit GPIO bus.
func NewGPIOCharacterDisplay(
	rs, en, d4, d5, d6, d7, backlight interface{},
	blPolarity Polarity,
	cols, rows int,
	modes ...ModeSetter,
) (*CharacterDisplay, error) {
	pinKeys := []interface{}{rs, en, d4, d5, d6, d7, backlight}
	pins := [7]embd.DigitalPin{}
	for idx, key := range pinKeys {
		if key == nil {
			continue
		}
		var digitalPin embd.DigitalPin
		if pin, ok := key.(embd.DigitalPin); ok {
			digitalPin = pin
		} else {
			var err error
			digitalPin, err = embd.NewDigitalPin(key)
			if err != nil {
				glog.V(1).Infof("hd44780: error creating digital pin %+v: %s", key, err)
				return nil, err
			}
		}
		pins[idx] = digitalPin
	}
	hd, err := NewGPIO(
		pins[0],
		pins[1],
		pins[2],
		pins[3],
		pins[4],
		pins[5],
		pins[6],
		blPolarity,
		append(DefaultModes, modes...)...,
	)
	if err != nil {
		return nil, err
	}
	return NewCharacterDisplay(hd, cols, rows)
}

// NewI2CCharacterDisplay creates a new CharacterDisplay connected by an I²C bus.
func NewI2CCharacterDisplay(
	i2c embd.I2CBus,
	addr byte,
	pinMap I2CPinMap,
	cols, rows int,
	modes ...ModeSetter,
) (*CharacterDisplay, error) {
	hd, err := NewI2C(i2c, addr, pinMap, append(DefaultModes, modes...)...)
	if err != nil {
		return nil, err
	}
	return NewCharacterDisplay(hd, cols, rows)
}

// NewCharacterDisplay creates a new character display abstraction for an
// HD44780-compatible controller.
func NewCharacterDisplay(hd *HD44780, cols, rows int) (*CharacterDisplay, error) {
	display := &CharacterDisplay{
		HD44780: hd,
		Cols:    cols,
		Rows:    rows,
		p:       &position{0, 0},
	}
	err := display.BacklightOn()
	if err != nil {
		return nil, err
	}
	return display, nil
}

// Home moves the cursor and all characters to the home position.
func (disp *CharacterDisplay) Home() error {
	disp.currentPosition(0, 0)
	return disp.HD44780.Home()
}

// Clear clears the display, preserving the mode settings and setting the correct home.
func (disp *CharacterDisplay) Clear() error {
	disp.currentPosition(0, 0)
	err := disp.HD44780.Clear()
	if err != nil {
		return err
	}
	err = disp.SetMode()
	if err != nil {
		return err
	}
	if !disp.isLeftToRight() {
		return disp.SetCursor(disp.Cols-1, 0)
	}
	return nil
}

// Message prints the given string on the display.
func (disp *CharacterDisplay) Message(message string) error {
	bytes := []byte(message)
	for _, b := range bytes {
		if b == byte('\n') {
			err := disp.Newline()
			if err != nil {
				return err
			}
			continue
		}
		err := disp.WriteChar(b)
		if err != nil {
			return err
		}
		if disp.isLeftToRight() {
			disp.p.col++
		} else {
			disp.p.col--
		}
		if disp.p.col >= disp.Cols || disp.p.col < 0 {
			err := disp.Newline()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Newline moves the input cursor to the beginning of the next line.
func (disp *CharacterDisplay) Newline() error {
	var col int
	if disp.isLeftToRight() {
		col = 0
	} else {
		col = disp.Cols - 1
	}
	return disp.SetCursor(col, disp.p.row+1)
}

func (disp *CharacterDisplay) isLeftToRight() bool {
	// EntryIncrement and EntryShiftOn is right-to-left
	// EntryDecrement and EntryShiftOn is left-to-right
	// EntryIncrement and EntryShiftOff is left-to-right
	// EntryDecrement and EntryShiftOff is right-to-left
	return disp.EntryIncrementEnabled() != disp.EntryShiftEnabled()
}

// SetCursor sets the input cursor to the given position.
func (disp *CharacterDisplay) SetCursor(col, row int) error {
	if row >= disp.Rows {
		row = disp.Rows - 1
	}
	disp.currentPosition(col, row)
	return disp.HD44780.SetCursor(byte(col) + disp.lcdRowOffset(row))
}

func (disp *CharacterDisplay) lcdRowOffset(row int) byte {
	// Offset for up to 4 rows
	if row > 3 {
		row = 3
	}
	switch disp.Cols {
	case 16:
		// 16-char line mappings
		return []byte{0x00, 0x40, 0x10, 0x50}[row]
	default:
		// default to the 20-char line mappings
		return []byte{0x00, 0x40, 0x14, 0x54}[row]
	}
}

func (disp *CharacterDisplay) currentPosition(col, row int) {
	disp.p.col = col
	disp.p.row = row
}

// Close closes the underlying HD44780 controller.
func (disp *CharacterDisplay) Close() error {
	return disp.HD44780.Close()
}
