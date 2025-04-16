package tinysa

import (
	"io"
	"log/slog"
)

// NewNoopLogger creates a no-op logger that discards all output.
func NewNoopLogger() *slog.Logger {
	handler := slog.NewTextHandler(io.Discard, nil)
	return slog.New(handler)
}
