// Copyright 2016 by Thorsten von Eicken

// The RFM69 package interfaces with a HopeRF RFM69 radio connected to an SPI bus. In addition,
// an interrupt capable GPIO pin may be used to avoid having to poll the radio.
package rfm69

import (
	"fmt"
	"log"

	"github.com/kidoman/embd"
)

// rfm69 represents a HopeRF RFM69 radio
type rfm69 struct {
	// configuration
	spi     embd.SPIBus       // bus where the radio is connected
	intrPin embd.InterruptPin // interrupt pin for RX and TX interrupts
	id      byte              // my RF ID/address
	group   byte              // RF address of group
	freq    int               // center frequency
	parity  byte              // ???
	// state
	mode byte // current operation mode
	// info about current RX packet
	rxInfo *RxInfo
}

type Packet struct {
	Length  uint8 // number of message bytes plus 1 for the address byte
	Address uint8 // destination address
	Message []byte
}

type RxInfo struct {
	rssi int // rssi value for current packet
	lna  int // low noise amp gain for current packet
	fei  int // frequency error for current packet
	afc  int // frequency correction applied for current packet
}

// New creates a connection to an rfm69 radio connected to the provided SPI bus and interrupt pin.
// the bufCount determines how many transmit buffers are allocated to allow for the queueing of
// transmit packets.
// For the RFM69 the SPI bus must be set to 10Mhz and mode 0.
func New(bus embd.SPIBus, intr embd.InterruptPin, id, group byte, freq int) *rfm69 {
	// bit 7 = b7^b5^b3^b1; bit 6 = b6^b4^b2^b0
	parity := group ^ (group << 4)
	parity = (parity ^ (parity << 2)) & 0xc0
	return &rfm69{spi: bus, intrPin: intr, id: id, group: group, freq: freq, parity: parity,
		mode: 255}
}

func (rf *rfm69) writeReg(addr, data byte) error {
	buf := []byte{addr | 0x80, data}
	return rf.spi.TransferAndReceiveData(buf)
}

func (rf *rfm69) readReg(addr byte) (byte, error) {
	buf := []byte{addr & 0x7f, 0}
	err := rf.spi.TransferAndReceiveData(buf)
	return buf[1], err
}

func (rf *rfm69) Init() error {
	// try to establish communication with the rfm69
	sync := func(pattern byte) error {
		n := 10
		for {
			rf.writeReg(REG_SYNCVALUE1, pattern)
			v, err := rf.readReg(REG_SYNCVALUE1)
			if err != nil {
				return err
			}
			if v == pattern {
				return nil
			}
			if n == 0 {
				return fmt.Errorf("Cannot sync with rfm69 chip")
			}
			n--
		}
	}
	if err := sync(0xaa); err != nil {
		return err
	}
	if err := sync(0x55); err != nil {
		return err
	}

	// write the configuration into the registers
	for i := 0; i < len(configRegs)-1; i += 2 {
		if err := rf.writeReg(configRegs[i], configRegs[i+1]); err != nil {
			return err
		}
	}

	rf.setFrequency(rf.freq)
	rf.writeReg(REG_SYNCVALUE2, rf.group)

	if gpio, ok := rf.intrPin.(embd.DigitalPin); ok {
		log.Printf("Set intr direction")
		gpio.SetDirection(embd.In)
	}
	if err := rf.intrPin.Watch(embd.EdgeRising, rf.intrHandler); err != nil {
		return err
	}

	return nil
}

func (rf *rfm69) setFrequency(freq int) {
	// accept any frequency scale as input, including KHz and MHz
	// multiply by 10 until freq >= 100 MHz
	for freq > 0 && freq < 100000000 {
		freq = freq * 10
	}

	// Frequency steps are in units of (32,000,000 >> 19) = 61.03515625 Hz
	// use multiples of 64 to avoid multi-precision arithmetic, i.e. 3906.25 Hz
	// due to this, the lower 6 bits of the calculated factor will always be 0
	// this is still 4 ppm, i.e. well below the radio's 32 MHz crystal accuracy
	// 868.0 MHz = 0xD90000, 868.3 MHz = 0xD91300, 915.0 MHz = 0xE4C000
	frf := (freq << 2) / (32000000 >> 11)
	rf.writeReg(REG_FRFMSB, byte(frf>>10))
	rf.writeReg(REG_FRFMSB+1, byte(frf>>2))
	rf.writeReg(REG_FRFMSB+2, byte(frf<<6))
}

