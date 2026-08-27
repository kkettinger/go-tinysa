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

	result, err := dev.SendCommand("version")
	if err != nil {
		panic(err)
	}

	fmt.Println("Result:", result)
}
