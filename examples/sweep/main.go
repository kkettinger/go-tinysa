package main

import (
	"github.com/kkettinger/go-tinysa"
)

func main() {
	dev, err := tinysa.FindDevice()
	if err != nil {
		panic(err)
	}

	if err := dev.SetSweepStartStop(100e6, 120e6); err != nil {
		panic(err)
	}
}
