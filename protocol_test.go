package tinysa

import (
	"bytes"
	"testing"
)

func TestHandleResponse(t *testing.T) {
	logger := newNoopLogger()

	tests := []struct {
		name      string
		fullCmd   string
		response  []byte
		want      []byte
		expectErr bool
	}{
		{
			name:     "command with string response",
			fullCmd:  "vbat\r\n",
			response: []byte("vbat\r\n4179 mV\r\nch> "),
			want:     []byte("4179 mV"),
		},
		{
			name:     "command with binary response",
			fullCmd:  "capture\r\n",
			response: []byte("capture\r\n\x01\x02\x03ch> "),
			want:     []byte("\x01\x02\x03"),
		},
		{
			name:     "command with no response",
			fullCmd:  "sweep start 120000000\r\n",
			response: []byte("sweep start 120000000\r\nch> "),
			want:     []byte(""),
		},
		{
			name:      "echo command received, missing command prompt",
			fullCmd:   "sweep start 120000000\r\n",
			response:  []byte("sweep start 120000000\r\n"),
			expectErr: true,
		},
		{
			name:      "nothing received",
			fullCmd:   "sweep start 120000000\r\n",
			response:  []byte(""),
			expectErr: true,
		},
		{
			name:      "does not start with echoed command",
			fullCmd:   "vbat\r\n",
			response:  []byte("garbage4179 mV\r\nch> "),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := handleResponse(logger, tt.fullCmd, tt.response)
			if (err != nil) != tt.expectErr {
				t.Fatalf("handleResponse() error = %v, expectErr %v", err, tt.expectErr)
			}
			if !tt.expectErr && !bytes.Equal(got, tt.want) {
				t.Fatalf("handleResponse() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
