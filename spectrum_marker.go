package tinysa

import (
	"fmt"
	"strings"
)

// Marker represents a marker with its marker id, position index, frequency in Hz, and associated measured value.
type Marker struct {
	Marker    uint    // Marker identifier
	Index     uint    // Position index
	Frequency uint64  // Frequency in Hz
	Value     float64 // Measured value at the given frequency
}

// GetMarker returns marker information for the given marker ID as Marker struct.
func (d *Device) GetMarker(markerID uint) (Marker, error) {
	d.logger.Info("requesting marker information", "marker_id", markerID)

	line, err := d.sendCommand(fmt.Sprintf("marker %d", markerID))
	if err != nil {
		return Marker{}, err
	}

	result, err := parseMarkerResponseLine(line)
	if err != nil {
		d.logger.Error("failed to parse marker result", "line", line, "err", err)
		return Marker{}, fmt.Errorf("failed to parse marker result: %s", err.Error())
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
		s, err := parseMarkerResponseLine(line)
		if err != nil {
			d.logger.Error("failed to parse marker result", "line", line, "err", err)
			return nil, fmt.Errorf("failed to parse marker result: %s", err.Error())
		}
		status = append(status, s)
	}

	return status, nil
}

// EnableMarker enables the marker for the specified markerId.
func (d *Device) EnableMarker(markerID uint) error {
	d.logger.Info("enabling marker", "marker_id", markerID)
	_, err := d.sendCommand(fmt.Sprintf("marker %d on", markerID))
	return err
}

// DisableMarker disables the marker for the specified markerId.
func (d *Device) DisableMarker(markerID uint) error {
	d.logger.Info("disabling marker", "marker_id", markerID)
	_, err := d.sendCommand(fmt.Sprintf("marker %d off", markerID))
	return err
}

// SetMarkerFreq sets the marker to the specified frequency.
func (d *Device) SetMarkerFreq(markerID uint, freqHz uint64) error {
	d.logger.Info("setting marker frequency", "marker_id", markerID, "freq", freqHz)
	_, err := d.sendCommand(fmt.Sprintf("marker %d %d", markerID, freqHz))
	return err
}

// SetMarkerTrace assigns the specified marker to the specified trace.
func (d *Device) SetMarkerTrace(markerID uint, traceID uint) error {
	d.logger.Info("assigning marker to trace", "marker_id", markerID, "trace_id", traceID)
	_, err := d.sendCommand(fmt.Sprintf("marker %d trace %d", markerID, traceID))
	return err
}

// MoveMarkerPeak moves the marker to the peak value of the assigned trace.
func (d *Device) MoveMarkerPeak(markerID uint) error {
	d.logger.Info("move marker peak", "marker_id", markerID)
	_, err := d.sendCommand(fmt.Sprintf("marker %d peak", markerID))
	return err
}

// EnableMarkerDelta sets the specified marker to delta mode, referencing the specified marker.
func (d *Device) EnableMarkerDelta(markerID uint, refMarkerID uint) error {
	d.logger.Info("enabling marker delta", "marker_id", markerID, "ref_marker_id", refMarkerID)
	_, err := d.sendCommand(fmt.Sprintf("marker %d delta %d", markerID, refMarkerID))
	return err
}

// DisableMarkerDelta disables delta mode for the specified marker.
func (d *Device) DisableMarkerDelta(markerID uint) error {
	d.logger.Info("disabling marker delta", "marker_id", markerID)
	_, err := d.sendCommand(fmt.Sprintf("marker %d delta off", markerID))
	return err
}

// EnableMarkerTracking enables tracking of the peak value for the assigned trace of the given marker.
func (d *Device) EnableMarkerTracking(markerID uint) error {
	d.logger.Info("enabling marker tracking", "marker_id", markerID)
	_, err := d.sendCommand(fmt.Sprintf("marker %d tracking on", markerID))
	return err
}

// DisableMarkerTracking disables tracking of the peak value for the assigned trace of the given marker.
func (d *Device) DisableMarkerTracking(markerID uint) error {
	d.logger.Info("disabling marker tracking", "marker_id", markerID)
	_, err := d.sendCommand(fmt.Sprintf("marker %d tracking off", markerID))
	return err
}
