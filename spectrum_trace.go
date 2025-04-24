package tinysa

import (
	"fmt"
	"strconv"
	"strings"
)

// TraceUnit represents the unit of measurement for a trace value.
type TraceUnit struct {
	value string
}

// String returns the string representation of the TraceUnit.
func (u TraceUnit) String() string {
	return u.value
}

// IsValid reports whether the TraceUnit contains a valid unit.
func (u TraceUnit) IsValid() bool {
	return u.value != ""
}

const (
	traceUnitRAW  string = "RAW"
	traceUnitDBm  string = "dBm"
	traceUnitDBmV string = "dBmV"
	traceUnitDBuV string = "dBuV"
	traceUnitV    string = "V"
	traceUnitVpp  string = "Vpp"
	traceUnitW    string = "W"
)

var (
	// TraceUnitRaw represents raw, unprocessed trace values.
	TraceUnitRaw = TraceUnit{traceUnitRAW}

	// TraceUnitDBm represents values in decibels relative to 1 milliwatt (dBm).
	TraceUnitDBm = TraceUnit{traceUnitDBm}

	// TraceUnitDBmV represents values in decibels relative to 1 millivolt (dBmV).
	TraceUnitDBmV = TraceUnit{traceUnitDBmV}

	// TraceUnitDBuV represents values in decibels relative to 1 microvolt (dBuV).
	TraceUnitDBuV = TraceUnit{traceUnitDBuV}

	// TraceUnitV represents values in volts (V).
	TraceUnitV = TraceUnit{traceUnitV}

	// TraceUnitVpp represents values in volts peak-to-peak (Vpp).
	TraceUnitVpp = TraceUnit{traceUnitVpp}

	// TraceUnitW represents values in watts (W).
	TraceUnitW = TraceUnit{traceUnitW}
)

var traceUnitMap = map[string]TraceUnit{
	traceUnitRAW:  TraceUnitRaw,
	traceUnitDBm:  TraceUnitDBm,
	traceUnitDBmV: TraceUnitDBmV,
	traceUnitDBuV: TraceUnitDBuV,
	traceUnitV:    TraceUnitV,
	traceUnitVpp:  TraceUnitVpp,
	traceUnitW:    TraceUnitW,
}

var traceUnitOptions = []string{
	traceUnitRAW,
	traceUnitDBm,
	traceUnitDBmV,
	traceUnitDBuV,
	traceUnitV,
	traceUnitVpp,
	traceUnitW,
}

// TraceUnitOptions returns a list of possible trace units like "dBm" or "Vpp".
func TraceUnitOptions() []string {
	return traceUnitOptions
}

// Trace represents a single trace status containing the trace id, trace unit, reference position and scale value.
type Trace struct {
	Trace  uint
	Unit   TraceUnit
	RefPos float64
	Scale  float64
}

// TraceValue represents a single trace data point containing the trace id, index point and signal value in dBm.
type TraceValue struct {
	Trace uint
	Point uint
	Value float64
}

// TraceData represents a single trace data point containing the trace id, index point, frequency value in Hz, and
// signal value in dBm.
type TraceData struct {
	Trace     uint
	Point     uint
	Frequency uint64
	Value     float64
}

// TraceCalc contains the string for the trace calculation option like "maxh" (see traceCalc*).
type TraceCalc struct {
	value string
}

func (c TraceCalc) String() string {
	return c.value
}

// IsValid returns true if a valid TraceCalc option is set.
func (c TraceCalc) IsValid() bool {
	return c.value != ""
}

const (
	traceCalcMinH   string = "minh"
	traceCalcMaxH   string = "maxh"
	traceCalcMaxD   string = "maxd"
	traceCalcAver4  string = "aver4"
	traceCalcAver16 string = "aver16"
	traceCalcQuasi  string = "quasi"
	// traceCalcLog    string = "log"
	// traceCalcLin    string = "lin"
)

