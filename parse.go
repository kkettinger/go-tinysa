package tinysa

import (
	"fmt"
	"strconv"
	"strings"
)

// parseBatteryResponse parses a battery response into a uint voltage (mV).
//
// Example response: `4191 mV`
func parseBatteryResponse(response string) (uint, error) {
	parts := strings.Split(response, " ")
	if len(parts) != 2 {
		return 0, fmt.Errorf("expected 2 fields, got %d", len(parts))
	}

	vbat, err := strconv.ParseInt(parts[0], 10, 16)
	if err != nil {
		return 0, fmt.Errorf("integer conversion failed: %s", err.Error())
	}

	if vbat < 0 {
		return 0, fmt.Errorf("negative voltage")
	}

	return uint(vbat), nil
}

// parseMarkerResponseLine parses a single line of a marker response into a Marker struct.
//
// Example response: `1 216 522167037 -9.08e+01`
func parseMarkerResponseLine(line string) (Marker, error) {
	fields := strings.Fields(line)
	if len(fields) != 4 {
		return Marker{}, fmt.Errorf("expected 4 fields, got %d", len(fields))
	}

	marker, err := strconv.ParseUint(fields[0], 10, 0)
	if err != nil {
		return Marker{}, fmt.Errorf("invalid marker %q: %s", fields[0], err.Error())
	}

	index, err := strconv.ParseUint(fields[1], 10, 0)
	if err != nil {
		return Marker{}, fmt.Errorf("invalid index %q: %s", fields[1], err.Error())
	}

	freq, err := strconv.ParseUint(fields[2], 10, 64)
	if err != nil {
		return Marker{}, fmt.Errorf("invalid frequency %q: %s", fields[2], err.Error())
	}

	val, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return Marker{}, fmt.Errorf("invalid value %q: %s", fields[3], err.Error())
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
func parseSweepResponse(response string) (Sweep, error) {
	parts := strings.Split(response, " ")
	if len(parts) != 3 {
		return Sweep{}, fmt.Errorf("expected 3 fields, got %d", len(parts))
	}

	sweepStart, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return Sweep{}, fmt.Errorf("integer conversion failed: %s", err.Error())
	}

	sweepStop, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return Sweep{}, fmt.Errorf("integer conversion failed: %s", err.Error())
	}

	sweepPoints, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return Sweep{}, fmt.Errorf("integer conversion failed: %s", err.Error())
	}

	return Sweep{
		Start:  uint64(sweepStart),
		Stop:   uint64(sweepStop),
		Points: uint(sweepPoints),
	}, nil
}

// parseTraceValueResponseLine parses a single line of a trace value response into a TraceValue struct.
//
// Example response: `trace 1 value 442 -108.88`
func parseTraceValueResponseLine(line string) (TraceValue, error) {
	fields := strings.Fields(line)
	if len(fields) != 5 {
		return TraceValue{}, fmt.Errorf("expected 5 fields, got %d", len(fields))
	}

	trace, err := strconv.ParseUint(fields[1], 10, 32)
	if err != nil {
		return TraceValue{}, fmt.Errorf("invalid trace: %s", err.Error())
	}

	point, err := strconv.ParseUint(fields[3], 10, 32)
	if err != nil {
		return TraceValue{}, fmt.Errorf("invalid point: %s", err.Error())
	}

	value, err := strconv.ParseFloat(fields[4], 64)
	if err != nil {
		return TraceValue{}, fmt.Errorf("invalid value: %s", err.Error())
	}

	return TraceValue{
		Trace: uint(trace),
		Point: uint(point),
		Value: value,
	}, nil
}

// parseTraceResponseLine parses a single line of a trace response into a Trace struct.
//
// Example response: `1: dBm -30.000000000 10.000000000`
func parseTraceResponseLine(line string) (Trace, error) {
	fields := strings.Fields(line)
	if len(fields) != 4 {
		return Trace{}, fmt.Errorf("expected 4 fields, got %d", len(fields))
	}

	trace, err := strconv.Atoi(strings.TrimSuffix(fields[0], ":"))
	if err != nil {
		return Trace{}, fmt.Errorf("failed to parse trace id: %s", err.Error())
	}

	unit, ok := TraceUnitFromString(fields[1])
	if !ok {
		return Trace{}, fmt.Errorf("failed to parse unit: %s", fields[1])
	}

	refPos, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return Trace{}, fmt.Errorf("failed to parse ref position: %s", err.Error())
	}

	scale, err := strconv.ParseFloat(fields[3], 64)
	if err != nil {
		return Trace{}, fmt.Errorf("failed to parse ref scale: %s", err.Error())
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
