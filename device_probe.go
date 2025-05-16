package tinysa

import (
	"fmt"
	"log/slog"
	"regexp"
	"sync"
	"time"

	"go.bug.st/serial"
)

// probeResult contains the parsed response of the `version` command, which is used to detect the tinySA model.
type probeResult struct {
	model     string
	version   string
	hwVersion string
}

// probeDevice tries to detect a tinySA device on the given port, returning a probeResult.
func probeDevice(logger *slog.Logger, port serial.Port, responseTimeout time.Duration) (probeResult, error) {
	logger.Debug("probing device")

	// We try multiple times to detect the tinySA, because directly after boot we find some malformed output.
	// This helps us detect the tinySA reliably and also clears the input buffer for further commands.
	i := 0
	for i < 3 {
		response, _ := sendCommand(logger, port, "version", responseTimeout)

		if pr, err := parseVersionResponse(response); err == nil {
			logger.Info("found valid device", "probe_result", pr)
			return pr, nil
		}

		i++
	}

	err := fmt.Errorf("no valid version response found, might not be a tinySA device")
	logger.Warn(err.Error())
	return probeResult{}, err
}

// parseVersionResponse matches the response of the `version` command and returns a probeResult.
func parseVersionResponse(response string) (probeResult, error) {
	var re = regexp.MustCompile(`^(tinySA\w+)_v?(\S+)?\s*HW Version:V(.*?)\s*$`)

	matches := re.FindStringSubmatch(response)
	if len(matches) != 4 {
		return probeResult{}, fmt.Errorf("invalid probe response")
	}

	if matches[1] == "" {
		return probeResult{}, fmt.Errorf("invalid probe response")
	}

	return probeResult{
		model:     matches[1],
		version:   matches[2],
		hwVersion: matches[3],
	}, nil
}

// createDeviceFromProbe creates a new *Device from a probeResult.
func createDeviceFromProbe(logger *slog.Logger, port serial.Port, pr probeResult, opts deviceOptions) (*Device, error) {
	cfg, ok := deviceModels[pr.model]
	if !ok {
		logger.Error("unknown model", "model", pr.model)
		return nil, fmt.Errorf("unknown model %s", pr.model)
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
