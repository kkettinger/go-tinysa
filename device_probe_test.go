package tinysa

import "testing"

func TestMatchProbeResponse(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		want        probeResult
		expectError bool
	}{
		{
			name:  "normal response zs-405",
			input: "tinySA4_v1.4-197-gaa78ccc\r\nHW Version:V0.4.5.1 ",
			want: probeResult{
				model:     "tinySA4",
				version:   "1.4-197-gaa78ccc",
				hwVersion: "0.4.5.1",
			},
			expectError: false,
		},
		{
			name:  "normal response zs-407",
			input: "tinySA4_v1.4-197-gaa78ccc\r\nHW Version:V0.5.4 max2871",
			want: probeResult{
				model:     "tinySA4",
				version:   "1.4-197-gaa78ccc",
				hwVersion: "0.5.4 max2871",
			},
			expectError: false,
		},
		{
			name:  "response with custom flashed firmware (no version)",
			input: "tinySA4_\r\nHW Version:V0.4.5.1 ",
			want: probeResult{
				model:     "tinySA4",
				version:   "",
				hwVersion: "0.4.5.1",
			},
			expectError: false,
		},
		{
			name:        "empty response",
			input:       "",
			want:        probeResult{},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseVersionResponse(tc.input)
			if (err != nil) != tc.expectError {
				t.Errorf("expected error: %v, got: %v", tc.expectError, err)
			}
			if got != tc.want {
				t.Errorf("expected result: %+v, got: %+v", tc.want, got)
			}
		})
	}
}
