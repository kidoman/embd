package hd44780

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/kidoman/embd"
)

const testAddr byte = 0x20

type mockDigitalPin struct {
	direction embd.Direction
	values    chan int
	closed    bool
}

func newMockDigitalPin() *mockDigitalPin {
	return &mockDigitalPin{
		values: make(chan int, 256),
		closed: false,
	}
}

func (pin *mockDigitalPin) Watch(edge embd.Edge, handler func(embd.DigitalPin)) error { return nil }
func (pin *mockDigitalPin) StopWatching() error                                       { return nil }
func (pin *mockDigitalPin) N() int                                                    { return 0 }
func (pin *mockDigitalPin) Read() (int, error)                                        { return 0, nil }
func (pin *mockDigitalPin) TimePulse(state int) (time.Duration, error)                { return time.Duration(0), nil }
func (pin *mockDigitalPin) ActiveLow(b bool) error                                    { return nil }
func (pin *mockDigitalPin) PullUp() error                                             { return nil }
func (pin *mockDigitalPin) PullDown() error                                           { return nil }

func (pin *mockDigitalPin) Write(val int) error {
	pin.values <- val
	return nil
}

func (pin *mockDigitalPin) SetDirection(dir embd.Direction) error {
	pin.direction = dir
	return nil
}

func (pin *mockDigitalPin) Close() error {
	pin.closed = true
	return nil
}

type mockGPIOConnection struct {
	rs, en         *mockDigitalPin
	d4, d5, d6, d7 *mockDigitalPin
	backlight      *mockDigitalPin
	writes         []instruction
}

type instruction struct {
	rs   int
	data byte
}

func (ins *instruction) printAsBinary() string {
	return fmt.Sprintf("RS:%d|Byte:%s", ins.rs, printByteAsBinary(ins.data))
}

func printInstructionsAsBinary(ins []instruction) string {
	var binary []string
	for _, i := range ins {
		binary = append(binary, i.printAsBinary())
	}
	return fmt.Sprintf("%+v", binary)
}

func newMockGPIOConnection() *mockGPIOConnection {
	be := &mockGPIOConnection{
		rs:        newMockDigitalPin(),
		en:        newMockDigitalPin(),
		d4:        newMockDigitalPin(),
		d5:        newMockDigitalPin(),
		d6:        newMockDigitalPin(),
		d7:        newMockDigitalPin(),
		backlight: newMockDigitalPin(),
	}
	go func() {
		for {
			var b byte = 0x00
			var rs int = 0
			// wait for EN low,high,low then read high nibble
			if <-be.en.values == embd.Low &&
				<-be.en.values == embd.High &&
				<-be.en.values == embd.Low {
				rs = <-be.rs.values
				b |= byte(<-be.d4.values) << 4
				b |= byte(<-be.d5.values) << 5
				b |= byte(<-be.d6.values) << 6
				b |= byte(<-be.d7.values) << 7
			}
			// wait for EN low,high,low then read low nibble
			if <-be.en.values == embd.Low &&
				<-be.en.values == embd.High &&
				<-be.en.values == embd.Low {
				b |= byte(<-be.d4.values)
				b |= byte(<-be.d5.values) << 1
				b |= byte(<-be.d6.values) << 2
				b |= byte(<-be.d7.values) << 3
				be.writes = append(be.writes, instruction{rs, b})
			}
		}
	}()
	return be
}

func (be *mockGPIOConnection) pins() []*mockDigitalPin {
	return []*mockDigitalPin{be.rs, be.en, be.d4, be.d5, be.d6, be.d7, be.backlight}
}

type mockI2CBus struct {
	writes []byte
	closed bool
}

func (bus *mockI2CBus) ReadByte(addr byte) (byte, error)                  { return 0x00, nil }
func (bus *mockI2CBus) WriteBytes(addr byte, value []byte) error          { return nil }
func (bus *mockI2CBus) ReadFromReg(addr, reg byte, value []byte) error    { return nil }
func (bus *mockI2CBus) ReadByteFromReg(addr, reg byte) (byte, error)      { return 0x00, nil }
func (bus *mockI2CBus) ReadWordFromReg(addr, reg byte) (uint16, error)    { return 0, nil }
func (bus *mockI2CBus) WriteToReg(addr, reg byte, value []byte) error     { return nil }
func (bus *mockI2CBus) WriteByteToReg(addr, reg, value byte) error        { return nil }
func (bus *mockI2CBus) WriteWordToReg(addr, reg byte, value uint16) error { return nil }

