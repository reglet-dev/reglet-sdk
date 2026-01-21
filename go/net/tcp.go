// Package net provides high-level network check functions for SDK plugins.
package sdknet

import (
	"context"
	"fmt"
	"time"

	"github.com/reglet-dev/reglet-sdk/go/application/config"
	"github.com/reglet-dev/reglet-sdk/go/domain/entities"
	"github.com/reglet-dev/reglet-sdk/go/domain/ports"
	"github.com/reglet-dev/reglet-sdk/go/infrastructure/wasm"
)

// TCPCheckOption is a functional option for configuring TCP checks.
type TCPCheckOption func(*tcpCheckConfig)

type tcpCheckConfig struct {
	dialer ports.TCPDialer
}

func defaultTCPCheckConfig() tcpCheckConfig {
	return tcpCheckConfig{
		dialer: wasm.NewTCPAdapter(),
	}
}

// WithTCPDialer sets the TCP dialer to use for the check.
// This is useful for injecting mocks during testing.
func WithTCPDialer(d ports.TCPDialer) TCPCheckOption {
	return func(c *tcpCheckConfig) {
		if d != nil {
			c.dialer = d
		}
	}
}

// RunTCPCheck performs a TCP connectivity check.
// It parses configuration, executes the connection test, and returns a structured Result.
//
// Expected config fields:
//   - host (string, required): Target hostname or IP address
//   - port (int, required): Target port number (1-65535)
//   - timeout_ms (int, optional): Connection timeout in milliseconds (default: 5000)
//
// Returns a Result with:
//   - Status: "success" if connected, "error" if failed
//   - Data: map containing "connected", "remote_addr", "latency_ms"
//   - Error: structured error details if connection failed
func RunTCPCheck(ctx context.Context, cfg config.Config, opts ...TCPCheckOption) (entities.Result, error) {
	// Parse required fields
	host, err := config.MustGetString(cfg, "host")
	if err != nil {
		return entities.ResultError(entities.NewErrorDetail("config", err.Error()).WithCode("MISSING_HOST")), nil
	}

	port, err := config.MustGetInt(cfg, "port")
	if err != nil {
		return entities.ResultError(entities.NewErrorDetail("config", err.Error()).WithCode("MISSING_PORT")), nil
	}

	// Validate port range
	if port < 1 || port > 65535 {
		return entities.ResultError(entities.NewErrorDetail("config", fmt.Sprintf("invalid port: %d (must be 1-65535)", port)).WithCode("INVALID_PORT")), nil
	}

	// Parse optional timeout
	timeoutMs := config.GetIntDefault(cfg, "timeout_ms", 5000)

	// Configure check
	checkCfg := defaultTCPCheckConfig()
	for _, opt := range opts {
		opt(&checkCfg)
	}

	address := fmt.Sprintf("%s:%d", host, port)

	// Execute TCP connect
	start := time.Now()
	conn, err := checkCfg.dialer.DialWithTimeout(ctx, address, timeoutMs)
	latency := time.Since(start)

	// Create metadata
	metadata := entities.NewRunMetadata(start, time.Now())

	if err != nil {
		// Connection failed - convert error to ErrorDetail
		// Note: We lost granular error codes from hostfuncs unless we parse error string
		// or if err satisfies an interface. For now we use generic "CONNECTION_FAILED"
		// or try to match common strings if critical.
		errDetail := entities.NewErrorDetail("network", err.Error()).WithCode("CONNECTION_FAILED")
		return entities.ResultError(errDetail).WithMetadata(metadata), errDetail
	}

	defer func() { _ = conn.Close() }()

	// Build result data
	resultData := map[string]any{
		"connected":  conn.IsConnected(),
		"latency_ms": latency.Milliseconds(),
	}

	if conn.RemoteAddr() != "" {
		resultData["remote_addr"] = conn.RemoteAddr()
	}

	if conn.IsConnected() {
		return entities.ResultSuccess("TCP connection successful", resultData).WithMetadata(metadata), nil
	}

	// Connected is false but no error? Should not happen normally if Dial returns nil err.
	return entities.ResultError(entities.NewErrorDetail("network", "TCP connection failed").WithCode("CONNECTION_FAILED")).WithMetadata(metadata), nil
}