var (
	// TraceCalcMinH sets the trace to hold the minimum value measured.
	TraceCalcMinH = TraceCalc{traceCalcMinH}

	// TraceCalcMaxH sets the trace to hold the maximum value measured.
	TraceCalcMaxH = TraceCalc{traceCalcMaxH}

	// TraceCalcMaxD holds the maximum value for a limited number of scans.
	TraceCalcMaxD = TraceCalc{traceCalcMaxD}

	// TraceCalcAver4 enables running average over 4 samples.
	TraceCalcAver4 = TraceCalc{traceCalcAver4}

	// TraceCalcAver16 enables running average over 16 samples.
	TraceCalcAver16 = TraceCalc{traceCalcAver16}

	// TraceCalcQuasi sets quasi-peak hold mode.
	TraceCalcQuasi = TraceCalc{traceCalcQuasi}

	// TODO: TraceCalcLog and TraceCalcLin are not currently supported (require extra arguments).
	// TraceCalcLog    = TraceCalc{traceCalcLog}
	// TraceCalcLin    = TraceCalc{traceCalcLin}
)

var traceCalcMap = map[string]TraceCalc{
	traceCalcMinH:   TraceCalcMinH,
	traceCalcMaxH:   TraceCalcMaxH,
	traceCalcMaxD:   TraceCalcMaxD,
	traceCalcAver4:  TraceCalcAver4,
	traceCalcAver16: TraceCalcAver16,
	traceCalcQuasi:  TraceCalcQuasi,
	//traceCalcLog:    TraceCalcLog,
	//traceCalcLin:    TraceCalcLin,
}

var traceCalcOptions = []string{
	traceCalcMinH,
	traceCalcMaxH,
	traceCalcMaxD,
	traceCalcAver4,
	traceCalcAver16,
	traceCalcQuasi,
	// traceCalcLog,
	// traceCalcLin,
}

// TraceCalcOptions returns a list of possible trace calculations options like "minh" or "quasi".
func TraceCalcOptions() []string {
	return traceCalcOptions
}

// GetTrace returns trace information for the specified trace id as Trace struct.
func (d *Device) GetTrace(traceID uint) (Trace, error) {
	d.logger.Info("requesting trace information", "trace_id", traceID)

	line, err := d.sendCommand(fmt.Sprintf("trace %d", traceID))
	if err != nil {
		return Trace{}, err
	}

	result, err := parseTraceResponseLine(line)
	if err != nil {
		d.logger.Error("failed to parse trace result", "trace_id", traceID, "line", line, "err", err)
		return Trace{}, fmt.Errorf("failed to parse trace result: %s", err.Error())
	}

	return result, nil
}

// GetTraceAll returns trace information for all traces as Trace slice.
func (d *Device) GetTraceAll() ([]Trace, error) {
	d.logger.Info("requesting trace information")

	statusStr, err := d.sendCommand("trace")
	if err != nil {
		return nil, err
	}

	var status []Trace
	lines := strings.Split(statusStr, commandTerminator)
	for _, line := range lines {
		s, err := parseTraceResponseLine(line)
		if err != nil {
			d.logger.Error("failed to parse trace result", "line", line, "err", err)
			return nil, fmt.Errorf("failed to parse trace result: %s", err.Error())
		}

		status = append(status, s)
	}

	return status, nil
}

// GetTraceFrequencies returns a list of frequencies as uint64 list for the current (or really all) traces.
func (d *Device) GetTraceFrequencies() ([]uint64, error) {
	d.logger.Info("getting trace frequencies")

	freqStr, err := d.sendCommand("frequencies")
	if err != nil {
		return nil, err
	}
	freqList := strings.Split(freqStr, commandTerminator)

	result := make([]uint64, len(freqList))
	for i, freq := range freqList {
		freqInt, err := strconv.ParseUint(freq, 10, 64)
		if err != nil {
			d.logger.Error("failed to parse trace frequency", "freq", freq, "err", err)
			return nil, fmt.Errorf("failed to parse frequency %q as int: %s", freq, err.Error())
		}

		result[i] = freqInt
	}

	return result, nil
}

