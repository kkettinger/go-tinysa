package tinysa

import (
	"fmt"
	"strconv"
)

// GetBatteryVoltage returns battery voltage in mV.
func (d *Device) GetBatteryVoltage() (uint, error) {
	d.logger.Info("retrieving battery voltage")

	line, err := d.sendCommand("vbat")
	if err != nil {
		return 0, err
	}

	result, err := parseBatteryResponse(line)
	if err != nil {
		d.logger.Error("failed to parse battery voltage", "err", err, "line", line)
		return 0, fmt.Errorf("failed to parse battery voltage: %s", err.Error())
	}

	return result, nil
}

// GetBatteryOffsetVoltage gets battery offset voltage in mV.
func (d *Device) GetBatteryOffsetVoltage() (uint, error) {
	d.logger.Info("retrieving battery voltage")

	res, err := d.sendCommand("vbat_offset")
	if err != nil {
		return 0, err
	}

	vbatOffset, err := strconv.ParseUint(res, 10, 0)
	if err != nil {
		return 0, fmt.Errorf("failed to parse battery offset voltage: %s", err.Error())
	}

	return uint(vbatOffset), nil
}

// SetBatteryOffsetVoltage sets battery offset voltage in mV.
func (d *Device) SetBatteryOffsetVoltage(voltage uint) error {
	d.logger.Info("setting battery voltage", "voltage", voltage)
	_, err := d.sendCommand(fmt.Sprintf("vbat_offset %d", voltage))
	return err
}
