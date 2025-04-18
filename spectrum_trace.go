package tinysa

import (
	"fmt"
	"strconv"
	"strings"
)

type TraceUnit struct {
	value string
}

func (u TraceUnit) String() string {
	return u.value
}

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
	TraceUnitRaw  = TraceUnit{traceUnitRAW}
	TraceUnitDBm  = TraceUnit{traceUnitDBm}
	TraceUnitDBmV = TraceUnit{traceUnitDBmV}
	TraceUnitDBuV = TraceUnit{traceUnitDBuV}
	TraceUnitV    = TraceUnit{traceUnitV}
	TraceUnitVpp  = TraceUnit{traceUnitVpp}
	TraceUnitW    = TraceUnit{traceUnitW}
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

func TraceUnitOptions() []string {
	return traceUnitOptions
}

type Trace struct {
	Trace  uint
	Unit   TraceUnit
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
	// traceCalcLog    string = "log"
	// traceCalcLin    string = "lin"
)

var (
	TraceCalcMinH   = TraceCalc{traceCalcMinH}
	TraceCalcMaxH   = TraceCalc{traceCalcMaxH}
	TraceCalcMaxD   = TraceCalc{traceCalcMaxD}
	TraceCalcAver4  = TraceCalc{traceCalcAver4}
	TraceCalcAver16 = TraceCalc{traceCalcAver16}
	TraceCalcQuasi  = TraceCalc{traceCalcQuasi}
	// TraceCalcLog    = TraceCalc{traceCalcLog} // TODO: Currently not supported (needs extra argument)
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
	// traceCalcLog,
	// traceCalcLin,
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

	result, err := parseTraceResponseLine(line)
	if err != nil {
		d.logger.Error("failed to parse trace result", "trace_id", traceId, "line", line, "err", err)
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
func (d *Device) GetTraceValues(traceId uint) ([]TraceValue, error) {
	d.logger.Info("getting trace values", "trace_id", traceId)

	dataStr, err := d.sendCommand(fmt.Sprintf("trace %d value", traceId))
	if err != nil {
		return nil, err
	}
	dataList := strings.Split(dataStr, commandTerminator)

	data := make([]TraceValue, len(dataList))
	for i, line := range dataList {
		data[i], err = parseTraceValueResponseLine(line)
		if err != nil {
			d.logger.Error("failed to parse trace value", "trace_id", traceId, "line", line)
			return nil, fmt.Errorf("failed to parse trace result: %s", err.Error())
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
func (d *Device) SetTraceScale(level int) error {
	d.logger.Info("setting trace scale", "level", level)
	_, err := d.sendCommand(fmt.Sprintf("trace scale %d", level))
	return err
}
