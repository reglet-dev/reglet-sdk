package hostfuncs

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

// TCPConnectRequest contains parameters for a TCP connection test.
type TCPConnectRequest struct {
	// Host is the target hostname or IP address.
	Host string `json:"host"`

	// Port is the target port number.
	Port int `json:"port"`

	// Timeout is the connection timeout in milliseconds. Default is 5000 (5s).
	Timeout int `json:"timeout_ms,omitempty"`
}

// TCPConnectResponse contains the result of a TCP connection test.
type TCPConnectResponse struct {
	// Error contains error information if the connection failed.
	Error *TCPError `json:"error,omitempty"`

	// RemoteAddr is the resolved remote address if connected.
	RemoteAddr string `json:"remote_addr,omitempty"`

	// LatencyMs is the connection latency in milliseconds.
	LatencyMs int64 `json:"latency_ms,omitempty"`

	// Connected indicates whether the connection was successful.
	Connected bool `json:"connected"`
}

// TCPError represents a TCP connection error.
type TCPError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error implements the error interface.
func (e *TCPError) Error() string {
	return e.Message
}

// TCPOption is a functional option for configuring TCP connection behavior.
type TCPOption func(*tcpConfig)

type tcpConfig struct {
	timeout time.Duration
}

func defaultTCPConfig() tcpConfig {
	return tcpConfig{
		timeout: 5 * time.Second,
	}
}

// WithTCPTimeout sets the TCP connection timeout.
func WithTCPTimeout(d time.Duration) TCPOption {
	return func(c *tcpConfig) {
		if d > 0 {
			c.timeout = d
		}
	}
}

// PerformTCPConnect tests TCP connectivity to the specified host and port.
// This is a pure Go implementation with no WASM runtime dependencies.
//
// Example usage from a WASM host:
//
//	func handleTCPConnect(req hostfuncs.TCPConnectRequest) hostfuncs.TCPConnectResponse {
//	    return hostfuncs.PerformTCPConnect(ctx, req)
//	}
func PerformTCPConnect(ctx context.Context, req TCPConnectRequest, opts ...TCPOption) TCPConnectResponse {
	cfg := defaultTCPConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	// Override config from request if specified
	if req.Timeout > 0 {
		cfg.timeout = time.Duration(req.Timeout) * time.Millisecond
	}

	// Validate request
	if req.Host == "" {
		return TCPConnectResponse{
			Connected: false,
			Error: &TCPError{
				Code:    "INVALID_REQUEST",
				Message: "host is required",
			},
		}
	}
	if req.Port <= 0 || req.Port > 65535 {
		return TCPConnectResponse{
			Connected: false,
			Error: &TCPError{
				Code:    "INVALID_REQUEST",
				Message: fmt.Sprintf("invalid port: %d", req.Port),
			},
		}
	}

	// Build address
	address := fmt.Sprintf("%s:%d", req.Host, req.Port)

	// Apply timeout to context
	ctx, cancel := context.WithTimeout(ctx, cfg.timeout)
	defer cancel()

	// Attempt connection
	start := time.Now()
	dialer := &net.Dialer{
		Timeout: cfg.timeout,
	}

	conn, err := dialer.DialContext(ctx, "tcp", address)
	latency := time.Since(start)

	if err != nil {
		code := "CONNECTION_FAILED"
		switch {
		case strings.Contains(err.Error(), "timeout"), ctx.Err() == context.DeadlineExceeded:
			code = "TIMEOUT"
		case strings.Contains(err.Error(), "refused"):
			code = "CONNECTION_REFUSED"
		case strings.Contains(err.Error(), "no such host"):
			code = "HOST_NOT_FOUND"
		}

		return TCPConnectResponse{
			Connected: false,
			LatencyMs: latency.Milliseconds(),
			Error: &TCPError{
				Code:    code,
				Message: err.Error(),
			},
		}
	}
	defer func() { _ = conn.Close() }()

	return TCPConnectResponse{
		Connected:  true,
		RemoteAddr: conn.RemoteAddr().String(),
		LatencyMs:  latency.Milliseconds(),
	}
}
