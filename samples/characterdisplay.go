// +build ignore

package main

import (
	"flag"
	"time"

	"github.com/kidoman/embd"
	"github.com/kidoman/embd/controller/hd44780"
	"github.com/kidoman/embd/interface/display/characterdisplay"

	_ "github.com/kidoman/embd/host/all"
)

func main() {
	flag.Parse()

	if err := embd.InitI2C(); err != nil {
		panic(err)
	}
	defer embd.CloseI2C()

	bus := embd.NewI2CBus(1)

	controller, err := hd44780.NewI2C(
		bus,
		0x20,
		hd44780.PCF8574PinMap,
		hd44780.RowAddress20Col,
		hd44780.TwoLine,
		hd44780.BlinkOn,
	)
	if err != nil {
		panic(err)
	}

	display := characterdisplay.New(controller, 20, 4)
	defer display.Close()

	display.Clear()
	display.Message("Hello, world!\n@embd | characterdisplay")
	time.Sleep(10 * time.Second)
	display.BacklightOff()
}
