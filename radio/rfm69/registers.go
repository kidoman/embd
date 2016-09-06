// Copyright 2016 by Thorsten von Eicken

package rfm69

const (
	REG_FIFO       = 0x00
	REG_OPMODE     = 0x01
	REG_FRFMSB     = 0x07
	REG_PALEVEL    = 0x11
	REG_LNAVALUE   = 0x18
	REG_AFCMSB     = 0x1F
	REG_AFCLSB     = 0x20
	REG_FEIMSB     = 0x21
	REG_FEILSB     = 0x22
	REG_RSSIVALUE  = 0x24
	REG_IRQFLAGS1  = 0x27
	REG_IRQFLAGS2  = 0x28
	REG_SYNCVALUE1 = 0x2F
	REG_SYNCVALUE2 = 0x30
	REG_NODEADDR   = 0x39
	REG_BCASTADDR  = 0x3A
	REG_FIFOTHRESH = 0x3C
	REG_PKTCONFIG2 = 0x3D
	REG_AESKEYMSB  = 0x3E

	MODE_SLEEP    = 0 << 2
	MODE_STANDBY  = 1 << 2
	MODE_TRANSMIT = 3 << 2
	MODE_RECEIVE  = 4 << 2

	START_TX = 0xC2
	STOP_TX  = 0x42

	RCCALSTART        = 0x80
	IRQ1_MODEREADY    = 1 << 7
	IRQ1_RXREADY      = 1 << 6
	IRQ1_SYNADDRMATCH = 1 << 0

	IRQ2_FIFONOTEMPTY = 1 << 6
	IRQ2_PACKETSENT   = 1 << 3
	IRQ2_PAYLOADREADY = 1 << 2
)

// register values to initialize the chip, this array has pairs of <address, data>
var configRegs = []byte{
	// POR value is better for first rf_sleep  0x01, 0x00, // OpMode = sleep
	0x02, 0x00, // DataModul = packet mode, fsk
	0x03, 0x02, // BitRateMsb, data rate = 49,261 khz
	0x04, 0x8A, // BitRateLsb, divider = 32 MHz / 650
	0x05, 0x02, // FdevMsb = 45 KHz
	0x06, 0xE1, // FdevLsb = 45 KHz
	0x0B, 0x20, // Low M
	0x19, 0x4A, // RxBw 100 KHz
	0x1A, 0x42, // AfcBw 125 KHz
	0x1E, 0x0C, // AfcAutoclearOn, AfcAutoOn
	//0x25, 0x40, //0x80, // DioMapping1 = SyncAddress (Rx)
	0x26, 0x07, // disable clkout
	0x29, 0xA0, // RssiThresh -80 dB
	0x2D, 0x05, // PreambleSize = 5
	0x2E, 0x88, // SyncConfig = sync on, sync size = 2
	0x2F, 0x2D, // SyncValue1 = 0x2D
	0x37, 0xD0, // PacketConfig1 = fixed, white, no filtering
	0x38, 0x42, // PayloadLength = 0, unlimited
	0x3C, 0x8F, // FifoTresh, not empty, level 15
	0x3D, 0x12, // 0x10, // PacketConfig2, interpkt = 1, autorxrestart off
	0x6F, 0x20, // TestDagc ...
	0x71, 0x02, // RegTestAfc
}
