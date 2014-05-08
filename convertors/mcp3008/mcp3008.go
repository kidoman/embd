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

func New(mode byte, spiChan, speed int) (*mcp3008, error) {
	if err := embd.InitSPI(); err != nil {
		return nil, err
	}
	glog.V(3).Infof("mcp3008: getting spiBus with mode: %v, channel: %v, speed: %v", mode, spiChan, speed)
	spiBus := embd.NewSPIBus(embd.SpiMode0, byte(spiChan), speed, 0, 0)
	return &mcp3008{mode, spiBus}, nil
}

func (m *mcp3008) AnalogValueAt(chanNum int) (int, error) {
	var data [3]uint8
	data[0] = 1
	data[1] = uint8(m.mode)<<7 | uint8(chanNum)<<4
	data[2] = 0

	if err := m.bus.TransferAndRecieveData(data); err != nil {
		return 0, err
	}

	return int(uint16(data[1]&0x03)<<8 | uint16(data[2])), nil
}

func (m *mcp3008) Close() error {
	glog.V(2).Infoln("mcp3008: performing cleanup")
	if err := m.bus.Close(); err != nil {
		return err
	}
	if err := embd.CloseSPI(); err != nil {
		return err
	}
	return nil
}
