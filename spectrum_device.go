package tinysa

import (
	"fmt"
	"strconv"
	"strings"
)

// GetVersion requests and returns the full version information.
func (d *Device) GetVersion() (string, error) {
	d.logger.Info("requesting device version")
	return d.sendCommand("version")
}

// Reset restarts the device, optionally entering DFU mode (supported only by the basic model).
func (d *Device) Reset(dfu bool) error {
	d.logger.Info("resetting device", "dfu", dfu)
	cmd := "reset"
	if dfu {
		if d.model != ModelBasic {
			return ErrCommandNotSupportedByModel
		}
		cmd = cmd + " dfu"
	}
	_, err := d.sendCommand(cmd)
	return err
}

// GetDeviceId returns the device id.
func (d *Device) GetDeviceId() (uint, error) {
	d.logger.Info("requesting device id")
	res, err := d.sendCommand("deviceid")
	if err != nil {
		return 0, fmt.Errorf("%w: failed to get device id: %v", ErrCommandFailed, err)
	}

	parts := strings.Split(res, " ")
	if len(parts) != 2 && parts[0] != "deviceid" {
		return 0, fmt.Errorf("%w: unexpected response for deviceid: %s", ErrCommandFailed, res)
	}
	id, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("%w: failed to parse device id: %v", ErrCommandFailed, err)
	}
	return uint(id), nil
}

// SetDeviceId sets the device id.
func (d *Device) SetDeviceId(id uint) error {
	d.logger.Info("setting device id", "id", id)
	_, err := d.sendCommand(fmt.Sprintf("deviceid %d", id))
	return err
}
