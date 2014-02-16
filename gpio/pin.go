package gpio

import (
	"fmt"
	"os"
	"path"
)

type pin struct {
	n int

	dir       *os.File
	val       *os.File
	activeLow *os.File
	edge      *os.File
}

func NewPin(n int) (p *pin, err error) {
	p = &pin{n: n}
	err = p.init()
	return
}

func (p *pin) init() (err error) {
	if p.dir, err = p.directionFile(); err != nil {
		return
	}
	if p.val, err = p.valueFile(); err != nil {
		return
	}
	if p.activeLow, err = p.activeLowFile(); err != nil {
		return
	}
	if p.edge, err = p.edgeFile(); err != nil {
		return
	}

	return
}

func (p *pin) basePath() string {
	return fmt.Sprintf("/sys/class/gpio/gpio%v", p.n)
}

func (p *pin) openFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDWR, os.ModeExclusive)
}

func (p *pin) directionPath() string {
	return path.Join(p.basePath(), "direction")
}

func (p *pin) directionFile() (*os.File, error) {
	return p.openFile(p.directionPath())
}

func (p *pin) valuePath() string {
	return path.Join(p.basePath(), "value")
}

func (p *pin) valueFile() (*os.File, error) {
	return p.openFile(p.valuePath())
}

func (p *pin) activeLowPath() string {
	return path.Join(p.basePath(), "active_low")
}

func (p *pin) activeLowFile() (*os.File, error) {
	return p.openFile(p.activeLowPath())
}

func (p *pin) edgePath() string {
	return path.Join(p.basePath(), "edge")
}

func (p *pin) edgeFile() (*os.File, error) {
	return p.openFile(p.edgePath())
}

func (p *pin) Mode(dir Direction) (err error) {
	str := "in"
	if dir == Output {
		str = "out"
	}
	_, err = p.dir.WriteString(str)
	return
}

func (p *pin) Input() error {
	return p.Mode(Input)
}

func (p *pin) Output() error {
	return p.Mode(Output)
}

func (p *pin) Read() (s State, err error) {
	buf := make([]byte, 1)
	if _, err = p.val.Read(buf); err != nil {
		return
	}
	s = Low
	if buf[0] == '1' {
		s = High
	}
	return
}

func (p *pin) Write(s State) (err error) {
	str := "0"
	if s == High {
		str = "1"
	}
	_, err = p.val.WriteString(str)
	return
}

func (p *pin) Low() error {
	return p.Write(Low)
}

func (p *pin) High() error {
	return p.Write(High)
}

func (p *pin) SetActiveLow(b bool) (err error) {
	str := "0"
	if b {
		str = "1"
	}
	_, err = p.activeLow.WriteString(str)
	return
}

func (p *pin) ActiveLow() error {
	return p.SetActiveLow(true)
}

func (p *pin) ActiveHigh() error {
	return p.SetActiveLow(false)
}

func (p *pin) Close() {
	p.dir.Close()
	p.val.Close()
	p.activeLow.Close()
	p.edge.Close()
}
