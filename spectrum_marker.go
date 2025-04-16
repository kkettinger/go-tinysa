package tinysa

import (
	"fmt"
	"strconv"
	"strings"
)

type Marker struct {
	Marker    uint
	Index     uint
	Frequency uint64
	Value     float64
}

// GetMarker returns marker information for the given marker ID as Marker struct.
func (d *Device) GetMarker(markerId uint) (Marker, error) {
	d.logger.Info("requesting marker information", "marker_id", markerId)

	line, err := d.sendCommand(fmt.Sprintf("marker %d", markerId))
	if err != nil {
		return Marker{}, err
	}

	result, err := parseMarkerResultLine(line)
	if err != nil {
		d.logger.Error("failed to parse marker result", "line", line, "err", err)
		return Marker{}, fmt.Errorf("%w: failed to parse marker result: %v", ErrCommandFailed, err)
	}

	return result, nil
}

// GetMarkerAll requests all marker information and returns a Marker slice containing index, frequency, and power.
func (d *Device) GetMarkerAll() ([]Marker, error) {
	d.logger.Info("requesting all marker information")

	statusStr, err := d.sendCommand("marker")
	if err != nil {
		return nil, err
	}

	var status []Marker
	lines := strings.Split(statusStr, commandTerminator)
	for _, line := range lines {
		if s, err := parseMarkerResultLine(line); err != nil {
			d.logger.Error("failed to parse marker result", "line", line, "err", err)
			return nil, fmt.Errorf("%w: failed to parse marker result: %v", ErrCommandFailed, err)
		} else {
			status = append(status, s)
		}
	}

	return status, nil
}

// EnableMarker enables the marker for the specified markerId.
func (d *Device) EnableMarker(markerId uint) error {
	d.logger.Info("enabling marker", "marker_id", markerId)
	_, err := d.sendCommand(fmt.Sprintf("marker %d on", markerId))
	return err
}

// DisableMarker disables the marker for the specified markerId.
func (d *Device) DisableMarker(markerId uint) error {
	d.logger.Info("disabling marker", "marker_id", markerId)
	_, err := d.sendCommand(fmt.Sprintf("marker %d off", markerId))
	return err
}

// SetMarkerFreq sets the marker to the specified frequency.
func (d *Device) SetMarkerFreq(markerId uint, freqHz uint64) error {
	d.logger.Info("setting marker frequency", "marker_id", markerId, "freq", freqHz)
	_, err := d.sendCommand(fmt.Sprintf("marker %d %d", markerId, freqHz))
	return err
}

// SetMarkerTrace assigns the specified marker to the specified trace.
func (d *Device) SetMarkerTrace(markerId uint, traceId uint) error {
	d.logger.Info("assigning marker to trace", "marker_id", markerId, "trace_id", traceId)
	_, err := d.sendCommand(fmt.Sprintf("marker %d trace %d", markerId, traceId))
	return err
}

// MoveMarkerPeak moves the marker to the peak value of the assigned trace.
func (d *Device) MoveMarkerPeak(markerId uint) error {
	d.logger.Info("move marker peak", "marker_id", markerId)
	_, err := d.sendCommand(fmt.Sprintf("marker %d peak", markerId))
	return err
}

// EnableMarkerDelta sets the specified marker to delta mode, referencing the specified marker.
func (d *Device) EnableMarkerDelta(markerId uint, refMarkerId uint) error {
	d.logger.Info("enabling marker delta", "marker_id", markerId, "ref_marker_id", refMarkerId)
	_, err := d.sendCommand(fmt.Sprintf("marker %d delta %d", markerId, refMarkerId))
	return err
}

// DisableMarkerDelta disables delta mode for the specified marker.
func (d *Device) DisableMarkerDelta(markerId uint) error {
	d.logger.Info("disabling marker delta", "marker_id", markerId)
	_, err := d.sendCommand(fmt.Sprintf("marker %d delta off", markerId))
	return err
}

// parseMarkerResultLine parses a marker response line like `1 216 522167037 -9.08e+01` into a Marker struct.
func parseMarkerResultLine(line string) (Marker, error) {
	fields := strings.Fields(line)
	if len(fields) != 4 {
		return Marker{}, fmt.Errorf("%w: expected 4 fields, got %d", ErrCommandFailed, len(fields))
	}

	marker, err := strconv.ParseUint(fields[0], 10, 0)
	if err != nil {
		return Marker{}, fmt.Errorf("%w: invalid marker %q: %b", ErrCommandFailed, fields[0], err)
	}

	index, err := strconv.ParseUint(fields[1], 10, 0)
	if err != nil {
		return Marker{}, fmt.Errorf("%w: ErrCommandFailedinvalid index %q: %b", ErrCommandFailed, fields[1], err)
	}

	freq, err := strconv.ParseUint(fields[2], 10, 64)
	if err != nil {
		return Marker{}, fmt.Errorf("%w: invalid frequency %q: %b", ErrCommandFailed, fields[2], err)
	}

	val, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return Marker{}, fmt.Errorf("%w: invalid value %q: %v", ErrCommandFailed, fields[3], err)
	}

	return Marker{
		Marker:    uint(marker),
		Index:     uint(index),
		Frequency: freq,
		Value:     val,
	}, nil
}
