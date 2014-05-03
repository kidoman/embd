package embd

import "sync"

type spiBusFactory func(byte, byte, byte, int, int, int, bool, func() error) SPIBus

type spiDriver struct {
	spiDevMinor      byte
	shouldInitialize bool
	initializer      func() error

	busMap     map[byte]SPIBus
	busMapLock sync.Mutex

	sbf spiBusFactory
}

func NewSPIDriver(spiDevMinor byte, sbf spiBusFactory, shouldInitialize bool, i func() error) SPIDriver {
	return &spiDriver{
		spiDevMinor:      spiDevMinor,
		sbf:              sbf,
		shouldInitialize: shouldInitialize,
		initializer:      i,
	}
}

func (s *spiDriver) Bus(mode, channel byte, speed, bpw, delay int) SPIBus {
	s.busMapLock.Lock()
	defer s.busMapLock.Unlock()

	b := s.sbf(s.spiDevMinor, mode, channel, speed, bpw, delay, s.shouldInitialize, s.initializer)
	s.busMap = make(map[byte]SPIBus)
	s.busMap[channel] = b
	return b
}

func (s *spiDriver) Close() error {
	for _, b := range s.busMap {
		b.Close()
	}

	return nil
}
