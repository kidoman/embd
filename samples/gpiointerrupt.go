// +build ignore

package main

import (
	"flag"
	"fmt"

	"github.com/kidoman/embd"

	_ "github.com/kidoman/embd/host/all"
)

func main() {
	flag.Parse()

	if err := embd.InitGPIO(); err != nil {
		panic(err)
	}
	defer embd.CloseGPIO()

	btn, err := embd.NewDigitalPin(10)
	if err != nil {
		panic(err)
	}
	defer btn.Close()

	if err := btn.SetDirection(embd.In); err != nil {
		panic(err)
	}
	btn.ActiveLow(false)

	quit := make(chan interface{})
	err = btn.Watch(embd.EdgeFalling, func(btn embd.DigitalPin) {
		quit <- btn
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Button %v was pressed.\n", <-quit)
}
