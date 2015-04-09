/*
Package hd44780 allows controlling an HD44780-compatible character LCD
controller. Currently the library is write-only and does not support
reading from the display controller.

Resources

This library is based three other HD44780 libraries:
	Adafruit	https://github.com/adafruit/Adafruit-Raspberry-Pi-Python-Code/blob/master/Adafruit_CharLCD/Adafruit_CharLCD.py
	hwio		https://github.com/mrmorphic/hwio/blob/master/devices/hd44780/hd44780_i2c.go
	LiquidCrystal	https://github.com/arduino/Arduino/blob/master/libraries/LiquidCrystal/LiquidCrystal.cpp
*/
package hd44780

import (
	"time"

	"github.com/golang/glog"
	"github.com/kidoman/embd"
)

type entryMode byte
type displayMode byte
type functionMode byte

// RowAddress defines the cursor (DDRAM) address of the first column of each row, up to 4 rows.
// You must use the RowAddress value that matches the number of columns on your character display
// for the SetCursor function to work correctly.
type RowAddress [4]byte

var (
	// RowAddress16Col are row addresses for a 16-column display
	RowAddress16Col RowAddress = [4]byte{0x00, 0x40, 0x10, 0x50}
	// RowAddress20Col are row addresses for a 20-column display
	RowAddress20Col RowAddress = [4]byte{0x00, 0x40, 0x14, 0x54}
)

// BacklightPolarity is used to set the polarity of the backlight switch, either positive or negative.
type BacklightPolarity bool

const (
	// Negative indicates that the backlight is active-low and must have a logical low value to enable.
	Negative BacklightPolarity = false
	// Positive indicates that the backlight is active-high and must have a logical high value to enable.
	Positive BacklightPolarity = true

	writeDelay = 37 * time.Microsecond
	pulseDelay = 1 * time.Microsecond
	clearDelay = 1520 * time.Microsecond

	// Initialize display
	lcdInit     byte = 0x33 // 00110011
	lcdInit4bit byte = 0x32 // 00110010

	// Commands
	lcdClearDisplay byte = 0x01 // 00000001
	lcdReturnHome   byte = 0x02 // 00000010
	lcdCursorShift  byte = 0x10 // 00010000
	lcdSetCGRamAddr byte = 0x40 // 01000000
	lcdSetDDRamAddr byte = 0x80 // 10000000

	// Cursor and display move flags
	lcdCursorMove  byte = 0x00 // 00000000
	lcdDisplayMove byte = 0x08 // 00001000
	lcdMoveLeft    byte = 0x00 // 00000000
	lcdMoveRight   byte = 0x04 // 00000100

	// Entry mode flags
	lcdSetEntryMode   entryMode = 0x04 // 00000100
	lcdEntryDecrement entryMode = 0x00 // 00000000
	lcdEntryIncrement entryMode = 0x02 // 00000010
	lcdEntryShiftOff  entryMode = 0x00 // 00000000
	lcdEntryShiftOn   entryMode = 0x01 // 00000001

	// Display mode flags
	lcdSetDisplayMode displayMode = 0x08 // 00001000
	lcdDisplayOff     displayMode = 0x00 // 00000000
	lcdDisplayOn      displayMode = 0x04 // 00000100
	lcdCursorOff      displayMode = 0x00 // 00000000
	lcdCursorOn       displayMode = 0x02 // 00000010
	lcdBlinkOff       displayMode = 0x00 // 00000000
	lcdBlinkOn        displayMode = 0x01 // 00000001

	// Function mode flags
	lcdSetFunctionMode functionMode = 0x20 // 00100000
	lcd4BitMode        functionMode = 0x00 // 00000000
	lcd8BitMode        functionMode = 0x10 // 00010000
	lcd1Line           functionMode = 0x00 // 00000000
	lcd2Line           functionMode = 0x08 // 00001000
	lcd5x8Dots         functionMode = 0x00 // 00000000
	lcd5x10Dots        functionMode = 0x04 // 00000100
)

// HD44780 represents an HD44780-compatible character LCD controller.
type HD44780 struct {
	Connection
	eMode   entryMode
	dMode   displayMode
	fMode   functionMode
	rowAddr RowAddress
}

