// +build ignore

package main

import (
	"fmt"

	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/all"
)

func main() {
	if err := embd.InitSPI(); err != nil {
		panic(err)
	}
	defer embd.CloseSPI()

	spiBus := embd.NewSPIBus(embd.SPIMode0, 0, 1000000, 8, 0)
	defer spiBus.Close()

	dataBuf := [3]uint8{1, 2, 3}

	if err := spiBus.TransferAndReceiveData(dataBuf[:]); err != nil {
		panic(err)
	}

	fmt.Println("received data is:", dataBuf)

	dataReceived, err := spiBus.ReceiveData(3)
	if err != nil {
		panic(err)
	}

	fmt.Println("received data is:", dataReceived)

	dataByte := byte(1)
	receivedByte, err := spiBus.TransferAndReceiveByte(dataByte)
	if err != nil {
		panic(err)
	}
	fmt.Println("received byte is:", receivedByte)

	receivedByte, err = spiBus.ReceiveByte()
	if err != nil {
		panic(err)
	}
	fmt.Println("received byte is:", receivedByte)
}
