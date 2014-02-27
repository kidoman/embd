// Package i2c enables gophers i2c speaking ability.
package i2c

import (
	"sync"

	"github.com/kidoman/embd/i2c"
)

type I2C struct {
	busMap     map[byte]*bus
	busMapLock sync.Mutex
}

func New() *I2C {
	return &I2C{
		busMap: make(map[byte]*bus),
	}
}

func (i *I2C) Bus(l byte) i2c.Bus {
	i.busMapLock.Lock()
	defer i.busMapLock.Unlock()

	var b *bus

	if b = i.busMap[l]; b == nil {
		b = &bus{l: l}
		i.busMap[l] = b
	}

	return b
}

func (i *I2C) Close() error {
	for _, b := range i.busMap {
		b.Close()

		delete(i.busMap, b.l)
	}

	return nil
}
