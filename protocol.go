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
		if r, err := sendCommandBinaryInner(logger, port, fullCmd, responseTimeout); err != nil {
			if errors.Is(err, ErrCommandResponseTimeout) && tries < responseTimeoutTries {
				logger.Debug("response timeout detected, re-trying", "tries", tries)
				tries++
				continue
			}
			return nil, err
		} else {
			response = r
			break
		}
	}

	var responseBytes = response.Bytes()

	// Check if we retrieved our echo.
	if len(responseBytes) < len(fullCmd) {
		logger.Debug("response is too short, missing echo", "len", len(responseBytes))
		return nil, fmt.Errorf("response is too short, missing echo")
	}

	// Remove the echoed command.
	responseBytes = responseBytes[len(fullCmd):]

	logger.Debug("removed echoed command", "response", string(responseBytes))

	// If the response is just the response prompt, it was command without any response,
	// and we are finished here.
	if len(responseBytes) == len(responsePrompt) && string(responseBytes) == responsePrompt {
		logger.Debug("found only response prompt, finished")
		return []byte{}, nil
	}

	// If the last received bytes match commandTerminator + responsePrompt, we received a string as response,
	// and remove both. Otherwise, it is a binary response (e.g. `capture`), where we only remove the prompt.
	responseSuffixStr := commandTerminator + responsePrompt
	if string(responseBytes[len(responseBytes)-len(responseSuffixStr):]) == responseSuffixStr {
		responseBytes = responseBytes[:len(responseBytes)-len(responseSuffixStr)]
		logger.Debug("parsed as string response", "response", string(responseBytes))
	} else {
		responseBytes = responseBytes[:len(responseBytes)-len(responsePrompt)]
		logger.Debug("parsed as binary response", "response", string(responseBytes))
	}

	return responseBytes, nil
}

// sendCommandBinaryInner sends a request over the serial port and reads the response.
func sendCommandBinaryInner(logger *slog.Logger, port serial.Port, fullCmd string, responseTimeout time.Duration) (bytes.Buffer, error) {
	logger.Debug("sending full command", "cmd", fullCmd)
	if _, err := port.Write([]byte(fullCmd)); err != nil {
		logger.Debug("failed to write command", "cmd", fullCmd, "err", err)
		return bytes.Buffer{}, fmt.Errorf("cmd write failed: %v", err)
	}

	buffer := make([]byte, 512)
	var response bytes.Buffer

	timeout := time.Now().Add(responseTimeout)

	logger.Debug("waiting for response")
	for {
		if time.Now().After(timeout) {
			logger.Debug("timeout occurred while reading response", "timeout", timeout.String())
			return bytes.Buffer{}, ErrCommandResponseTimeout
		}

		n, err := port.Read(buffer)
		if err != nil {
			if err == io.EOF {
				logger.Debug("port eof occurred while reading response")
				break
			}
			logger.Debug("failed to read response", "err", err)
			return bytes.Buffer{}, fmt.Errorf("failed to read response: %v", err)
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
