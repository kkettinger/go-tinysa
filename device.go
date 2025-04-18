package tinysa

import (
	"go.bug.st/serial"
	"log/slog"
	"sync"
	"time"
)

type Device struct {
	port            serial.Port
	mutex           sync.Mutex
	model           Model
	version         string
	hwVersion       string
	width           int
	height          int
	logger          *slog.Logger
	readTimeout     time.Duration
	responseTimeout time.Duration
}

// Close closes the open device.
func (d *Device) Close() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.logger.Info("closing port")
	return d.port.Close()
}

// Model returns the detected tinySA device model.
func (d *Device) Model() Model {
	return d.model
}

// Version returns the detected firmware version.
func (d *Device) Version() string {
	return d.version
}

// HardwareVersion returns the detected hardware version.
func (d *Device) HardwareVersion() string {
	return d.hwVersion
}

// ScreenResolution returns the screen width and height in pixels for the detected device model.
func (d *Device) ScreenResolution() (width, height int) {
	return d.width, d.height
}

// SendCommand sends a command to the device and returns the parsed response as string.
func (d *Device) SendCommand(cmd string) (string, error) {
	return d.sendCommand(cmd)
}

// SendCommandBinary sends a command to the device and returns the parsed response as []byte.
func (d *Device) SendCommandBinary(cmd string) ([]byte, error) {
	return d.sendCommandBinary(cmd)
}

// sendCommand is the internal method for requesting commands and returning a string response.
func (d *Device) sendCommand(cmd string) (string, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	res, err := sendCommand(d.logger, d.port, cmd, d.responseTimeout)
	if err != nil {
		return "", err
	}

	return res, nil
}

// sendCommand is the internal method for requesting commands and returning a binary response.
func (d *Device) sendCommandBinary(cmd string) ([]byte, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	res, err := sendCommandBinary(d.logger, d.port, cmd, d.responseTimeout)
	if err != nil {
		return []byte{}, err
	}

	return res, nil
}
