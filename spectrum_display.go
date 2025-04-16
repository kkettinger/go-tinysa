package tinysa

import "fmt"

type DisplayUnit struct {
	value string
}

func (u DisplayUnit) String() string {
	return u.value
}

func (u DisplayUnit) IsValid() bool {
	return u.value != ""
}

const (
	displayUnitRAW  string = "RAW"
	displayUnitDBm  string = "dBm"
	displayUnitDBmV string = "dBmV"
	displayUnitDBuV string = "dBuV"
	displayUnitV    string = "V"
	displayUnitVpp  string = "Vpp"
	displayUnitW    string = "W"
)

var (
	DisplayUnitRaw  = DisplayUnit{displayUnitRAW}
	DisplayUnitDBm  = DisplayUnit{displayUnitDBm}
	DisplayUnitDBmV = DisplayUnit{displayUnitDBmV}
	DisplayUnitDBuV = DisplayUnit{displayUnitDBuV}
	DisplayUnitV    = DisplayUnit{displayUnitV}
	DisplayUnitVpp  = DisplayUnit{displayUnitVpp}
	DisplayUnitW    = DisplayUnit{displayUnitW}
)

var displayUnitMap = map[string]DisplayUnit{
	displayUnitRAW:  DisplayUnitRaw,
	displayUnitDBm:  DisplayUnitDBm,
	displayUnitDBmV: DisplayUnitDBmV,
	displayUnitDBuV: DisplayUnitDBuV,
	displayUnitV:    DisplayUnitV,
	displayUnitVpp:  DisplayUnitVpp,
	displayUnitW:    DisplayUnitW,
}

var displayUnitOptions = []string{
	displayUnitRAW,
	displayUnitDBm,
	displayUnitDBmV,
	displayUnitDBuV,
	displayUnitV,
	displayUnitVpp,
	displayUnitW,
}

func DisplayUnitOptions() []string {
	return displayUnitOptions
}

// SetDisplayUnit sets the display unit to the specified value.
func (d *Device) SetDisplayUnit(unit DisplayUnit) error {
	d.logger.Info("setting display unit", "unit", unit)
	_, err := d.sendCommand(fmt.Sprintf("trace %s", unit.value))
	return err
}

// SetDisplayRefLevel sets the display ref level to the specified value in dBm.
func (d *Device) SetDisplayRefLevel(levelDbm int) error {
	d.logger.Info("setting display ref level", "level", levelDbm)
	_, err := d.sendCommand(fmt.Sprintf("trace reflevel %d", levelDbm))
	return err
}

// SetDisplayRefLevelAuto sets the display ref level to auto.
func (d *Device) SetDisplayRefLevelAuto() error {
	d.logger.Info("setting display ref level auto")
	_, err := d.sendCommand("trace reflevel auto")
	return err
}

// SetDisplayScale sets the display scale to the specified value.
func (d *Device) SetDisplayScale(level int) error {
	d.logger.Info("setting display scale", "level", level)
	_, err := d.sendCommand(fmt.Sprintf("trace scale %d", level))
	return err
}

// SetDisplayScaleAuto sets the display scale to auto.
// TODO: doesn't seem to work
func (d *Device) SetDisplayScaleAuto() error {
	d.logger.Info("setting display scale auto")
	_, err := d.sendCommand("trace scale auto")
	return err
}
