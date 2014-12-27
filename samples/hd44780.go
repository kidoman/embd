// +build ignore

package main

import (
	"flag"
	"time"

	"github.com/kidoman/embd"
	"github.com/kidoman/embd/controller/hd44780"

	_ "github.com/kidoman/embd/host/all"
)

func main() {
	flag.Parse()

	if err := embd.InitI2C(); err != nil {
		panic(err)
	}
	defer embd.CloseI2C()

	bus := embd.NewI2CBus(1)

	display, err := hd44780.NewI2CCharacterDisplay(
		bus,
		0x20,
		hd44780.PCF8574PinMap,
		20,
		4,
		hd44780.TwoLine,
		hd44780.BlinkOn,
	)
	if err != nil {
		panic(err)
	}
	defer display.Close()

	display.Clear()
	display.Message("Hello, world!\n@embd")
	time.Sleep(10 * time.Second)
	display.BacklightOff()
}
