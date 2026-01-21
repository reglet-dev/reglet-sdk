// Package net provides high-level network check functions for SDK plugins.
package sdknet

import (
	"context"
	"fmt"
	"time"

	"github.com/reglet-dev/reglet-sdk/go/application/config"
	"github.com/reglet-dev/reglet-sdk/go/domain/entities"
	"github.com/reglet-dev/reglet-sdk/go/hostfuncs"
)

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
func RunTCPCheck(ctx context.Context, cfg config.Config) (entities.Result, error) {
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

	// Create request
	req := hostfuncs.TCPConnectRequest{
		Host:    host,
		Port:    port,
		Timeout: timeoutMs,
	}

	// Execute TCP connect
	start := time.Now()
	resp := hostfuncs.PerformTCPConnect(ctx, req)

	// Build result data
	resultData := map[string]any{
		"connected":  resp.Connected,
		"latency_ms": resp.LatencyMs,
	}

	if resp.RemoteAddr != "" {
		resultData["remote_addr"] = resp.RemoteAddr
	}

	// Create metadata
	metadata := entities.NewRunMetadata(start, time.Now())

	// Return result based on connection status
	if resp.Connected {
		return entities.ResultSuccess("TCP connection successful", resultData).WithMetadata(metadata), nil
	}

	// Connection failed - convert TCP error to ErrorDetail
	if resp.Error != nil {
		errDetail := entities.NewErrorDetail("network", resp.Error.Message).WithCode(resp.Error.Code)
		return entities.ResultError(errDetail).WithMetadata(metadata), errDetail
	}

	// Fallback error
	return entities.ResultError(entities.NewErrorDetail("network", "TCP connection failed").WithCode("CONNECTION_FAILED")).WithMetadata(metadata), nil
}
