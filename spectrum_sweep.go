package tinysa

import (
	"fmt"
)

type Sweep struct {
	Start  uint64
	Stop   uint64
	Points uint
}

type SweepMode struct {
	value string
}

const (
	sweepModeNormal  string = "normal"
	sweepModePrecise string = "precise"
	sweepModeFast    string = "fast"
	sweepModeNoise   string = "noise"
)

var (
	SweepModeNormal  = SweepMode{sweepModeNormal}
	SweepModePrecise = SweepMode{sweepModePrecise}
	SweepModeFast    = SweepMode{sweepModeFast}
	SweepModeNoise   = SweepMode{sweepModeNoise}
)

var sweepModeMap = map[string]SweepMode{
	sweepModeNormal:  SweepModeNormal,
	sweepModePrecise: SweepModePrecise,
	sweepModeFast:    SweepModeFast,
	sweepModeNoise:   SweepModeNoise,
}

var sweepModeOptions = []string{
	sweepModeNormal,
	sweepModePrecise,
	sweepModeFast,
	sweepModeNoise,
}

func SweepModeOptions() []string {
	return sweepModeOptions
}

func (m SweepMode) String() string {
	return m.value
}

func (m SweepMode) IsValid() bool {
	return m.value != ""
}

type SweepStatus string

const (
	SweepStatusPaused  SweepStatus = "paused"
	SweepStatusResumed SweepStatus = "resumed"
	SweepStatusUnknown SweepStatus = "unknown"
)

// GetSweep returns the current start and stop frequencies and sweep points as a Sweep struct.
func (d *Device) GetSweep() (Sweep, error) {
	d.logger.Info("requesting sweep")

	line, err := d.sendCommand("sweep")
	if err != nil {
		return Sweep{}, err
	}
	
	sweep, err := parseSweepResponse(line)
	if err != nil {
		d.logger.Error("failed to parse sweep response", "line", line, "err", err)
		return Sweep{}, fmt.Errorf("%w: failed to parse sweep response: %v", ErrCommandFailed, err)
	}

	return sweep, nil
}

// GetSweepStatus returns the current sweep status as SweepStatus type.
func (d *Device) GetSweepStatus() (SweepStatus, error) {
	d.logger.Debug("requesting sweep status")

	res, err := d.sendCommand("status")
	if err != nil {
		return SweepStatusUnknown, err
	}

	switch res {
	case "Paused":
		return SweepStatusPaused, nil
	case "Resumed":
		return SweepStatusResumed, nil
	default:
		d.logger.Error("unexpected sweep status", "sweep_status", res)
		return SweepStatusUnknown, fmt.Errorf("%w: unexpected sweep status %s", ErrCommandFailed, res)
	}
}

// SetSweepMode sets the sweep mode.
func (d *Device) SetSweepMode(mode SweepMode) error {
	d.logger.Info("setting sweep mode", "mode", mode)
	_, err := d.sendCommand(fmt.Sprintf("sweep %s", mode))
	return err
}

// SetSweepStart sets the sweep start frequency in Hz.
func (d *Device) SetSweepStart(freqHz uint64) error {
	d.logger.Info("setting sweep start", "freq", freqHz)
	_, err := d.sendCommand(fmt.Sprintf("sweep start %d", freqHz))
	return err
}

// SetSweepStop sets the sweep stop frequency in Hz.
func (d *Device) SetSweepStop(freqHz uint64) error {
	d.logger.Info("setting sweep stop", "freq", freqHz)
	_, err := d.sendCommand(fmt.Sprintf("sweep stop %d", freqHz))
	return err
}

// SetSweepCenter sets the sweep center frequency in Hz.
func (d *Device) SetSweepCenter(freqHz uint64) error {
	d.logger.Info("setting sweep center", "freq", freqHz)
	_, err := d.sendCommand(fmt.Sprintf("sweep center %d", freqHz))
	return err
}

// SetSweepSpan sets the sweep span frequency in Hz.
func (d *Device) SetSweepSpan(freqHz uint64) error {
	d.logger.Info("setting sweep span", "freq", freqHz)
	_, err := d.sendCommand(fmt.Sprintf("sweep span %d", freqHz))
	return err
}

// SetSweepContinuousWave sets the sweep to continuous wave mode at the specified frequency in Hz.
func (d *Device) SetSweepContinuousWave(freqHz uint64) error {
	d.logger.Info("setting sweep continuous wave", "freq", freqHz)
	_, err := d.sendCommand(fmt.Sprintf("sweep cw %d", freqHz))
	return err
}

// SetSweepStartStop sets the sweep start and stop frequency in Hz.
func (d *Device) SetSweepStartStop(freqStartHz uint64, freqStopHz uint64) error {
	d.logger.Info("set sweep start and stop", "freq_start", freqStartHz, "freq_stop", freqStopHz)
	_, err := d.sendCommand(fmt.Sprintf("sweep %d %d", freqStartHz, freqStopHz))
	return err
}

// SetSweepStartStopWithPoints sets the sweep start and stop frequencies in Hz and the number of sweep points.
func (d *Device) SetSweepStartStopWithPoints(freqStartHz uint64, freqStopHz uint64, points uint) error {
	d.logger.Info("set sweep start, stop and points", "freq_start", freqStartHz, "freq_stop", freqStopHz, "points", points)
	_, err := d.sendCommand(fmt.Sprintf("sweep %d %d %d", freqStartHz, freqStopHz, points))
	return err
}

// SetSweepTime sets the sweep time in microseconds.
// TODO: use time type
func (d *Device) SetSweepTime(timeUs uint64) error {
	d.logger.Info("setting sweep time", "time", timeUs)
	_, err := d.sendCommand(fmt.Sprintf("sweeptime %du", timeUs))
	return err
}

// SetSweepPoints updates the number of sweep points while preserving the current start and stop frequencies.
func (d *Device) SetSweepPoints(points uint) error {
	// Sweep points can't be set directly; retrieve current sweep settings
	// and apply them together with the new point count.
	sweep, err := d.GetSweep()
	if err != nil {
		return err
	}
	return d.SetSweepStartStopWithPoints(sweep.Start, sweep.Stop, points)
}

// PauseSweep pauses the ongoing sweep operation.
func (d *Device) PauseSweep() error {
	d.logger.Info("pausing sweep")
	_, err := d.sendCommand("pause")
	return err
}

// ResumeSweep resumes the paused sweep operation.
func (d *Device) ResumeSweep() error {
	d.logger.Info("resuming sweep")
	_, err := d.sendCommand("resume")
	return err
}
