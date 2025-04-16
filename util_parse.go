package tinysa

import "strings"

// SweepModeFromString parses a string into a SweepMode (case-insensitive).
func SweepModeFromString(s string) (SweepMode, bool) {
	for k, v := range sweepModeMap {
		if strings.EqualFold(k, s) {
			return v, true
		}
	}
	return SweepMode{}, false
}

// TraceUnitFromString parses a string into a DisplayUnit (case-insensitive).
func TraceUnitFromString(s string) (DisplayUnit, bool) {
	for k, v := range displayUnitMap {
		if strings.EqualFold(k, s) {
			return v, true
		}
	}
	return DisplayUnit{}, false
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
