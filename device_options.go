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

func defaultDeviceOptions() deviceOptions {
	return deviceOptions{
		logger:          nil,
		baudrate:        115200,
		readTimeout:     500 * time.Millisecond,
		responseTimeout: 2000 * time.Millisecond,
	}
}

type DeviceOption func(*deviceOptions)

func WithBaudRate(baudRate int) DeviceOption {
	return func(opts *deviceOptions) {
		opts.baudrate = baudRate
	}
}

func WithLogger(logger *slog.Logger) DeviceOption {
	return func(opts *deviceOptions) {
		opts.logger = logger
	}
}

func WithReadTimeout(timeout time.Duration) DeviceOption {
	return func(opts *deviceOptions) {
		opts.readTimeout = timeout
	}
}

func WithResponseTimeout(timeout time.Duration) DeviceOption {
	return func(opts *deviceOptions) {
		opts.responseTimeout = timeout
	}
}
