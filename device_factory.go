package tinysa

import (
	"fmt"
	"go.bug.st/serial"
)

// NewDevice creates a *Device from the specified port name.
func NewDevice(portName string, opts ...DeviceOption) (*Device, error) {
	options := defaultDeviceOptions()
	for _, opt := range opts {
		opt(&options)
	}

	if options.logger == nil {
		options.logger = NewNoopLogger()
	}

	logger := options.logger

	logger.Debug("initializing new device", "options", options)

	// open serial port
	mode := &serial.Mode{BaudRate: options.baudrate}
	logger.Debug("opening port", "port", portName, "baudrate", options.baudrate)
	port, err := serial.Open(portName, mode)
	if err != nil {
		logger.Error("failed to open port", "err", err)
		return nil, fmt.Errorf("failed to open port %s: %s", portName, err.Error())
	}

	// set read timeout
	if err = port.SetReadTimeout(options.readTimeout); err != nil {
		logger.Error("failed to set read timeout", "err", err)
		return nil, fmt.Errorf("failed to set read timeout: %s", err.Error())
	}

	// probe device
	logger.Debug("probing device", "port", port)
	pr, err := probeDevice(logger, port, options.responseTimeout)
	if err != nil {
		logger.Error("failed to probe device", "err", err)
		return nil, fmt.Errorf("failed to probe device: %s", err.Error())
	}

	return createDeviceFromProbe(logger, port, pr, options)
}

// FindDevice iterates over all serial ports and creates a *Device from the first valid tinySA device found.
func FindDevice(opts ...DeviceOption) (*Device, error) {
	options := defaultDeviceOptions()
	for _, opt := range opts {
		opt(&options)
	}

	if options.logger == nil {
		options.logger = NewNoopLogger()
	}

	logger := options.logger

	logger.Debug("finding device", "options", options)

	// list serial ports
	ports, err := serial.GetPortsList()
	if err != nil {
		logger.Error("failed to list serial ports", "err", err)
		return nil, fmt.Errorf("failed to list serial ports: %s", err.Error())
	}
	logger.Debug("list serial ports", "ports", ports)

	// try to create a device from each one
	for _, portName := range ports {
		logger.Debug("trying to create device", "port", portName)
		device, err := NewDevice(portName, opts...)
		if err != nil {
			logger.Error("failed to create device, skipping", "port", portName, "err", err)
			continue
		}
		return device, nil
	}

	return nil, fmt.Errorf("no device found")
}
