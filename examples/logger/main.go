package main

import (
	"fmt"
	"github.com/kkettinger/go-tinysa"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	dev, err := tinysa.FindDevice(
		tinysa.WithLogger(logger))
	if err != nil {
		panic(err)
	}

	fmt.Println("Model:", dev.Model())
	fmt.Println("Version:", dev.Version())
	fmt.Println("Hardware Version:", dev.HardwareVersion())
}