func (rf *rfm69) setMode(mode byte) error {
	reg, err := rf.readReg(REG_OPMODE)
	if err != nil {
		return err
	}
	reg = (reg & 0xE3) | mode
	err = rf.writeReg(REG_OPMODE, reg)
	if err != nil {
		return err
	}
	for {
		val, err := rf.readReg(REG_IRQFLAGS1)
		if err != nil {
			rf.mode = 255
			return err
		}
		if val&IRQ1_MODEREADY != 0 {
			rf.mode = mode
			return nil
		}
	}
}

func (rf *rfm69) Send(header byte, message []byte) error {
	if len(message) > 62 {
		return fmt.Errorf("message too long")
	}
	rf.setMode(MODE_SLEEP)

	buf := make([]byte, len(message)+4)
	buf[0] = REG_FIFO | 0x80
	buf[1] = byte(len(message) + 2)
	buf[2] = (header & 0x3f) | rf.parity
	buf[3] = (header & 0xC0) | rf.id
	copy(buf[4:], message)
	err := rf.spi.TransferAndReceiveData(buf)
	if err != nil {
		return err
	}
	rf.setMode(MODE_TRANSMIT)
	for {
		val, err := rf.readReg(REG_IRQFLAGS2)
		if err != nil {
			return err
		}
		if val&IRQ2_PACKETSENT != 0 {
			break
		}
	}
	rf.setMode(MODE_STANDBY)
	return nil
}

func (rf *rfm69) readInfo() *RxInfo {
	// collect rxinfo, start with rssi
	rxInfo := &RxInfo{}
	rssi, _ := rf.readReg(REG_RSSIVALUE)
	rxInfo.rssi = 0 - int(rssi)/2
	// low noise amp gain
	lna, _ := rf.readReg(REG_LNAVALUE)
	rxInfo.lna = int((lna >> 3) & 0x7)
	// auto freq correction applied, caution: signed value
	buf := []byte{REG_AFCMSB, 0, 0}
	rf.spi.TransferAndReceiveData(buf)
	rxInfo.afc = int(int8(buf[1]))<<8 | int(buf[2])
	// freq error detected, caution: signed value
	buf = []byte{REG_FEIMSB, 0, 0}
	rf.spi.TransferAndReceiveData(buf)
	rxInfo.fei = int(int8(buf[1]))<<8 | int(buf[2])
	return rxInfo
}

func (rf *rfm69) Receive() (header byte, message []byte, info *RxInfo, err error) {
	// if we're not in receive mode, then switch, this also flushes the FIFO
	if rf.mode != MODE_RECEIVE {
		rf.setMode(MODE_RECEIVE)
		return
	}

	// if we don't have rxinfo check whether we have RX_READY, which means that we've
	// started receiving a packet so we can collect info
	if rf.rxInfo == nil {
		irq1, err := rf.readReg(REG_IRQFLAGS1)
		if err != nil {
			return 0, nil, nil, err
		}
		if irq1&IRQ1_RXREADY != 0 {
			rf.rxInfo = rf.readInfo()
		}
	}

	// see whether we have a full packet
	irq2, err := rf.readReg(REG_IRQFLAGS2)
	if err != nil {
		return 0, nil, nil, err
	}
	if irq2&IRQ2_PAYLOADREADY == 0 {
		return
	}
	i2 := rf.readInfo()
	if rf.rxInfo != nil && i2 != nil &&
		(rf.rxInfo.rssi != i2.rssi || rf.rxInfo.lna != i2.lna ||
			rf.rxInfo.afc != i2.afc || rf.rxInfo.fei != i2.fei) {
		fmt.Printf("\nrxInfo mismatch: %+v vs %+v\n", *rf.rxInfo, *i2)
	}
	// got packet, read it by fetching the entire FIFO, should be faster than first
	// looking at the length
	buf := make([]byte, 67)
	buf[0] = REG_FIFO
	err = rf.spi.TransferAndReceiveData(buf)
	if err != nil {
		return 0, nil, nil, err
	}
	// return the packet
	info = rf.rxInfo
	rf.rxInfo = nil
	l := buf[1]
	if l > 66 {
		l = 66 // or error?
	}
	header = buf[2]
	message = buf[3 : 2+l]
	return
}

func (rf *rfm69) intrHandler(pin embd.DigitalPin) {
	log.Printf("Interrupt called!")
}