// NewGPIO creates a new HD44780 connected by a 4-bit GPIO bus.
func NewGPIO(
	rs, en, d4, d5, d6, d7, backlight interface{},
	blPolarity BacklightPolarity,
	rowAddr RowAddress,
	modes ...ModeSetter,
) (*HD44780, error) {
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
	for _, pin := range pins {
		if pin == nil {
			continue
		}
		err := pin.SetDirection(embd.Out)
		if err != nil {
			glog.Errorf("hd44780: error setting pin %+v to out direction: %s", pin, err)
			return nil, err
		}
	}
	return New(
		NewGPIOConnection(
			pins[0],
			pins[1],
			pins[2],
			pins[3],
			pins[4],
			pins[5],
			pins[6],
			blPolarity),
		rowAddr,
		modes...,
	)
}

// NewI2C creates a new HD44780 connected by an I²C bus.
func NewI2C(
	i2c embd.I2CBus,
	addr byte,
	pinMap I2CPinMap,
	rowAddr RowAddress,
	modes ...ModeSetter,
) (*HD44780, error) {
	return New(NewI2CConnection(i2c, addr, pinMap), rowAddr, modes...)
}

// New creates a new HD44780 connected by a Connection bus.
func New(bus Connection, rowAddr RowAddress, modes ...ModeSetter) (*HD44780, error) {
	controller := &HD44780{
		Connection: bus,
		eMode:      0x00,
		dMode:      0x00,
		fMode:      0x00,
		rowAddr:    rowAddr,
	}
	err := controller.lcdInit()
	if err != nil {
		return nil, err
	}
	err = controller.SetMode(append(DefaultModes, modes...)...)
	if err != nil {
		return nil, err
	}
	return controller, nil
}

func (controller *HD44780) lcdInit() error {
	glog.V(2).Info("hd44780: initializing display")
	err := controller.WriteInstruction(lcdInit)
	if err != nil {
		return err
	}
	glog.V(2).Info("hd44780: initializing display in 4-bit mode")
	return controller.WriteInstruction(lcdInit4bit)
}

// DefaultModes are the default initialization modes for an HD44780.
// ModeSetters passed in to a constructor will override these default values.
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

// ModeSetter defines a function used for setting modes on an HD44780.
// ModeSetters must be used with the SetMode function or in a constructor.
type ModeSetter func(*HD44780)

// EntryDecrement is a ModeSetter that sets the HD44780 to entry decrement mode.
func EntryDecrement(hd *HD44780) { hd.eMode &= ^lcdEntryIncrement }

// EntryIncrement is a ModeSetter that sets the HD44780 to entry increment mode.
func EntryIncrement(hd *HD44780) { hd.eMode |= lcdEntryIncrement }

// EntryShiftOff is a ModeSetter that sets the HD44780 to entry shift off mode.
func EntryShiftOff(hd *HD44780) { hd.eMode &= ^lcdEntryShiftOn }

// EntryShiftOn is a ModeSetter that sets the HD44780 to entry shift on mode.
func EntryShiftOn(hd *HD44780) { hd.eMode |= lcdEntryShiftOn }

// DisplayOff is a ModeSetter that sets the HD44780 to display off mode.
func DisplayOff(hd *HD44780) { hd.dMode &= ^lcdDisplayOn }

// DisplayOn is a ModeSetter that sets the HD44780 to display on mode.
func DisplayOn(hd *HD44780) { hd.dMode |= lcdDisplayOn }

// CursorOff is a ModeSetter that sets the HD44780 to cursor off mode.
func CursorOff(hd *HD44780) { hd.dMode &= ^lcdCursorOn }

// CursorOn is a ModeSetter that sets the HD44780 to cursor on mode.
func CursorOn(hd *HD44780) { hd.dMode |= lcdCursorOn }

// BlinkOff is a ModeSetter that sets the HD44780 to cursor blink off mode.
func BlinkOff(hd *HD44780) { hd.dMode &= ^lcdBlinkOn }

// BlinkOn is a ModeSetter that sets the HD44780 to cursor blink on mode.
func BlinkOn(hd *HD44780) { hd.dMode |= lcdBlinkOn }

// FourBitMode is a ModeSetter that sets the HD44780 to 4-bit bus mode.
func FourBitMode(hd *HD44780) { hd.fMode &= ^lcd8BitMode }

// EightBitMode is a ModeSetter that sets the HD44780 to 8-bit bus mode.
func EightBitMode(hd *HD44780) { hd.fMode |= lcd8BitMode }

// OneLine is a ModeSetter that sets the HD44780 to 1-line display mode.
func OneLine(hd *HD44780) { hd.fMode &= ^lcd2Line }

// TwoLine is a ModeSetter that sets the HD44780 to 2-line display mode.
func TwoLine(hd *HD44780) { hd.fMode |= lcd2Line }

