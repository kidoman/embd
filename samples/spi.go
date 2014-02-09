package main

import (
	"fmt"
	"github.com/kidoman/embd/spi"
)

func main() {
	var rx_data uint8

	fmt.Println("Hello")

	bus, _ := spi.NewSpiBus()

	rx_data, _ = bus.TransferAndRecieveByteData(8`	)

	fmt.Printf("Received %v \n", rx_data)
}
