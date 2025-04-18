package tinysa

import "errors"

// ErrCommandResponseTimeout is returned when a command does not receive a response within the expected timeframe.
var ErrCommandResponseTimeout = errors.New("command response timeout")
