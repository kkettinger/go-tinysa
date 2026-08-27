//revive:disable:package-comments
package main

import (
	"fmt"
	"github.com/kkettinger/go-tinysa"
	"time"
)

func main() {
	dev, err := tinysa.FindDevice(
		tinysa.WithBaudRate(9600),
		tinysa.WithReadTimeout(500*time.Millisecond))
	if err != nil {
		panic(err)
	}

	fmt.Println("Model:", dev.Model())
	fmt.Println("Version:", dev.Version())
	fmt.Println("Hardware Version:", dev.HardwareVersion())

	width, height := dev.ScreenResolution()
	fmt.Println("Screen resolution:", width, height)
}
