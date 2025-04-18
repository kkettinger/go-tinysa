package tinysa

import "fmt"

// LoadPreset loads a configuration from internal storage of the device.
func (d *Device) LoadPreset(presetID uint) error {
	d.logger.Info("loading preset", "preset_id", presetID)
	_, err := d.sendCommand(fmt.Sprintf("load %d", presetID))
	return err
}

// SavePreset saves the current configuration to the internal storage of the device.
func (d *Device) SavePreset(presetID uint) error {
	d.logger.Info("saving preset", "preset_id", presetID)
	_, err := d.sendCommand(fmt.Sprintf("save %d", presetID))
	return err
}
