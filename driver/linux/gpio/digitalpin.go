package gpio

import (
	"fmt"
	"os"
	"path"

	"github.com/kidoman/embd/gpio"
)

type digitalPin struct {
	n int

	dir       *os.File
	val       *os.File
	activeLow *os.File
	edge      *os.File
}

func newDigitalPin(n int) (p *digitalPin, err error) {
	p = &digitalPin{n: n}
	err = p.init()
	return
}

func (p *digitalPin) init() (err error) {
	if p.dir, err = p.directionFile(); err != nil {
		return
	}
	if p.val, err = p.valueFile(); err != nil {
		return
	}
	if p.activeLow, err = p.activeLowFile(); err != nil {
		return
	}

	return
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

func (p *digitalPin) SetDir(dir gpio.Direction) (err error) {
	str := "in"
	if dir == gpio.Out {
		str = "out"
	}
	_, err = p.dir.WriteString(str)
	return
}

func (p *digitalPin) Read() (val int, err error) {
	buf := make([]byte, 1)
	if _, err = p.val.Read(buf); err != nil {
		return
	}
	val = 0
	if buf[0] == '1' {
		val = 1
	}
	return
}

func (p *digitalPin) Write(val int) (err error) {
	str := "0"
	if val == gpio.High {
		str = "1"
	}
	_, err = p.val.WriteString(str)
	return
}

func (p *digitalPin) ActiveLow(b bool) (err error) {
	str := "0"
	if b {
		str = "1"
	}
	_, err = p.activeLow.WriteString(str)
	return
}

func (p *digitalPin) Close() error {
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

	return nil
}