// Dots5x8 is a ModeSetter that sets the HD44780 to 5x8-pixel character mode.
func Dots5x8(hd *HD44780) { hd.fMode &= ^lcd5x10Dots }

// Dots5x10 is a ModeSetter that sets the HD44780 to 5x10-pixel character mode.
func Dots5x10(hd *HD44780) { hd.fMode |= lcd5x10Dots }

// EntryIncrementEnabled returns true if entry increment mode is enabled.
func (hd *HD44780) EntryIncrementEnabled() bool { return hd.eMode&lcdEntryIncrement > 0 }

// EntryShiftEnabled returns true if entry shift mode is enabled.
func (hd *HD44780) EntryShiftEnabled() bool { return hd.eMode&lcdEntryShiftOn > 0 }

// DisplayEnabled returns true if the display is on.
func (hd *HD44780) DisplayEnabled() bool { return hd.dMode&lcdDisplayOn > 0 }

// CursorEnabled returns true if the cursor is on.
func (hd *HD44780) CursorEnabled() bool { return hd.dMode&lcdCursorOn > 0 }

// BlinkEnabled returns true if the cursor blink mode is enabled.
func (hd *HD44780) BlinkEnabled() bool { return hd.dMode&lcdBlinkOn > 0 }

// EightBitModeEnabled returns true if 8-bit bus mode is enabled and false if 4-bit
// bus mode is enabled.
func (hd *HD44780) EightBitModeEnabled() bool { return hd.fMode&lcd8BitMode > 0 }

// TwoLineEnabled returns true if 2-line display mode is enabled and false if 1-line
// display mode is enabled.
func (hd *HD44780) TwoLineEnabled() bool { return hd.fMode&lcd2Line > 0 }

// Dots5x10Enabled returns true if 5x10-pixel characters are enabled.
func (hd *HD44780) Dots5x10Enabled() bool { return hd.fMode&lcd5x8Dots > 0 }

// SetModes modifies the entry mode, display mode, and function mode with the
// given mode setter functions.
func (hd *HD44780) SetMode(modes ...ModeSetter) error {
	for _, m := range modes {
		m(hd)
	}
	functions := []func() error{
		func() error { return hd.setEntryMode() },
		func() error { return hd.setDisplayMode() },
		func() error { return hd.setFunctionMode() },
	}
	for _, f := range functions {
		err := f()
		if err != nil {
			return err
		}
	}
	return nil
}

func (hd *HD44780) setEntryMode() error {
	return hd.WriteInstruction(byte(lcdSetEntryMode | hd.eMode))
}

func (hd *HD44780) setDisplayMode() error {
	return hd.WriteInstruction(byte(lcdSetDisplayMode | hd.dMode))
}

func (hd *HD44780) setFunctionMode() error {
	return hd.WriteInstruction(byte(lcdSetFunctionMode | hd.fMode))
}

// DisplayOff sets the display mode to off.
func (hd *HD44780) DisplayOff() error {
	DisplayOff(hd)
	return hd.setDisplayMode()
}

// DisplayOn sets the display mode to on.
func (hd *HD44780) DisplayOn() error {
	DisplayOn(hd)
	return hd.setDisplayMode()
}

// CursorOff turns the cursor off.
func (hd *HD44780) CursorOff() error {
	CursorOff(hd)
	return hd.setDisplayMode()
}

// CursorOn turns the cursor on.
func (hd *HD44780) CursorOn() error {
	CursorOn(hd)
	return hd.setDisplayMode()
}

// BlinkOff sets cursor blink mode off.
func (hd *HD44780) BlinkOff() error {
	BlinkOff(hd)
	return hd.setDisplayMode()
}

// BlinkOn sets cursor blink mode on.
func (hd *HD44780) BlinkOn() error {
	BlinkOn(hd)
	return hd.setDisplayMode()
}

// ShiftLeft shifts the cursor and all characters to the left.
func (hd *HD44780) ShiftLeft() error {
	return hd.WriteInstruction(lcdCursorShift | lcdDisplayMove | lcdMoveLeft)
}

// ShiftRight shifts the cursor and all characters to the right.
func (hd *HD44780) ShiftRight() error {
	return hd.WriteInstruction(lcdCursorShift | lcdDisplayMove | lcdMoveRight)
}

// Home moves the cursor and all characters to the home position.
func (hd *HD44780) Home() error {
	err := hd.WriteInstruction(lcdReturnHome)
	time.Sleep(clearDelay)
	return err
}

