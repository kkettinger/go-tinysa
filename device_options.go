package tinysa

import (
	"log/slog"
	"time"
)

type deviceOptions struct {
	logger *slog.Logger

	// baudrate is the serial port baudrate.
	baudrate int

	// readTimeout is the maximum time to wait for the initial read operation.
	readTimeout time.Duration

	// responseTimeout is the maximum time to wait for the full response.
	responseTimeout time.Duration
}

// defaultDeviceOptions returns a deviceOptions struct initialized with default values.
func defaultDeviceOptions() deviceOptions {
	return deviceOptions{
		logger:          nil,
		baudrate:        115200,
		readTimeout:     500 * time.Millisecond,
		responseTimeout: 2000 * time.Millisecond,
	}
}

// DeviceOption defines a function type that modifies deviceOptions.
type DeviceOption func(*deviceOptions)

// WithBaudRate sets a custom baud rate for the device communication.
func WithBaudRate(baudRate int) DeviceOption {
	return func(opts *deviceOptions) {
		opts.baudrate = baudRate
	}
}

// WithLogger sets a custom slog.Logger for logging device interactions.
func WithLogger(logger *slog.Logger) DeviceOption {
	return func(opts *deviceOptions) {
		opts.logger = logger
	}
}

// WithReadTimeout sets the timeout duration for reading from the device.
func WithReadTimeout(timeout time.Duration) DeviceOption {
	return func(opts *deviceOptions) {
		opts.readTimeout = timeout
	}
}

// WithResponseTimeout sets the timeout duration for waiting for a device response.
func WithResponseTimeout(timeout time.Duration) DeviceOption {
	return func(opts *deviceOptions) {
		opts.responseTimeout = timeout
	}
}
