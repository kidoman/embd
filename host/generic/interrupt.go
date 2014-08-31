// Generic Interrupt Pins.

package generic

import (
	"errors"
	"fmt"
	"sync"
	"syscall"

	"github.com/kidoman/embd"
)

const (
	MaxGPIOInterrupt = 64
)

var ErrorPinAlreadyRegistered = errors.New("pin interrupt already registered")

type interrupt struct {
	pin            embd.DigitalPin
	initialTrigger bool
	handler        func(embd.DigitalPin)
}

func (i *interrupt) Signal() {
	if !i.initialTrigger {
		i.initialTrigger = true
		return
	}
	i.handler(i.pin)
}

type ePollListener struct {
	epollFd           int
	interruptablePins map[int]*interrupt
	mu                sync.Mutex
}

var ePollListenerInstance *ePollListener

func getEPollListenerInstance() *ePollListener {
	if ePollListenerInstance == nil {
		ePollListenerInstance = initEPollListener()
	}
	return ePollListenerInstance
}

func initEPollListener() *ePollListener {
	epollFd, err := syscall.EpollCreate1(0)
	if err != nil {
		panic(fmt.Sprintf("Unable to create epoll: %v", err))
	}
	listener := &ePollListener{epollFd: epollFd, interruptablePins: make(map[int]*interrupt)}

	go func() {
		var epollEvents [MaxGPIOInterrupt]syscall.EpollEvent

		for {
			n, err := syscall.EpollWait(listener.epollFd, epollEvents[:], -1)
			if err != nil {
				panic(fmt.Sprintf("EpollWait error: %v", err))
			}
			for i := 0; i < n; i++ {
				if irq, ok := listener.interruptablePins[int(epollEvents[i].Fd)]; ok {
					irq.Signal()
				}
			}
		}
	}()
	return listener
}

func registerInterrupt(pin *digitalPin, handler func(embd.DigitalPin)) error {
	l := getEPollListenerInstance()

	pinFd := int(pin.val.Fd())

	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.interruptablePins[pinFd]; ok {
		return ErrorPinAlreadyRegistered
	}

	var event syscall.EpollEvent
	event.Events = syscall.EPOLLIN | (syscall.EPOLLET & 0xffffffff) | syscall.EPOLLPRI

	if err := syscall.SetNonblock(pinFd, true); err != nil {
		return err
	}

	event.Fd = int32(pinFd)

	if err := syscall.EpollCtl(l.epollFd, syscall.EPOLL_CTL_ADD, pinFd, &event); err != nil {
		return err
	}

	l.interruptablePins[pinFd] = &interrupt{pin: pin, handler: handler}

	return nil
}

func unregisterInterrupt(pin *digitalPin) error {
	l := getEPollListenerInstance()

	pinFd := int(pin.val.Fd())

	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.interruptablePins[pinFd]; !ok {
		return nil
	}

	if err := syscall.EpollCtl(l.epollFd, syscall.EPOLL_CTL_DEL, pinFd, nil); err != nil {
		return err
	}

	if err := syscall.SetNonblock(pinFd, false); err != nil {
		return err
	}

	delete(l.interruptablePins, pinFd)
	return nil
}