// Clear clears the display and mode settings sets the cursor to the home position.
func (hd *HD44780) Clear() error {
	err := hd.WriteInstruction(lcdClearDisplay)
	if err != nil {
		return err
	}
	time.Sleep(clearDelay)
	// have to set mode here because clear also clears some mode settings
	return hd.SetMode()
}

// SetCursor sets the input cursor to the given position.
func (hd *HD44780) SetCursor(col, row int) error {
	return hd.SetDDRamAddr(byte(col) + hd.lcdRowOffset(row))
}

func (hd *HD44780) lcdRowOffset(row int) byte {
	// Offset for up to 4 rows
	if row > 3 {
		row = 3
	}
	return hd.rowAddr[row]
}

// SetDDRamAddr sets the input cursor to the given address.
func (hd *HD44780) SetDDRamAddr(value byte) error {
	return hd.WriteInstruction(lcdSetDDRamAddr | value)
}

// WriteInstruction writes a byte to the bus with register select in data mode.
func (hd *HD44780) WriteChar(value byte) error {
	return hd.Write(true, value)
}

// WriteInstruction writes a byte to the bus with register select in command mode.
func (hd *HD44780) WriteInstruction(value byte) error {
	return hd.Write(false, value)
}

// Close closes the underlying Connection.
func (hd *HD44780) Close() error {
	return hd.Connection.Close()
}

// Connection abstracts the different methods of communicating with an HD44780.
type Connection interface {
	// Write writes a byte to the HD44780 controller with the register select
	// flag either on or off.
	Write(rs bool, data byte) error

	// BacklightOff turns the optional backlight off.
	BacklightOff() error

	// BacklightOn turns the optional backlight on.
	BacklightOn() error

	// Close closes all open resources.
	Close() error
}

// GPIOConnection implements Connection using a 4-bit GPIO bus.
type GPIOConnection struct {
	RS, EN         embd.DigitalPin
	D4, D5, D6, D7 embd.DigitalPin
	Backlight      embd.DigitalPin
	BLPolarity     BacklightPolarity
}

// NewGPIOConnection returns a new Connection based on a 4-bit GPIO bus.
func NewGPIOConnection(
	rs, en, d4, d5, d6, d7, backlight embd.DigitalPin,
	blPolarity BacklightPolarity,
) *GPIOConnection {
	return &GPIOConnection{
		RS:         rs,
		EN:         en,
		D4:         d4,
		D5:         d5,
		D6:         d6,
		D7:         d7,
		Backlight:  backlight,
		BLPolarity: blPolarity,
	}
}

// BacklightOff turns the optional backlight off.
func (conn *GPIOConnection) BacklightOff() error {
	if conn.Backlight != nil {
		return conn.Backlight.Write(conn.backlightSignal(false))
	}
	return nil
}

// BacklightOn turns the optional backlight on.
func (conn *GPIOConnection) BacklightOn() error {
	if conn.Backlight != nil {
		return conn.Backlight.Write(conn.backlightSignal(true))
	}
	return nil
}

func (conn *GPIOConnection) backlightSignal(state bool) int {
	if state == bool(conn.BLPolarity) {
		return embd.High
	} else {
		return embd.Low
	}
}

// Write writes a register select flag and byte to the 4-bit GPIO connection.
func (conn *GPIOConnection) Write(rs bool, data byte) error {
	glog.V(3).Infof("hd44780: writing to GPIO RS: %t, data: %#x", rs, data)
	rsInt := embd.Low
	if rs {
		rsInt = embd.High
	}
	functions := []func() error{
		func() error { return conn.RS.Write(rsInt) },
		func() error { return conn.D4.Write(int((data >> 4) & 0x01)) },
		func() error { return conn.D5.Write(int((data >> 5) & 0x01)) },
		func() error { return conn.D6.Write(int((data >> 6) & 0x01)) },
		func() error { return conn.D7.Write(int((data >> 7) & 0x01)) },
		func() error { return conn.pulseEnable() },
		func() error { return conn.D4.Write(int(data & 0x01)) },
		func() error { return conn.D5.Write(int((data >> 1) & 0x01)) },
		func() error { return conn.D6.Write(int((data >> 2) & 0x01)) },
		func() error { return conn.D7.Write(int((data >> 3) & 0x01)) },
		func() error { return conn.pulseEnable() },
	}
	for _, f := range functions {
		err := f()
		if err != nil {
			return err
		}
	}
	time.Sleep(writeDelay)
	return nil
}

