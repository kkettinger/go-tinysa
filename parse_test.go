package tinysa

import (
	"testing"
)

func TestParseBatteryVoltageLine(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      uint
		shouldErr bool
	}{
		{name: "valid input", input: "3700 mV", want: 3700, shouldErr: false},
		{name: "missing voltage suffix", input: "3700", want: 0, shouldErr: true},
		{name: "extra field", input: "3700 mV ch>", want: 0, shouldErr: true},
		{name: "non-integer voltage", input: "bad mV", want: 0, shouldErr: true},
		{name: "empty input", input: "", want: 0, shouldErr: true},
		{name: "negative voltage", input: "-10 mV", want: 0, shouldErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseBatteryVoltageLine(tt.input)
			if (err != nil) != tt.shouldErr {
				t.Errorf("parseBatteryVoltageLine(%q) error = %v, wantErr = %v", tt.input, err, tt.shouldErr)
			}
			if got != tt.want {
				t.Errorf("parseBatteryVoltageLine(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseMarkerResultLine(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      Marker
		shouldErr bool
	}{
		{
			name:  "valid input",
			input: "1 216 522167037 -9.08e+01",
			want: Marker{
				Marker:    1,
				Index:     216,
				Frequency: 522167037,
				Value:     -90.8,
			},
			shouldErr: false,
		},
		{
			name:      "too few fields",
			input:     "1 216 522167037",
			shouldErr: true,
		},
		{
			name:      "too many fields",
			input:     "1 216 522167037 -90.8 extra",
			shouldErr: true,
		},
		{
			name:      "non-integer marker",
			input:     "x 216 522167037 -90.8",
			shouldErr: true,
		},
		{
			name:      "non-integer index",
			input:     "1 abc 522167037 -90.8",
			shouldErr: true,
		},
		{
			name:      "non-integer frequency",
			input:     "1 216 xyz -90.8",
			shouldErr: true,
		},
		{
			name:      "non-float value",
			input:     "1 216 522167037 not_a_float",
			shouldErr: true,
		},
		{
			name:      "empty input",
			input:     "",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseMarkerResultLine(tt.input)
			if (err != nil) != tt.shouldErr {
				t.Errorf("parseMarkerResultLine(%q) error = %v, wantErr = %v", tt.input, err, tt.shouldErr)
			}
			if !tt.shouldErr && got != tt.want {
				t.Errorf("parseMarkerResultLine(%q) = %+v, want %+v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseSweepResponse(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      Sweep
		shouldErr bool
	}{
		{
			name:  "valid input",
			input: "450000000 600000000 450",
			want: Sweep{
				Start:  450000000,
				Stop:   600000000,
				Points: 450,
			},
			shouldErr: false,
		},
		{
			name:      "too few fields",
			input:     "450000000 600000000",
			shouldErr: true,
		},
		{
			name:      "too many fields",
			input:     "450000000 600000000 450 extra",
			shouldErr: true,
		},
		{
			name:      "non-integer start",
			input:     "start 600000000 450",
			shouldErr: true,
		},
		{
			name:      "non-integer stop",
			input:     "450000000 stop 450",
			shouldErr: true,
		},
		{
			name:      "non-integer points",
			input:     "450000000 600000000 many",
			shouldErr: true,
		},
		{
			name:      "empty input",
			input:     "",
			shouldErr: true,
		},
		{
			name:  "zero values",
			input: "0 0 0",
			want: Sweep{
				Start:  0,
				Stop:   0,
				Points: 0,
			},
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSweepResponse(tt.input)
			if (err != nil) != tt.shouldErr {
				t.Errorf("parseSweepResponse(%q) error = %v, wantErr = %v", tt.input, err, tt.shouldErr)
			}
			if !tt.shouldErr && got != tt.want {
				t.Errorf("parseSweepResponse(%q) = %+v, want %+v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseTraceValue(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      TraceValue
		shouldErr bool
	}{
		{
			name:  "valid input",
			input: "trace 1 at 216 -123.45",
			want: TraceValue{
				Trace: 1,
				Point: 216,
				Value: -123.45,
			},
			shouldErr: false,
		},
		{
			name:      "too few fields",
			input:     "trace 1 at 216",
			shouldErr: true,
		},
		{
			name:      "too many fields",
			input:     "trace 1 at 216 -123.45 extra",
			shouldErr: true,
		},
		{
			name:      "non-numeric trace",
			input:     "trace x at 216 -123.45",
			shouldErr: true,
		},
		{
			name:      "non-numeric point",
			input:     "trace 1 at x -123.45",
			shouldErr: true,
		},
		{
			name:      "non-numeric value",
			input:     "trace 1 at 216 bad",
			shouldErr: true,
		},
		{
			name:      "empty input",
			input:     "",
			shouldErr: true,
		},
		{
			name:  "zero values",
			input: "trace 0 at 0 0.0",
			want: TraceValue{
				Trace: 0,
				Point: 0,
				Value: 0.0,
			},
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTraceValue(tt.input)
			if (err != nil) != tt.shouldErr {
				t.Errorf("parseTraceValue(%q) error = %v, wantErr = %v", tt.input, err, tt.shouldErr)
			}
			if !tt.shouldErr && got != tt.want {
				t.Errorf("parseTraceValue(%q) = %+v, want %+v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseTraceStatusLine(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      Trace
		shouldErr bool
	}{
		{
			name:  "valid input dBm",
			input: "1: dBm -30.000000000 10.000000000",
			want: Trace{
				Trace:  1,
				Unit:   TraceUnitDBm,
				RefPos: -30.0,
				Scale:  10.0,
			},
			shouldErr: false,
		},
		{
			name:  "valid input Vpp",
			input: "2: Vpp 0.000004500 0.000000450",
			want: Trace{
				Trace:  2,
				Unit:   TraceUnitVpp,
				RefPos: 0.0000045,
				Scale:  0.00000045,
			},
			shouldErr: false,
		},
		{
			name:      "invalid unit",
			input:     "1: UNKNOWN 0.0 10.0",
			shouldErr: true,
		},
		{
			name:      "non-integer trace ID",
			input:     "abc: dBm 0.0 10.0",
			shouldErr: true,
		},
		{
			name:      "non-float refPos",
			input:     "1: dBm notafloat 10.0",
			shouldErr: true,
		},
		{
			name:      "non-float scale",
			input:     "1: dBm 0.0 notafloat",
			shouldErr: true,
		},
		{
			name:      "too few fields",
			input:     "1: dBm 0.0",
			shouldErr: true,
		},
		{
			name:      "too many fields",
			input:     "1: dBm 0.0 10.0 extra",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTraceStatusLine(tt.input)
			if (err != nil) != tt.shouldErr {
				t.Errorf("parseTraceStatusLine(%q) error = %v, wantErr = %v", tt.input, err, tt.shouldErr)
			}
			if !tt.shouldErr && got != tt.want {
				t.Errorf("parseTraceStatusLine(%q) = %+v, want %+v", tt.input, got, tt.want)
			}
		})
	}
}

func TestSweepModeFromString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected SweepMode
		valid    bool
	}{
		// Valid inputs with exact case
		{name: "normal exact case", input: "normal", expected: SweepModeNormal, valid: true},
		{name: "precise exact case", input: "precise", expected: SweepModePrecise, valid: true},
		{name: "fast exact case", input: "fast", expected: SweepModeFast, valid: true},
		{name: "noise exact case", input: "noise", expected: SweepModeNoise, valid: true},

		// Valid inputs with different case
		{name: "normal uppercase", input: "NORMAL", expected: SweepModeNormal, valid: true},
		{name: "precise mixed case", input: "PrEcIsE", expected: SweepModePrecise, valid: true},
		{name: "fast lowercase", input: "fast", expected: SweepModeFast, valid: true},
		{name: "noise title case", input: "Noise", expected: SweepModeNoise, valid: true},

		// Invalid inputs
		{name: "invalid input", input: "invalid", expected: SweepMode{}, valid: false},
		{name: "empty string", input: "", expected: SweepMode{}, valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := SweepModeFromString(tt.input)
			if ok != tt.valid {
				t.Errorf("SweepModeFromString(%q) returned valid=%v, want %v", tt.input, ok, tt.valid)
				return
			}
			if result.String() != tt.expected.String() {
				t.Errorf("SweepModeFromString(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTraceUnitFromString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected TraceUnit
		valid    bool
	}{
		// Valid inputs with exact case
		{name: "RAW exact case", input: "RAW", expected: TraceUnitRaw, valid: true},
		{name: "dBm exact case", input: "dBm", expected: TraceUnitDBm, valid: true},
		{name: "dBmV exact case", input: "dBmV", expected: TraceUnitDBmV, valid: true},
		{name: "dBuV exact case", input: "dBuV", expected: TraceUnitDBuV, valid: true},
		{name: "V exact case", input: "V", expected: TraceUnitV, valid: true},
		{name: "Vpp exact case", input: "Vpp", expected: TraceUnitVpp, valid: true},
		{name: "W exact case", input: "W", expected: TraceUnitW, valid: true},

		// Valid inputs with different case
		{name: "raw lowercase", input: "raw", expected: TraceUnitRaw, valid: true},
		{name: "DBM uppercase", input: "DBM", expected: TraceUnitDBm, valid: true},
		{name: "dbmv lowercase", input: "dbmv", expected: TraceUnitDBmV, valid: true},
		{name: "DbUv mixed case", input: "DbUv", expected: TraceUnitDBuV, valid: true},
		{name: "v lowercase", input: "v", expected: TraceUnitV, valid: true},
		{name: "VPP uppercase", input: "VPP", expected: TraceUnitVpp, valid: true},
		{name: "w lowercase", input: "w", expected: TraceUnitW, valid: true},

		// Invalid inputs
		{name: "invalid input", input: "invalid", expected: TraceUnit{}, valid: false},
		{name: "empty string", input: "", expected: TraceUnit{}, valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := TraceUnitFromString(tt.input)
			if ok != tt.valid {
				t.Errorf("TraceUnitFromString(%q) returned valid=%v, want %v", tt.input, ok, tt.valid)
				return
			}
			if result.String() != tt.expected.String() {
				t.Errorf("TraceUnitFromString(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTraceCalcFromString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected TraceCalc
		valid    bool
	}{
		// Valid inputs with exact case
		{name: "minh exact case", input: "minh", expected: TraceCalcMinH, valid: true},
		{name: "maxh exact case", input: "maxh", expected: TraceCalcMaxH, valid: true},
		{name: "maxd exact case", input: "maxd", expected: TraceCalcMaxD, valid: true},
		{name: "aver4 exact case", input: "aver4", expected: TraceCalcAver4, valid: true},
		{name: "aver16 exact case", input: "aver16", expected: TraceCalcAver16, valid: true},
		{name: "quasi exact case", input: "quasi", expected: TraceCalcQuasi, valid: true},

		// Valid inputs with different case
		{name: "MINH uppercase", input: "MINH", expected: TraceCalcMinH, valid: true},
		{name: "MaxH mixed case", input: "MaxH", expected: TraceCalcMaxH, valid: true},
		{name: "maxd lowercase", input: "maxd", expected: TraceCalcMaxD, valid: true},
		{name: "AVER4 uppercase", input: "AVER4", expected: TraceCalcAver4, valid: true},
		{name: "Aver16 title case", input: "Aver16", expected: TraceCalcAver16, valid: true},
		{name: "QUASI uppercase", input: "QUASI", expected: TraceCalcQuasi, valid: true},

		// Invalid inputs
		{name: "invalid input", input: "invalid", expected: TraceCalc{}, valid: false},
		{name: "empty string", input: "", expected: TraceCalc{}, valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := TraceCalcFromString(tt.input)
			if ok != tt.valid {
				t.Errorf("TraceCalcFromString(%q) returned valid=%v, want %v", tt.input, ok, tt.valid)
				return
			}
			if result.String() != tt.expected.String() {
				t.Errorf("TraceCalcFromString(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
