package tinysa

import (
	"fmt"
	"go.bug.st/serial"
	"log/slog"
	"regexp"
	"sync"
	"time"
)

// probeResult contains the parsed response of the `version` command, which is used to detect the tinySA model.
type probeResult struct {
	model     string
	version   string
	hwVersion string
}

// probeDevice tries to detect a tinySA device on the given port, returning a probeResult.
//
// We try multiple times to detect the tinySA, because directly after boot we find some malformed output.
// This helps us detect the tinySA reliably, and also clears the input buffer for further commands.
func probeDevice(logger *slog.Logger, port serial.Port, responseTimeout time.Duration) (probeResult, error) {
	var re = regexp.MustCompile(`^(tinySA\w+)_+(\S+)\s*HW Version:(\S+)`)

	var probeResult probeResult

	logger.Debug("probing device")

	found := false
	i := 0
	for i < 3 {
		response, err := sendCommand(logger, port, "version", responseTimeout)
		if err != nil {
			logger.Debug("failed to send version command", "err", err)
			return probeResult, fmt.Errorf("%w: failed sending version command: %w", ErrProbeFailed, err)
		}

		matches := re.FindStringSubmatch(response)
		if len(matches) == 4 {
			found = true
			probeResult.model = matches[1]
			probeResult.version = matches[2]
			probeResult.hwVersion = matches[3]
			logger.Info("found valid device", "probe_result", probeResult)
			break
		}

		i++
	}

	if !found {
		logger.Warn("no valid version response found, might not be a tinySA device")
		return probeResult, fmt.Errorf("%w: no valid version response found, might not be a tinySA device", ErrProbeFailed)
	}

	return probeResult, nil
}

// createDeviceFromProbe creates a new *Device from a probeResult.
func createDeviceFromProbe(logger *slog.Logger, port serial.Port, pr probeResult, opts deviceOptions) (*Device, error) {
	cfg, ok := deviceModels[pr.model]
	if !ok {
		logger.Error("unknown model", "model", pr.model)
		return nil, fmt.Errorf("%w: unknown model %s", ErrConnectionFailed, pr.model)
	}

	return &Device{
		port:            port,
		mutex:           sync.Mutex{},
		model:           cfg.model,
		version:         pr.version,
		hwVersion:       pr.hwVersion,
		width:           cfg.width,
		height:          cfg.height,
		logger:          logger,
		readTimeout:     opts.readTimeout,
		responseTimeout: opts.responseTimeout,
	}, nil
}