// GetTraceValues returns the list of the trace values as a TraceValue slice.
func (d *Device) GetTraceValues(traceID uint) ([]TraceValue, error) {
	d.logger.Info("getting trace values", "trace_id", traceID)

	dataStr, err := d.sendCommand(fmt.Sprintf("trace %d value", traceID))
	if err != nil {
		return nil, err
	}
	dataList := strings.Split(dataStr, commandTerminator)

	data := make([]TraceValue, len(dataList))
	for i, line := range dataList {
		data[i], err = parseTraceValueResponseLine(line)
		if err != nil {
			d.logger.Error("failed to parse trace value", "trace_id", traceID, "line", line)
			return nil, fmt.Errorf("failed to parse trace result: %s", err.Error())
		}
	}

	return data, nil
}

// GetTraceData returns a combined list of frequencies and values as a TraceData slice.
func (d *Device) GetTraceData(traceID uint) ([]TraceData, error) {
	d.logger.Info("getting trace data", "trace_id", traceID)

	values, err := d.GetTraceValues(traceID)
	if err != nil {
		return nil, err
	}

	frequencies, err := d.GetTraceFrequencies()
	if err != nil {
		return nil, err
	}

	lenValues := len(values)
	lenFreq := len(frequencies)
	if lenValues != lenFreq {
		d.logger.Error("value and frequency values lengths do not match", "trace_id", traceID, "values", values, "frequencies", frequencies)
		return nil, fmt.Errorf("value and frequency values lengths do not match, %d != %d", lenValues, lenFreq)
	}

	data := make([]TraceData, lenValues)
	for i := range values {
		data[i] = TraceData{
			Trace:     values[i].Trace,
			Point:     values[i].Point,
			Value:     values[i].Value,
			Frequency: frequencies[i],
		}
	}

	return data, nil
}

// EnableTrace enables the display of the specified trace.
func (d *Device) EnableTrace(traceID uint) error {
	d.logger.Info("enabling trace", "trace_id", traceID)
	_, err := d.sendCommand(fmt.Sprintf("trace %d view on", traceID))
	return err
}

// DisableTrace disables the display of the specified trace.
func (d *Device) DisableTrace(traceID uint) error {
	_, err := d.sendCommand(fmt.Sprintf("trace %d view off", traceID))
	return err
}

// EnableTraceCalc enables trace calculations like TraceCalcMaxH or TraceCalcQuasi for the specified trace.
func (d *Device) EnableTraceCalc(traceID uint, calc TraceCalc) error {
	d.logger.Info("enabling trace calculations", "trace_id", traceID, "calc", calc)

	// calc log and lin is only supported on ultra
	/*if d.model != ModelUltra {
		if calc == TraceCalcLin || calc == TraceCalcLog {
			return ErrOptionNotSupportedByModel
		}
	}*/
	_, err := d.sendCommand(fmt.Sprintf("calc %d %s", traceID, calc.String()))
	return err
}

// DisableTraceCalc disables calculation for the specified trace.
func (d *Device) DisableTraceCalc(traceID uint) error {
	d.logger.Info("disabling trace calculations", "trace_id", traceID)
	_, err := d.sendCommand(fmt.Sprintf("calc %d off", traceID))
	return err
}

// SetTraceUnit sets the display unit to the specified value.
func (d *Device) SetTraceUnit(unit TraceUnit) error {
	d.logger.Info("setting display unit", "unit", unit)
	_, err := d.sendCommand(fmt.Sprintf("trace %s", unit.value))
	return err
}

// SetTraceRefLevel sets the display ref level to the specified value in dBm.
func (d *Device) SetTraceRefLevel(levelDbm int) error {
	d.logger.Info("setting trace ref level", "level", levelDbm)
	_, err := d.sendCommand(fmt.Sprintf("trace reflevel %d", levelDbm))
	return err
}

// SetTraceRefLevelAuto sets the display ref level to auto.
func (d *Device) SetTraceRefLevelAuto() error {
	d.logger.Info("setting trace ref level auto")
	_, err := d.sendCommand("trace reflevel auto")
	return err
}

// SetTraceScale sets the display scale to the specified value.
func (d *Device) SetTraceScale(level float64) error {
	d.logger.Info("setting trace scale", "level", level)
	_, err := d.sendCommand(fmt.Sprintf("trace scale %.3f", level))
	return err
}
