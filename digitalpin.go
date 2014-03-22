package embd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
)

type digitalPin struct {
	n int

	dir       *os.File
	val       *os.File
	activeLow *os.File
	edge      *os.File

	initialized bool
}

func newDigitalPin(n int) *digitalPin {
	return &digitalPin{n: n}
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

func (p *digitalPin) SetDirection(dir Direction) error {
	if err := p.init(); err != nil {
		return err
	}

	str := "in"
	if dir == Out {
		str = "out"
	}
	_, err := p.dir.WriteString(str)
	return err
}

func (p *digitalPin) Read() (int, error) {
	if err := p.init(); err != nil {
		return 0, err
	}

	buf := make([]byte, 1)
	if _, err := p.val.Read(buf); err != nil {
		return 0, err
	}
	var val int
	if buf[0] == '1' {
		val = 1
	}
	return val, nil
}

func (p *digitalPin) Write(val int) error {
	if err := p.init(); err != nil {
		return err
	}

	str := "0"
	if val == High {
		str = "1"
	}
	_, err := p.val.WriteString(str)
	return err
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
	if err := p.edge.Close(); err != nil {
		return err
	}
	if err := p.unexport(); err != nil {
		return err
	}

	p.initialized = false

	return nil
}
