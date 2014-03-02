// Package matrix4x3 allows interfacing 4x3 keypad with Raspberry pi.
package matrix4x3

import (
	"strconv"
	"sync"
	"time"

	"github.com/kidoman/embd"
)

type Key int

func (k Key) String() string {
	switch k {
	case KStar:
		return "*"
	case KHash:
		return "#"
	default:
		return strconv.Itoa(int(k) - 1)
	}
}

const (
	KNone Key = iota
	K0
	K1
	K2
	K3
	K4
	K5
	K6
	K7
	K8
	K9
	KStar
	KHash

	debounce = 20 * time.Millisecond

	pollDelay = 150

	rows = 4
	cols = 3
)

var keyMap [][]Key

func init() {
	keyMap = make([][]Key, rows)
	for i := 0; i < rows; i++ {
		keyMap[i] = make([]Key, cols)
	}
	keyMap[0][0] = K1
	keyMap[0][1] = K2
	keyMap[0][2] = K3
	keyMap[1][0] = K4
	keyMap[1][1] = K5
	keyMap[1][2] = K6
	keyMap[2][0] = K7
	keyMap[2][1] = K8
	keyMap[2][2] = K9
	keyMap[3][0] = KStar
	keyMap[3][1] = K0
	keyMap[3][2] = KHash
}

// A Matrix4x3 struct represents access to the keypad.
type Matrix4x3 struct {
	rowPins, colPins []embd.DigitalPin

	initialized bool
	mu          sync.RWMutex

	poll int

	keyPressed chan Key
	quit       chan bool
}

// New creates a new interface for matrix4x3.
func New(rowPins, colPins []int) (*Matrix4x3, error) {
	m := &Matrix4x3{
		rowPins: make([]embd.DigitalPin, rows),
		colPins: make([]embd.DigitalPin, cols),
		poll:    pollDelay,
	}

	var err error
	for i := 0; i < rows; i++ {
		m.rowPins[i], err = embd.NewDigitalPin(rowPins[i])
		if err != nil {
			return nil, err
		}
	}
	for i := 0; i < cols; i++ {
		m.colPins[i], err = embd.NewDigitalPin(colPins[i])
		if err != nil {
			return nil, err
		}
	}

	return m, nil
}

// SetPollDelay sets the delay between run of key scan acquisition loop.
func (d *Matrix4x3) SetPollDelay(delay int) {
	d.poll = delay
}

func (d *Matrix4x3) setup() error {
	d.mu.RLock()
	if d.initialized {
		d.mu.RUnlock()
		return nil
	}
	d.mu.RUnlock()

	d.mu.Lock()
	defer d.mu.Unlock()

	for i := 0; i < rows; i++ {
		if err := d.rowPins[i].SetDirection(embd.In); err != nil {
			return err
		}
		if err := d.rowPins[i].PullUp(); err != nil {
			return err
		}
	}

	for i := 0; i < cols; i++ {
		if err := d.colPins[i].SetDirection(embd.Out); err != nil {
			return err
		}
		if err := d.colPins[i].Write(embd.High); err != nil {
			return err
		}
	}

	d.initialized = true

	return nil
}

func (d *Matrix4x3) findPressedKey() (Key, error) {
	if err := d.setup(); err != nil {
		return 0, err
	}

	for col := 0; col < cols; col++ {
		if err := d.colPins[col].Write(embd.Low); err != nil {
			return KNone, err
		}
		for row := 0; row < rows; row++ {
			value, err := d.rowPins[row].Read()
			if err != nil {
				return KNone, err
			}
			if value == embd.Low {
				time.Sleep(debounce)

				value, err = d.rowPins[row].Read()
				if err != nil {
					return KNone, err
				}
				if value == embd.Low {
					if err := d.colPins[col].Write(embd.High); err != nil {
						return KNone, err
					}
					return keyMap[row][col], nil
				}
			}
		}
		if err := d.colPins[col].Write(embd.High); err != nil {
			return KNone, err
		}
	}

	return KNone, nil
}

// Pressed key returns the current key pressed on the keypad.
func (d *Matrix4x3) PressedKey() (key Key, err error) {
	select {
	case key = <-d.keyPressed:
		return
	default:
		return d.findPressedKey()
	}
}

// Run starts the continuous key scan loop.
func (d *Matrix4x3) Run() {
	d.quit = make(chan bool)

	go func() {
		timer := time.Tick(time.Duration(d.poll) * time.Millisecond)
		var key Key

		for {
			var keyUpdates chan Key

			select {
			case <-timer:
				var err error
				if key, err = d.findPressedKey(); err == nil {
					keyUpdates = d.keyPressed
				}
			case keyUpdates <- key:
				keyUpdates = nil
			case <-d.quit:
				d.keyPressed = nil

				return
			}
		}
	}()
}

// Close.
func (d *Matrix4x3) Close() {
	if d.quit != nil {
		d.quit <- true
	}
}
