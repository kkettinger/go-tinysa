package main

import (
	"fmt"
	"github.com/kkettinger/go-tinysa"
	"golang.org/x/image/bmp"
	"os"
)

func main() {
	output := "tinysa_screenshot.bmp"

	dev, err := tinysa.FindDevice()
	if err != nil {
		panic(err)
	}

	img, err := dev.Capture()
	if err != nil {
		panic(err)
	}

	file, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err = bmp.Encode(file, img); err != nil {
		panic(err)
	}

	fmt.Println("Screenshot saved to", output)
}