func (conn *GPIOConnection) pulseEnable() error {
	values := []int{embd.Low, embd.High, embd.Low}
	for _, v := range values {
		time.Sleep(pulseDelay)
		err := conn.EN.Write(v)
		if err != nil {
			return err
		}
	}
	return nil
}

// Close closes all open DigitalPins.
func (conn *GPIOConnection) Close() error {
	glog.V(2).Info("hd44780: closing all GPIO pins")
	pins := []embd.DigitalPin{
		conn.RS,
		conn.EN,
		conn.D4,
		conn.D5,
		conn.D6,
		conn.D7,
		conn.Backlight,
	}

	for _, pin := range pins {
		err := pin.Close()
		if err != nil {
			glog.Errorf("hd44780: error closing pin %+v: %s", pin, err)
			return err
		}
	}
	return nil
}

// I2CConnection implements Connection using an I²C bus.
type I2CConnection struct {
	I2C       embd.I2CBus
	Addr      byte
	PinMap    I2CPinMap
	Backlight bool
}

// I2CPinMap represents a mapping between the pins on an I²C port expander and
// the pins on the HD44780 controller.
type I2CPinMap struct {
	RS, RW, EN     byte
	D4, D5, D6, D7 byte
	Backlight      byte
	BLPolarity     BacklightPolarity
}

var (
	// MJKDZPinMap is the standard pin mapping for an MJKDZ-based I²C backpack.
	MJKDZPinMap I2CPinMap = I2CPinMap{
		RS: 6, RW: 5, EN: 4,
		D4: 0, D5: 1, D6: 2, D7: 3,
		Backlight:  7,
		BLPolarity: Negative,
	}
	// PCF8574PinMap is the standard pin mapping for a PCF8574-based I²C backpack.
	PCF8574PinMap I2CPinMap = I2CPinMap{
		RS: 0, RW: 1, EN: 2,
		D4: 4, D5: 5, D6: 6, D7: 7,
		Backlight:  3,
		BLPolarity: Positive,
	}
)

// NewI2CConnection returns a new Connection based on an I²C bus.
func NewI2CConnection(i2c embd.I2CBus, addr byte, pinMap I2CPinMap) *I2CConnection {
	return &I2CConnection{
		I2C:    i2c,
		Addr:   addr,
		PinMap: pinMap,
	}
}

// BacklightOff turns the optional backlight off.
func (conn *I2CConnection) BacklightOff() error {
	conn.Backlight = false
	return conn.Write(false, 0x00)
}

// BacklightOn turns the optional backlight on.
func (conn *I2CConnection) BacklightOn() error {
	conn.Backlight = true
	return conn.Write(false, 0x00)
}

// Write writes a register select flag and byte to the I²C connection.
func (conn *I2CConnection) Write(rs bool, data byte) error {
	var instructionHigh byte = 0x00
	instructionHigh |= ((data >> 4) & 0x01) << conn.PinMap.D4
	instructionHigh |= ((data >> 5) & 0x01) << conn.PinMap.D5
	instructionHigh |= ((data >> 6) & 0x01) << conn.PinMap.D6
	instructionHigh |= ((data >> 7) & 0x01) << conn.PinMap.D7

	var instructionLow byte = 0x00
	instructionLow |= (data & 0x01) << conn.PinMap.D4
	instructionLow |= ((data >> 1) & 0x01) << conn.PinMap.D5
	instructionLow |= ((data >> 2) & 0x01) << conn.PinMap.D6
	instructionLow |= ((data >> 3) & 0x01) << conn.PinMap.D7

	instructions := []byte{instructionHigh, instructionLow}
	for _, ins := range instructions {
		if rs {
			ins |= 0x01 << conn.PinMap.RS
		}
		if conn.Backlight == bool(conn.PinMap.BLPolarity) {
			ins |= 0x01 << conn.PinMap.Backlight
		}
		glog.V(3).Infof("hd44780: writing to I2C: %#x", ins)
		err := conn.pulseEnable(ins)
		if err != nil {
			return err
		}
	}
	time.Sleep(writeDelay)
	return nil
}

func (conn *I2CConnection) pulseEnable(data byte) error {
	bytes := []byte{data, data | (0x01 << conn.PinMap.EN), data}
	for _, b := range bytes {
		time.Sleep(pulseDelay)
		err := conn.I2C.WriteByte(conn.Addr, b)
		if err != nil {
			return err
		}
	}
	return nil
}

// Close closes the I²C connection.
func (conn *I2CConnection) Close() error {
	glog.V(2).Info("hd44780: closing I2C bus")
	return conn.I2C.Close()
}
