package tinysa

import "errors"

// TODO: use errors in codebase

var ErrConnectionFailed = errors.New("connection failed")
var ErrCommandFailed = errors.New("command failed")
var ErrProbeFailed = errors.New("probe failed")
var ErrUnexpectedCommandResponse = errors.New("unexpected command response")
var ErrCommandNotSupportedByModel = errors.New("command not supported by model")
var ErrOptionNotSupportedByModel = errors.New("option not supported by model")
var ErrCommandResponseTimeout = errors.New("command response timeout")
