package tinysa

import (
	"errors"
	"fmt"
	"net"
	"time"
)

// tcpTransport adapts a plain net.Conn (typically a TCP connection to a serial-to-TCP
// bridge such as socat) to the transport interface, replicating the non-blocking-read
// timeout semantics that go.bug.st/serial.Port provides over a real serial port: when no
// data arrives within the configured read timeout, Read returns (0, nil) rather than an
// error, so callers in protocol.go can keep polling until their own response deadline.
type tcpTransport struct {
	conn        net.Conn
	readTimeout time.Duration
}

func (t *tcpTransport) Read(p []byte) (int, error) {
	if t.readTimeout > 0 {
		_ = t.conn.SetReadDeadline(time.Now().Add(t.readTimeout))
	} else {
		_ = t.conn.SetReadDeadline(time.Time{})
	}

	n, err := t.conn.Read(p)
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) && netErr.Timeout() {
			return n, nil
		}
		return n, err
	}

	return n, nil
}

func (t *tcpTransport) Write(p []byte) (int, error) {
	return t.conn.Write(p)
}

func (t *tcpTransport) Close() error {
	return t.conn.Close()
}

func (t *tcpTransport) SetReadTimeout(d time.Duration) error {
	t.readTimeout = d
	return nil
}

// NewDeviceTCP creates a *Device by connecting over TCP to a serial-to-TCP bridge (e.g.
// `socat /dev/ttyACM0,raw,echo=0,ispeed=115200,ospeed=115200 TCP-LISTEN:9001,reuseaddr,fork`)
// instead of opening a local serial port directly. addr is a "host:port" pair, e.g. "localhost:9001".
func NewDeviceTCP(addr string, opts ...DeviceOption) (*Device, error) {
	options := defaultDeviceOptions()
	for _, opt := range opts {
		opt(&options)
	}

	if options.logger == nil {
		options.logger = newNoopLogger()
	}

	logger := options.logger

	logger.Debug("initializing new tcp device", "addr", addr, "options", options)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		logger.Error("failed to dial tcp bridge", "addr", addr, "err", err)
		return nil, fmt.Errorf("failed to dial %s: %s", addr, err.Error())
	}

	port := &tcpTransport{conn: conn}

	// set read timeout
	if err = port.SetReadTimeout(options.readTimeout); err != nil {
		logger.Error("failed to set read timeout", "err", err)
		_ = port.Close()
		return nil, fmt.Errorf("failed to set read timeout: %s", err.Error())
	}

	// probe device
	logger.Debug("probing device", "addr", addr)
	pr, err := probeDevice(logger, port, options.responseTimeout)
	if err != nil {
		logger.Error("failed to probe device", "err", err)
		_ = port.Close()
		return nil, fmt.Errorf("failed to probe device: %s", err.Error())
	}

	return createDeviceFromProbe(logger, port, pr, options)
}
