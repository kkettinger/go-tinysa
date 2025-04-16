package tinysa

import "fmt"

// LoadPreset loads a configuration from internal storage of the device.
func (d *Device) LoadPreset(presetId uint) error {
	d.logger.Info("loading preset", "preset_id", presetId)
	_, err := d.sendCommand(fmt.Sprintf("load %d", presetId))
	return err
}

// SavePreset saves the current configuration to the internal storage of the device.
func (d *Device) SavePreset(presetId uint) error {
	d.logger.Info("saving preset", "preset_id", presetId)
	_, err := d.sendCommand(fmt.Sprintf("save %d", presetId))
	return err
}
