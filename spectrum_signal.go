package tinysa

// EnableSpurRemoval enables spur removal.
func (d *Device) EnableSpurRemoval() error {
	d.logger.Info("enabling spur removal")
	_, err := d.sendCommand("spur on")
	return err
}

// DisableSpurRemoval disables spur removal.
func (d *Device) DisableSpurRemoval() error {
	d.logger.Info("disabling spur removal")
	_, err := d.sendCommand("spur off")
	return err
}

// EnableAutoSpurRemoval sets spur removal to auto.
func (d *Device) EnableAutoSpurRemoval() error {
	d.logger.Info("enabling auto spur removal")
	_, err := d.sendCommand("spur auto")
	return err
}
