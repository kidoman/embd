package mcp3008

import (
	"github.com/golang/glog"
	"github.com/kidoman/embd"
)

// MCP3008 represents a mcp3008 8bit DAC
type MCP3008 struct {
	mode byte

	bus embd.SPIBus
}

var SingleMode byte = 1
var DifferenceMode byte = 0

func New(mode byte, bus embd.SPIBus) *MCP3008 {
	return &MCP3008{mode, bus}
}

const (
	startBit = 1
)

func (m *MCP3008) AnalogValueAt(chanNum int) (int, error) {
	var data [3]uint8
	data[0] = startBit
	data[1] = uint8(m.mode)<<7 | uint8(chanNum)<<4
	data[2] = 0

	glog.V(2).Infof("mcp3008: sendingdata buffer %v", data)
	if err := m.bus.TransferAndRecieveData(data[:]); err != nil {
		return 0, err
	}

	return int(uint16(data[1]&0x03)<<8 | uint16(data[2])), nil
}
