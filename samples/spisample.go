package main

import (
	"fmt"

	"github.com/kidoman/embd"

	_ "github.com/kidoman/embd/host/all"
)

func main() {
	var err error
	err = embd.InitSPI()
	if err != nil {
		panic(err)
	}
	defer embd.CloseSPI()

	spiBus := embd.NewSPIBus(embd.SPI_MODE_0, 0, 1000000, 8, 0)
	defer spiBus.Close()

	dataBuf := make([]uint8, 3)
	dataBuf[0] = uint8(1)
	dataBuf[1] = uint8(2)
	dataBuf[2] = uint8(3)

	err = spiBus.TransferAndRecieveData(dataBuf)
	if err != nil {
		panic(err)
	}

	fmt.Println("Recived data is: %v", dataBuf)

	dataReceived, err := spiBus.ReceiveData(3)
	if err != nil {
		panic(err)
	}

	fmt.Println("Recived data is: %v", dataReceived)

	dataByte := byte(1)
	receivedByte, err := spiBus.TransferAndReceiveByte(dataByte)
	if err != nil {
		panic(err)
	}
	fmt.Println("Recived byte is: %v", receivedByte)

	receivedByte, err = spiBus.ReceiveByte()
	if err != nil {
		panic(err)
	}
	fmt.Println("Recived byte is: %v", receivedByte)
}
