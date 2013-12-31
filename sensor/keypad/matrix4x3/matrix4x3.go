// Package matrix4x3 allows interfacing 4x3 keypad with Raspberry pi.
package matrix4x3

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/stianeikeland/go-rpio"
)

type Key int

func (k Key) String() string {
	switch k {
	case KStar:
		return "*"
	case KHash:
		return "#"
	default:
		return strconv.Itoa(int(k))
	}
}

const (
	K0 Key = iota
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

// A Matrix4x3 interface implements access to the keypad.
type Matrix4x3 interface {
	// Run starts the continuous key scan loop.
	Run()

	// SetPollDelay sets the delay between runs of key scan acquisition loop.
	SetPollDelay(delay int)

	// Pressed key returns the current key pressed on the keypad.
	PressedKey() (Key, error)

	// Close.
	Close() error
}

type matrix4x3 struct {
	rpioRowPins, rpioColPins []rpio.Pin

	initialized bool
	mu          sync.RWMutex

	poll int

	keyPressed chan Key
	quit       chan bool
}

// New creates a new interface for matrix4x3.
func New(rowPins, colPins []int) Matrix4x3 {
	m := &matrix4x3{
		rpioRowPins: make([]rpio.Pin, rows),
		rpioColPins: make([]rpio.Pin, cols),
		poll:        pollDelay,
	}

	for i := 0; i < rows; i++ {
		m.rpioRowPins[i] = rpio.Pin(rowPins[i])
	}
	for i := 0; i < cols; i++ {
		m.rpioColPins[i] = rpio.Pin(colPins[i])
	}

	return m
}

// SetPollDelay sets the delay between run of key scan acquisition loop.
func (d *matrix4x3) SetPollDelay(delay int) {
	d.poll = delay
}

func (d *matrix4x3) setup() (err error) {
	d.mu.RLock()
	if d.initialized {
		d.mu.RUnlock()
		return
	}
	d.mu.RUnlock()

	d.mu.Lock()
	defer d.mu.Unlock()

	if err = rpio.Open(); err != nil {
		return
	}

	for i := 0; i < rows; i++ {
		d.rpioRowPins[i].Input()
		d.rpioRowPins[i].PullUp()
	}

	for i := 0; i < cols; i++ {
		d.rpioColPins[i].Output()
		d.rpioColPins[i].High()
	}

	d.initialized = true

	return
}

func (d *matrix4x3) findPressedKey() (key Key, err error) {
	if err = d.setup(); err != nil {
		return
	}

	err = errors.New("no key pressed")

	for col := 0; col < cols; col++ {
		d.rpioColPins[col].Low()
		for row := 0; row < rows; row++ {
			if d.rpioRowPins[row].Read() == rpio.Low {
				time.Sleep(debounce)

				if d.rpioRowPins[row].Read() == rpio.Low {
					key = keyMap[row][col]
					err = nil
				}
			}
		}
		d.rpioColPins[col].High()
	}

	return
}

// Pressed key returns the current key pressed on the keypad.
func (d *matrix4x3) PressedKey() (key Key, err error) {
	select {
	case key = <-d.keyPressed:
		return
	default:
		return d.findPressedKey()
	}
}

// Run starts the continuous key scan loop.
func (d *matrix4x3) Run() {
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
func (d *matrix4x3) Close() (err error) {
	if d.quit != nil {
		d.quit <- true
	}

	rpio.Close()

	return
}
