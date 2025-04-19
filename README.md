[![Golang Test](https://github.com/kkettinger/go-tinysa/actions/workflows/test.yml/badge.svg)](https://github.com/kkettinger/go-tinysa/actions/workflows/test.yml)
[![Golang CI Lint](https://github.com/kkettinger/go-tinysa/actions/workflows/lint.yml/badge.svg)](https://github.com/kkettinger/go-tinysa/actions/workflows/lint.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/kkettinger/go-tinysa.svg)](https://pkg.go.dev/github.com/kkettinger/go-tinysa)
[![License: MIT](https://img.shields.io/github/license/kkettinger/go-tinysa)](/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/kkettinger/go-tinysa)](https://goreportcard.com/report/github.com/kkettinger/go-tinysa)

# go-tinysa

`go-tinysa` is an SDK for controlling and interacting with the [tinySA](https://www.tinysa.org/) spectrum analyzer via
its USB serial interface. It contains methods that allows you to:

- Configure sweep parameters (frequency range, center, span, ...)
- Configure markers and traces
- Export screenshots as `image.Image`
- Export trace frequencies and values
- Open menus (e.g., enable waterfall view)
- Reset device (DFU mode for basic model)
- Load/save presets
- Send raw commands
- And more, check out the [go reference](https://pkg.go.dev/github.com/kkettinger/go-tinysa)

_Note:_ The SDK was developed by using a tinySA Ultra with firmware version `v1.4-197`.
If you encounter issues with the basic model or a different firmware version, please report them.

## Installation

To use the sdk, run `go get github.com/kkettinger/go-tinysa` inside your golang project folder.

## Usage

For quick and easy connection, use `FindDevice()` that iterates over all serial ports
and [probes](#probing--model-detection) each port for a tinySA device. The first valid device is used.

```go
dev, err := tinysa.FindDevice()
if err != nil { panic(err) }

fmt.Println("Model:", dev.Model())
fmt.Println("Version:", dev.Version())
fmt.Println("Hardware Version:", dev.HardwareVersion())

width, height := dev.ScreenResolution()
fmt.Println("Screen resolution:", width, height)
```

Options like baudrate or timeouts can be specified like this:

```go
dev, _ := tinysa.FindDevice(
    tinysa.WithBaudRate(9600),
    tinysa.WithReadTimeout(500 * time.Millisecond))
```

To directly connect to a device, use the `NewDevice()` method:

```go
dev, _ := tinysa.NewDevice("/dev/ttyACM0")
```

To have more insight about what happens inside, you can pass on a logger instance:

```go
logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))

dev, _ := tinysa.FindDevice(
    tinysa.WithLogger(logger))
```

Internally only `LevelInfo` and `LevelDebug` is used.

## Examples

Examples can be found in the [examples](examples) folder, which you can run directly with `go run ./examples/<example>`.

### Setting sweep parameters

```go
// Set sweep to 100Mhz to 120Mhz
dev.SetSweepStartStop(100e6, 120e6)

// Set span to 30.5Mhz
dev.SetSweepCenter(30.5e6)
```

All frequency arguments are `uint64` values specified in Hz.

### Getting trace data

```go
data, _ := dev.GetTraceData(1)
for _, d := range data {
    fmt.Println(d.Frequency, " ", d.Value)
}
```

### Sending raw commands

If a method for a specific command is missing, you can always send raw commands:

```go
result, _ := dev.SendCommand("sd_list")
fmt.Println("Result:", result)
```

### Capturing the screen

```go
img, _ := dev.Capture()
fmt.Println(img.Bounds())
```

## Probing / Model detection

The `FindDevice()` and `NewDevice()` methods both probe the serial device by issuing a `version` command and trying to
parse the response:

```go
tinySA4_v1.4-197-gaa78ccc
HW Version:V0.4.5.1
```

The prefix of the version response, e.g. `tinySA` or `tinySA4`, is used to decide if it's a basic or ultra model.
The probe result can then be access with `Model()`, `Version()` and `HardwareVersion()`.
The `ScreenResolution()` method will return the width and height of the screen based on the model.

The model is used for methods and options that are only valid for specific tinySA models.
For example, the `dfu` argument in `Reset(dfu bool)` is only valid for the basic model, and will return an
`ErrOptionNotSupportedByModel` error when the method is called by an ultra device.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgments

- The tinySA team for creating a great spectrum analyzer
- Contributors to the Go serial library