func (bus *mockI2CBus) WriteByte(addr, value byte) error {
	bus.writes = append(bus.writes, value)
	return nil
}

func (bus *mockI2CBus) Close() error {
	bus.closed = true
	return nil
}

func newMockI2CBus() *mockI2CBus {
	return &mockI2CBus{closed: false}
}

func printByteAsBinary(b byte) string {
	return fmt.Sprintf("%08b(%#x)", b, b)
}

func printBytesAsBinary(bytes []byte) string {
	var binary []string
	for _, w := range bytes {
		binary = append(binary, printByteAsBinary(w))
	}
	return fmt.Sprintf("%+v", binary)
}

func TestInitialize4Bit_directionOut(t *testing.T) {
	be := newMockGPIOConnection()
	NewGPIO(be.rs, be.en, be.d4, be.d5, be.d6, be.d7, be.backlight, Negative)
	for idx, pin := range be.pins() {
		if pin.direction != embd.Out {
			t.Errorf("Pin %d not set to direction Out", idx)
		}
	}
}

func TestInitialize4Bit_lcdInit(t *testing.T) {
	be := newMockGPIOConnection()
	NewGPIO(be.rs, be.en, be.d4, be.d5, be.d6, be.d7, be.backlight, Negative)
	instructions := []instruction{
		instruction{embd.Low, lcdInit},
		instruction{embd.Low, lcdInit4bit},
		instruction{embd.Low, byte(lcdSetEntryMode)},
		instruction{embd.Low, byte(lcdSetDisplayMode)},
		instruction{embd.Low, byte(lcdSetFunctionMode)},
	}

	if !reflect.DeepEqual(instructions, be.writes) {
		t.Errorf(
			"\nExpected\t%s\nActual\t\t%+v",
			printInstructionsAsBinary(instructions),
			printInstructionsAsBinary(be.writes))
	}
}

func TestGPIOConnectionClose(t *testing.T) {
	be := newMockGPIOConnection()
	bus, _ := NewGPIO(be.rs, be.en, be.d4, be.d5, be.d6, be.d7, be.backlight, Negative)
	bus.Close()
	for idx, pin := range be.pins() {
		if !pin.closed {
			t.Errorf("Pin %d was not closed", idx)
		}
	}
}

func TestI2CConnectionPinMap(t *testing.T) {
	cases := []map[string]interface{}{
		map[string]interface{}{
			"instruction": lcdDisplayMove | lcdMoveRight,
			"pinMap":      MJKDZPinMap,
			"expected": []byte{
				0x0,  // 00000000 high nibble
				0x10, // 00010000
				0x0,  // 00000000
				0xc,  // 00001100 low nibble
				0x1c, // 00011100
				0xc,  // 00001100
			},
		},
		map[string]interface{}{
			"instruction": lcdDisplayMove | lcdMoveRight,
			"pinMap":      PCF8574PinMap,
			"expected": []byte{
				0x8,  // 00001000 high nibble
				0xc,  // 00001100
				0x8,  // 00001000
				0xc8, // 11001000 low nibble
				0xcc, // 11001100
				0xc8, // 11001000
			},
		},
	}

	for idx, c := range cases {
		instruction := c["instruction"].(byte)
		pinMap := c["pinMap"].(I2CPinMap)
		expected := c["expected"].([]byte)

		i2c := newMockI2CBus()
		conn := NewI2CConnection(i2c, testAddr, pinMap)
		rawInstruction := instruction
		// instructions (RS = false) with backlight on
		conn.Backlight = true
		conn.Write(false, rawInstruction)

		if !reflect.DeepEqual(expected, i2c.writes) {
			t.Errorf(
				"Case %d:\nExpected\t%s\nActual\t\t%s",
				idx+1,
				printBytesAsBinary(expected),
				printBytesAsBinary(i2c.writes))
		}
	}
}

func TestI2CConnectionClose(t *testing.T) {
	i2c := newMockI2CBus()
	conn := NewI2CConnection(i2c, testAddr, MJKDZPinMap)
	conn.Close()
	if !i2c.closed {
		t.Error("I2C bus was not closed")
	}
}
