package tinysa

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"time"

	"go.bug.st/serial"
)

const (
	// responsePrompt is the prompt indicating a response from tinySA.
	responsePrompt = "ch> "

	// commandTerminator is the default terminator used by tinySA for commands, requests, and responses.
	commandTerminator = "\r\n"

	// responseTimeoutTries is the maximum number of retries after a response timeout.
	responseTimeoutTries = 3
)

// sendCommand wraps sendCommandBinary, converting its []byte response to a string.
func sendCommand(logger *slog.Logger, port serial.Port, cmd string, responseTimeout time.Duration) (string, error) {
	response, err := sendCommandBinary(logger, port, cmd, responseTimeout)
	if err != nil {
		return "", err
	}
	return string(response), nil
}

// sendCommandBinary sends a command to the tinySA and handles both string and binary responses.
func sendCommandBinary(logger *slog.Logger, port serial.Port, cmd string, responseTimeout time.Duration) ([]byte, error) {
	fullCmd := cmd + commandTerminator

	logger.Debug("sending command", "cmd", cmd)

	// Send full command and re-try if we run into timeout.
	// This sometimes happens with the tinySA for whatever reason.
	var response bytes.Buffer
	tries := 0
	for {
		r, err := sendCommandAndRead(logger, port, fullCmd, responseTimeout)
		if err != nil {
			if errors.Is(err, ErrCommandResponseTimeout) && tries < responseTimeoutTries {
				logger.Warn("response timeout detected, re-trying", "tries", tries)
				tries++
				continue
			}
			return nil, err
		}
		response = r
		break
	}

	return handleResponse(logger, fullCmd, response.Bytes())
}

// handleResponse handles both string and binary responses.
func handleResponse(logger *slog.Logger, fullCmd string, response []byte) ([]byte, error) {
	// Check if we retrieved (at least) our own command as echo.
	if len(response) < len(fullCmd) {
		logger.Error("response is too short, missing echo", "len", len(response), "response", string(response))
		return nil, fmt.Errorf("response is too short, missing echo")
	}

	// Ensure the response starts with the echoed command.
	if !bytes.HasPrefix(response, []byte(fullCmd)) {
		logger.Error("response does not start with echoed command", "expected", fullCmd, "got", string(response[:len(fullCmd)]))
		return nil, fmt.Errorf("response does not start with echoed command")
	}

	// Remove the echoed command from response.
	response = response[len(fullCmd):]
	logger.Debug("removed echoed command", "response", string(response))

	// Check if the response is long enough to contain the response prompt.
	if len(response) < len(responsePrompt) {
		logger.Error("response is too short, missing response prompt", "len", len(response), "response", string(response))
		return nil, fmt.Errorf("response is too short, missing response prompt")
	}

	// If the response is just the response prompt, it was command without any response, and we are finished here.
	if bytes.Equal(response, []byte(responsePrompt)) {
		logger.Debug("only response prompt found, no additional response")
		return []byte{}, nil
	}

	// If the last received bytes match commandTerminator + responsePrompt, we received a string as response,
	// and remove both. Otherwise, it is a binary response (e.g. `capture`), where we only remove the prompt.
	suffix := commandTerminator + responsePrompt
	if bytes.HasSuffix(response, []byte(suffix)) {
		response = response[:len(response)-len(suffix)]
		logger.Debug("parsed as string response", "response", string(response))
	} else {
		response = response[:len(response)-len(responsePrompt)]
		logger.Debug("parsed as binary response", "response", string(response))
	}

	return response, nil
}

// sendCommandAndRead sends a request over the serial port and reads the response.
func sendCommandAndRead(logger *slog.Logger, port serial.Port, fullCmd string, responseTimeout time.Duration) (bytes.Buffer, error) {
	logger.Debug("sending full command", "cmd", fullCmd)
	if _, err := port.Write([]byte(fullCmd)); err != nil {
		logger.Error("failed to write command", "cmd", fullCmd, "err", err)
		return bytes.Buffer{}, fmt.Errorf("cmd write failed: %s", err.Error())
	}

	buffer := make([]byte, 512)
	var response bytes.Buffer

	timeout := time.Now().Add(responseTimeout)

	logger.Debug("waiting for response")
	for {
		if time.Now().After(timeout) {
			logger.Error("timeout occurred while reading response", "timeout", timeout.String())
			return bytes.Buffer{}, ErrCommandResponseTimeout
		}

		n, err := port.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				logger.Error("port eof occurred while reading response")
				break
			}
			logger.Error("failed to read response", "err", err)
			return bytes.Buffer{}, fmt.Errorf("failed to read response: %s", err.Error())
		}

		response.Write(buffer[:n])
		logger.Debug("read bytes", "n", n, "len", response.Len())

		// Check if we received the response prompt.
		if bytes.HasSuffix(response.Bytes(), []byte(responsePrompt)) {
			logger.Debug("response prompt detected, reading complete", "buffer", string(buffer[:n]))
			break
		}
	}

	logger.Debug("finished reading response", "len", response.Len(), "response", response.String())

	return response, nil
}
