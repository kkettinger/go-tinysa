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

// EnableLNA enables the low noise amplifier.
func (d *Device) EnableLNA() error {
	d.logger.Info("enabling lna")
	_, err := d.sendCommand("lna on")
	return err
}

// DisableLNA disables the low noise amplifier.
func (d *Device) DisableLNA() error {
	d.logger.Info("disabling lna")
	_, err := d.sendCommand("lna off")
	return err
}
