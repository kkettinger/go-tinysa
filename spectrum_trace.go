package tinysa

import (
	"fmt"
	"strconv"
	"strings"
)

type Trace struct {
	Trace  int // TODO: uint
	Unit   DisplayUnit
	RefPos float64
	Scale  float64
}

type TraceValue struct {
	Trace uint
	Point uint
	Value float64
}

type TraceData struct {
	Trace     uint
	Point     uint
	Frequency uint64
	Value     float64
}

type TraceCalc struct {
	value string
}

func (c TraceCalc) String() string {
	return c.value
}

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
	//traceCalcLog    string = "log"
	// traceCalcLin    string = "lin"
)

var (
	TraceCalcMinH   = TraceCalc{traceCalcMinH}
	TraceCalcMaxH   = TraceCalc{traceCalcMaxH}
	TraceCalcMaxD   = TraceCalc{traceCalcMaxD}
	TraceCalcAver4  = TraceCalc{traceCalcAver4}
	TraceCalcAver16 = TraceCalc{traceCalcAver16}
	TraceCalcQuasi  = TraceCalc{traceCalcQuasi}
	//TraceCalcLog    = TraceCalc{traceCalcLog} // TODO: Currently not supported (needs extra argument)
	// TraceCalcLin    = TraceCalc{traceCalcLin} // TODO: Currently not supported (needs extra argument)
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
	//traceCalcLog,
	//traceCalcLin,
}

func TraceCalcOptions() []string {
	return traceCalcOptions
}

// GetTrace returns trace information for the specified trace id as Trace struct.
func (d *Device) GetTrace(traceId uint) (Trace, error) {
	d.logger.Info("requesting trace information", "trace_id", traceId)

	line, err := d.sendCommand(fmt.Sprintf("trace %d", traceId))
	if err != nil {
		return Trace{}, err
	}

	result, err := parseTraceStatusLine(line)
	if err != nil {
		d.logger.Error("failed to parse trace result", "trace_id", traceId, "line", line, "err", err)
		return Trace{}, fmt.Errorf("%w: failed to parse trace result: %v", ErrCommandFailed, err)
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
		s, err := parseTraceStatusLine(line)
		if err != nil {
			d.logger.Error("failed to parse trace result", "line", line, "err", err)
			return nil, fmt.Errorf("%w: failed to parse trace result: %v", ErrCommandFailed, err)
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
			return nil, fmt.Errorf("%w: failed to parse frequency %q as int: %v", ErrCommandFailed, freq, err)
		}

		result[i] = freqInt
	}

	return result, nil
}

// GetTraceValues returns the list of the trace values as a TraceValue slice.
func (d *Device) GetTraceValues(traceId uint) ([]TraceValue, error) {
	d.logger.Info("getting trace values", "trace_id", traceId)

	dataStr, err := d.sendCommand(fmt.Sprintf("trace %d value", traceId))
	if err != nil {
		return nil, err
	}
	dataList := strings.Split(dataStr, commandTerminator)

	data := make([]TraceValue, len(dataList))
	for i, line := range dataList {
		data[i], err = parseTraceValue(line)
		if err != nil {
			d.logger.Error("failed to parse trace value", "trace_id", traceId, "line", line)
			return nil, fmt.Errorf("%w: failed to parse trace result: %v", ErrCommandFailed, err)
		}
	}

	return data, nil
}

// GetTraceData returns a combined list of frequencies and values as a TraceData slice.
func (d *Device) GetTraceData(traceId uint) ([]TraceData, error) {
	d.logger.Info("getting trace data", "trace_id", traceId)

	values, err := d.GetTraceValues(traceId)
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
		d.logger.Error("value and frequency values lengths do not match", "trace_id", traceId, "values", values, "frequencies", frequencies)
		return nil, fmt.Errorf("%w: value and frequency values lengths do not match, %d != %d", ErrCommandFailed, lenValues, lenFreq)
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
func (d *Device) EnableTrace(traceId uint) error {
	d.logger.Info("enabling trace", "trace_id", traceId)
	_, err := d.sendCommand(fmt.Sprintf("trace %d view on", traceId))
	return err
}

// DisableTrace disables the display of the specified trace.
func (d *Device) DisableTrace(traceId uint) error {
	_, err := d.sendCommand(fmt.Sprintf("trace %d view off", traceId))
	return err
}

// EnableTraceCalc enables trace calculations like TraceCalcMaxH or TraceCalcQuasi for the specified trace.
func (d *Device) EnableTraceCalc(traceId uint, calc TraceCalc) error {
	d.logger.Info("enabling trace calculations", "trace_id", traceId, "calc", calc)

	// calc log and lin is only supported on ultra
	/*if d.model != ModelUltra {
		if calc == TraceCalcLin || calc == TraceCalcLog {
			return ErrOptionNotSupportedByModel
		}
	}*/
	_, err := d.sendCommand(fmt.Sprintf("calc %d %s", traceId, calc.String()))
	return err
}

// DisableTraceCalc disables calculation for the specified trace.
func (d *Device) DisableTraceCalc(traceId uint) error {
	d.logger.Info("disabling trace calculations", "trace_id", traceId)
	_, err := d.sendCommand(fmt.Sprintf("calc %d off", traceId))
	return err
}

// parseTraceValue parses a trace response like `trace <n> value` into a TraceValue struct.
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

// parseTraceStatusLine parses a trace response line like `1: dBm 0.000000000 10.000000000` into a Trace struct.
func parseTraceStatusLine(line string) (Trace, error) {
	fields := strings.Fields(line)
	if len(fields) != 4 {
		return Trace{}, fmt.Errorf("%w: expected 4 fields, got %d", ErrCommandFailed, len(fields))
	}

	trace, err := strconv.Atoi(strings.TrimSuffix(fields[0], ":"))
	if err != nil {
		return Trace{}, fmt.Errorf("%w: failed to parse trace id: %v", ErrCommandFailed, err)
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
		Unit:   DisplayUnit{fields[1]},
		RefPos: refPos,
		Scale:  scale,
	}, nil
}
