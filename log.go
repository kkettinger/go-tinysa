package tinysa

import (
	"io"
	"log/slog"
)

// newNoopLogger creates a no-op logger that discards all output.
func newNoopLogger() *slog.Logger {
	handler := slog.NewTextHandler(io.Discard, nil)
	return slog.New(handler)
}
