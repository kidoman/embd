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

	hd, err := hd44780.NewI2C(
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
	defer hd.Close()

	hd.Clear()
	message := "Hello, world!"
	bytes := []byte(message)
	for _, b := range bytes {
		hd.WriteChar(b)
	}
	hd.SetCursor(0, 1)

	message = "@embd | hd44780"
	bytes = []byte(message)
	for _, b := range bytes {
		hd.WriteChar(b)
	}
	time.Sleep(10 * time.Second)
	hd.BacklightOff()
}
