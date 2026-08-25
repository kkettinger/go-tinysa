//revive:disable:package-comments
package main

import (
	"fmt"

	"github.com/cjheath/go-tinysa"
)

// This connects to a tinySA that has its serial port bridged to TCP by socat, e.g.:
//
//	socat /dev/cu.usbmodem4001,raw,echo=0,ispeed=115200,ospeed=115200 TCP-LISTEN:9001,reuseaddr,fork
func main() {
	dev, err := tinysa.NewDeviceTCP("localhost:9001")
	if err != nil {
		panic(err)
	}

	fmt.Println("Model:", dev.Model())
	fmt.Println("Version:", dev.Version())
	fmt.Println("Hardware Version:", dev.HardwareVersion())

	width, height := dev.ScreenResolution()
	fmt.Println("Screen resolution:", width, height)
}
