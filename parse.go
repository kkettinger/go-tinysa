package tinysa

import (
	"fmt"
	"strconv"
	"strings"
)

// parseBatteryVoltageLine parses a battery response into a uint.
//
// Example response: `4191 mV`
func parseBatteryVoltageLine(line string) (uint, error) {
	parts := strings.Split(line, " ")
	if len(parts) != 2 {
		return 0, fmt.Errorf("%w: expected 2 fields, got %d", ErrCommandFailed, len(parts))
	}

	vbat, err := strconv.ParseInt(parts[0], 10, 16)
	if err != nil {
		return 0, fmt.Errorf("%w: integer conversion failed: %v", ErrCommandFailed, err)
	}

	if vbat < 0 {
		return 0, fmt.Errorf("%w: negative voltage", ErrCommandFailed)
	}

	return uint(vbat), nil
}

// parseMarkerResultLine parses a marker response into a Marker struct.
//
// Example response: `1 216 522167037 -9.08e+01` (single line)
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

// parseSweepResponse parses a sweep response into a Sweep struct.
//
// Example response: `450000000 600000000 450`
func parseSweepResponse(line string) (Sweep, error) {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return Sweep{}, fmt.Errorf("%w: expected 3 fields, got %d", ErrCommandFailed, len(parts))
	}

	sweepStart, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return Sweep{}, fmt.Errorf("%w: integer conversion failed: %v", ErrCommandFailed, err)
	}

	sweepStop, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return Sweep{}, fmt.Errorf("%w: integer conversion failed: %v", ErrCommandFailed, err)
	}

	sweepPoints, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return Sweep{}, fmt.Errorf("%w: integer conversion failed: %v", ErrCommandFailed, err)
	}

	return Sweep{
		Start:  uint64(sweepStart),
		Stop:   uint64(sweepStop),
		Points: uint(sweepPoints),
	}, nil
}

// parseTraceValue parses a trace response into a TraceValue struct.
//
// Example response: `trace 1 value 442 -108.88` (single line)
func parseTraceValue(line string) (TraceValue, error) {
	fields := strings.Fields(line)
	if len(fields) != 5 {
		return TraceValue{}, fmt.Errorf("%w: expected 5 fields, got %d", ErrCommandFailed, len(fields))
	}

	trace, err := strconv.ParseUint(fields[1], 10, 32)
	if err != nil {
		return TraceValue{}, fmt.Errorf("%w: invalid trace: %v", ErrCommandFailed, err)
	}

	point, err := strconv.ParseUint(fields[3], 10, 32)
	if err != nil {
		return TraceValue{}, fmt.Errorf("%w: invalid point: %v", ErrCommandFailed, err)
	}

	value, err := strconv.ParseFloat(fields[4], 64)
	if err != nil {
		return TraceValue{}, fmt.Errorf("%w: invalid value: %v", ErrCommandFailed, err)
	}

	return TraceValue{
		Trace: uint(trace),
		Point: uint(point),
		Value: value,
	}, nil
}

// parseTraceStatusLine parses a trace response into a Trace struct.
//
// Example response: `1: dBm -30.000000000 10.000000000` (single line)
func parseTraceStatusLine(line string) (Trace, error) {
	fields := strings.Fields(line)
	if len(fields) != 4 {
		return Trace{}, fmt.Errorf("%w: expected 4 fields, got %d", ErrCommandFailed, len(fields))
	}

	trace, err := strconv.Atoi(strings.TrimSuffix(fields[0], ":"))
	if err != nil {
		return Trace{}, fmt.Errorf("%w: failed to parse trace id: %v", ErrCommandFailed, err)
	}

	unit, ok := TraceUnitFromString(fields[1])
	if !ok {
		return Trace{}, fmt.Errorf("%w: failed to parse unit: %v", ErrCommandFailed, fields[1])
	}

	refPos, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return Trace{}, fmt.Errorf("%w: failed to parse ref position: %v", ErrCommandFailed, err)
	}

	scale, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return Trace{}, fmt.Errorf("%w: failed to parse ref scale: %v", ErrCommandFailed, err)
	}

	return Trace{
		Trace:  trace,
		Unit:   unit,
		RefPos: refPos,
		Scale:  scale,
	}, nil
}

// SweepModeFromString parses a string into a SweepMode (case-insensitive).
func SweepModeFromString(s string) (SweepMode, bool) {
	for k, v := range sweepModeMap {
		if strings.EqualFold(k, s) {
			return v, true
		}
	}
	return SweepMode{}, false
}

// TraceUnitFromString parses a string into a TraceUnit (case-insensitive).
func TraceUnitFromString(s string) (TraceUnit, bool) {
	for k, v := range traceUnitMap {
		if strings.EqualFold(k, s) {
			return v, true
		}
	}
	return TraceUnit{}, false
}

// TraceCalcFromString parses a string into a TraceCalc (case-insensitive).
func TraceCalcFromString(s string) (TraceCalc, bool) {
	for k, v := range traceCalcMap {
		if strings.EqualFold(k, s) {
			return v, true
		}
	}
	return TraceCalc{}, false
}
