/*
Package characterdisplay provides an ease-of-use layer on top of a character
display controller.
*/
package characterdisplay

// Controller is an interface that describes the basic functionality of a character
// display controller.
type Controller interface {
	DisplayOff() error            // turns the display off
	DisplayOn() error             // turns the display on
	CursorOff() error             // sets the cursor visibility to off
	CursorOn() error              // sets the cursor visibility to on
	BlinkOff() error              // sets the cursor blink off
	BlinkOn() error               // sets the cursor blink on
	ShiftLeft() error             // moves the cursor and text one column to the left
	ShiftRight() error            // moves the cursor and text one column to the right
	BacklightOff() error          // turns the display backlight off
	BacklightOn() error           // turns the display backlight on
	Home() error                  // moves the cursor to the home position
	Clear() error                 // clears the display and moves the cursor to the home position
	WriteChar(byte) error         // writes a character to the display
	SetCursor(col, row int) error // sets the cursor position
	Close() error                 // closes the controller resources
}

// Display represents an abstract character display and provides a
// ease-of-use layer on top of a character display controller.
type Display struct {
	Controller
	cols, rows int
	p          *position
}

type position struct {
	col int
	row int
}

// New creates a new Display
func New(controller Controller, cols, rows int) *Display {
	return &Display{
		Controller: controller,
		cols:       cols,
		rows:       rows,
		p:          &position{0, 0},
	}
}

// Home moves the cursor and all characters to the home position.
func (disp *Display) Home() error {
	disp.setCurrentPosition(0, 0)
	return disp.Controller.Home()
}

// Clear clears the display, preserving the mode settings and setting the correct home.
func (disp *Display) Clear() error {
	disp.setCurrentPosition(0, 0)
	return disp.Controller.Clear()
}

// Message prints the given string on the display, including interpreting newline
// characters and wrapping at the end of lines.
func (disp *Display) Message(message string) error {
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
		disp.p.col++
		if disp.p.col >= disp.cols || disp.p.col < 0 {
			err := disp.Newline()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Newline moves the input cursor to the beginning of the next line.
func (disp *Display) Newline() error {
	return disp.SetCursor(0, disp.p.row+1)
}

// SetCursor sets the input cursor to the given position.
func (disp *Display) SetCursor(col, row int) error {
	if row >= disp.rows {
		row = disp.rows - 1
	}
	disp.setCurrentPosition(col, row)
	return disp.Controller.SetCursor(col, row)
}

func (disp *Display) setCurrentPosition(col, row int) {
	disp.p.col = col
	disp.p.row = row
}
