// +build ignore

// PWM example, works OOTB on a BBB.

package main

import (
	"flag"
	"time"

	"github.com/kidoman/embd"
)

func main() {
	flag.Parse()

	if err := embd.InitGPIO(); err != nil {
		panic(err)
	}
	defer embd.CloseGPIO()

	pwm, err := embd.NewPWMPin("P9_14")
	if err != nil {
		panic(err)
	}
	defer pwm.Close()

	if err := pwm.SetDuty(embd.BBBPWMDefaultPeriod / 2); err != nil {
		panic(err)
	}

	time.Sleep(1 * time.Second)
}
