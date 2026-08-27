//revive:disable:package-comments
package main

import (
	"fmt"
	"github.com/kkettinger/go-tinysa"
)

func main() {
	dev, err := tinysa.FindDevice()
	if err != nil {
		panic(err)
	}

	data, err := dev.GetTraceData(1)
	if err != nil {
		panic(err)
	}

	for _, d := range data {
		fmt.Println(d.Frequency, " ", d.Value)
	}
}
