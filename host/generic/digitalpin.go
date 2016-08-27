// Digital IO support.
// This driver requires kernel version 3.8+ and should work uniformly
// across all supported devices.

package generic

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"syscall"
	"time"

	"github.com/kidoman/embd"
)

type digitalPin struct {
	id string
	n  int

	drv embd.GPIODriver

	dir       *os.File
	val       *os.File
	activeLow *os.File

	readBuf []byte

	initialized bool
}

func NewDigitalPin(pd *embd.PinDesc, drv embd.GPIODriver) embd.DigitalPin {
	return &digitalPin{id: pd.ID, n: pd.DigitalLogical, drv: drv, readBuf: make([]byte, 1)}
}

func (p *digitalPin) N() int {
	return p.n
}

func (p *digitalPin) init() error {
	if p.initialized {
		return nil
	}

	var err error
	if err = p.export(); err != nil {
		return err
	}
	if p.dir, err = p.directionFile(); err != nil {
		return err
	}
	if p.val, err = p.valueFile(); err != nil {
		return err
	}
	if p.activeLow, err = p.activeLowFile(); err != nil {
		return err
	}

	p.initialized = true

	return nil
}

func (p *digitalPin) export() error {
	exporter, err := os.OpenFile("/sys/class/gpio/export", os.O_WRONLY, os.ModeExclusive)
	if err != nil {
		return err
	}
	defer exporter.Close()
	_, err = exporter.WriteString(strconv.Itoa(p.n))
	if e, ok := err.(*os.PathError); ok && e.Err == syscall.EBUSY {
		return nil // EBUSY -> the pin has already been exported
	}
	return err
}

func (p *digitalPin) unexport() error {
	unexporter, err := os.OpenFile("/sys/class/gpio/unexport", os.O_WRONLY, os.ModeExclusive)
	if err != nil {
		return err
	}
	defer unexporter.Close()
	_, err = unexporter.WriteString(strconv.Itoa(p.n))
	return err
}

func (p *digitalPin) basePath() string {
	return fmt.Sprintf("/sys/class/gpio/gpio%v", p.n)
}

func (p *digitalPin) openFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDWR, os.ModeExclusive)
}

func (p *digitalPin) directionFile() (*os.File, error) {
	return p.openFile(path.Join(p.basePath(), "direction"))
}

func (p *digitalPin) valueFile() (*os.File, error) {
	return p.openFile(path.Join(p.basePath(), "value"))
}

func (p *digitalPin) activeLowFile() (*os.File, error) {
	return p.openFile(path.Join(p.basePath(), "active_low"))
}

func (p *digitalPin) SetDirection(dir embd.Direction) error {
	if err := p.init(); err != nil {
		return err
	}

	str := "in"
	if dir == embd.Out {
		str = "out"
	}
	_, err := p.dir.WriteString(str)
	return err
}

func (p *digitalPin) read() (int, error) {
	if _, err := p.val.ReadAt(p.readBuf, 0); err != nil {
		return 0, err
	}
	if p.readBuf[0] == 49 {
		return 1, nil
	}
	return 0, nil
}

func (p *digitalPin) Read() (int, error) {
	if err := p.init(); err != nil {
		return 0, err
	}

	return p.read()
}

var (
	lowBytes  = []byte{48}
	highBytes = []byte{49}
)

func (p *digitalPin) write(val int) error {
	bytes := lowBytes
	if val == embd.High {
		bytes = highBytes
	}
	_, err := p.val.Write(bytes)
	return err
}

func (p *digitalPin) Write(val int) error {
	if err := p.init(); err != nil {
		return err
	}

	return p.write(val)
}

func (p *digitalPin) TimePulse(state int) (time.Duration, error) {
	if err := p.init(); err != nil {
		return 0, err
	}

	aroundState := embd.Low
	if state == embd.Low {
		aroundState = embd.High
	}

	// Wait for any previous pulse to end
	for {
		v, err := p.read()
		if err != nil {
			return 0, err
		}

		if v == aroundState {
			break
		}
	}

	// Wait until ECHO goes high
	for {
		v, err := p.read()
		if err != nil {
			return 0, err
		}

		if v == state {
			break
		}
	}

	startTime := time.Now() // Record time when ECHO goes high

	// Wait until ECHO goes low
	for {
		v, err := p.read()
		if err != nil {
			return 0, err
		}

		if v == aroundState {
			break
		}
	}

	return time.Since(startTime), nil // Calculate time lapsed for ECHO to transition from high to low
}

func (p *digitalPin) ActiveLow(b bool) error {
	if err := p.init(); err != nil {
		return err
	}

	str := "0"
	if b {
		str = "1"
	}
	_, err := p.activeLow.WriteString(str)
	return err
}

func (p *digitalPin) PullUp() error {
	return errors.New("gpio: not implemented")
}

func (p *digitalPin) PullDown() error {
	return errors.New("gpio: not implemented")
}

func (p *digitalPin) Close() error {
	if err := p.StopWatching(); err != nil {
		return err
	}

	if err := p.drv.Unregister(p.id); err != nil {
		return err
	}

	if !p.initialized {
		return nil
	}

	if err := p.dir.Close(); err != nil {
		return err
	}
	if err := p.val.Close(); err != nil {
		return err
	}
	if err := p.activeLow.Close(); err != nil {
		return err
	}
	if err := p.unexport(); err != nil {
		return err
	}

	p.initialized = false

	return nil
}

func (p *digitalPin) setEdge(edge embd.Edge) error {
	file, err := p.openFile(path.Join(p.basePath(), "edge"))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte(edge))
	return err
}

func (p *digitalPin) Watch(edge embd.Edge, handler func(embd.DigitalPin)) error {
	if err := p.setEdge(edge); err != nil {
		return err
	}
	return registerInterrupt(p, handler)
}

func (p *digitalPin) StopWatching() error {
	return unregisterInterrupt(p)
}
