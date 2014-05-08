package mcp3008

import (
	"github.com/golang/glog"
	"github.com/kidoman/embd"

	_ "github.com/kidoman/embd/host/all"
)

type mcp3008 struct {
	mode byte

	bus embd.SPIBus
}

var SingleMode byte = 1
var DifferenceMode byte = 0

func New(mode byte, bus embd.SPIBus) *mcp3008 {
	return &mcp3008{mode, bus}
}

func (m *mcp3008) AnalogValueAt(chanNum int) (int, error) {
	data := make([]uint8, 3)
	data[0] = 1
	data[1] = uint8(m.mode)<<7 | uint8(chanNum)<<4
	data[2] = 0

	glog.V(2).Infof("mcp3008: sendingdata buffer %v", data)
	if err := m.bus.TransferAndRecieveData(data); err != nil {
		return 0, err
	}

	return int(uint16(data[1]&0x03)<<8 | uint16(data[2])), nil
}
