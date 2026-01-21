package sdknet

import (
	"context"
	"fmt"
	"time"

	"github.com/reglet-dev/reglet-sdk/go/application/config"
	"github.com/reglet-dev/reglet-sdk/go/domain/entities"
	"github.com/reglet-dev/reglet-sdk/go/hostfuncs"
)

// RunSMTPCheck performs an SMTP connectivity check.
// It parses configuration, executes the SMTP connection test, and returns a structured Result.
//
// Expected config fields:
//   - host (string, required): SMTP server hostname
//   - port (int, required): SMTP server port (typically 25, 465, or 587)
//   - use_tls (bool, optional): Use implicit TLS (port 465). Default: false
//   - use_starttls (bool, optional): Upgrade to TLS via STARTTLS (port 587). Default: false
//   - timeout_ms (int, optional): Connection timeout in milliseconds (default: 30000)
//
// Returns a Result with:
//   - Status: "success" if connected, "error" if failed
//   - Data: map containing "connected", "banner", "tls_version", "latency_ms"
//   - Error: structured error details if connection failed
func RunSMTPCheck(ctx context.Context, cfg config.Config) (entities.Result, error) {
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

	// Parse optional fields
	useTLS := config.GetBoolDefault(cfg, "use_tls", false)
	useSTARTTLS := config.GetBoolDefault(cfg, "use_starttls", false)
	timeoutMs := config.GetIntDefault(cfg, "timeout_ms", 30000)

	// Create request
	req := hostfuncs.SMTPConnectRequest{
		Host:        host,
		Port:        port,
		UseTLS:      useTLS,
		UseSTARTTLS: useSTARTTLS,
		Timeout:     timeoutMs,
	}

	// Execute SMTP connect
	start := time.Now()
	resp := hostfuncs.PerformSMTPConnect(ctx, req)
	metadata := entities.NewRunMetadata(start, time.Now())

	// Build result data
	resultData := map[string]any{
		"connected":  resp.Connected,
		"latency_ms": resp.LatencyMs,
	}

	if resp.Banner != "" {
		resultData["banner"] = resp.Banner
	}
	if resp.TLSVersion != "" {
		resultData["tls_version"] = resp.TLSVersion
	}

	// Return result based on connection status
	if resp.Connected {
		message := fmt.Sprintf("SMTP connection successful to %s:%d", host, port)
		return entities.ResultSuccess(message, resultData).WithMetadata(metadata), nil
	}

	// Connection failed - convert SMTP error to ErrorDetail
	if resp.Error != nil {
		errDetail := entities.NewErrorDetail("network", resp.Error.Message).WithCode(resp.Error.Code)
		return entities.ResultError(errDetail).WithMetadata(metadata), errDetail
	}

	// Fallback error
	return entities.ResultError(entities.NewErrorDetail("network", "SMTP connection failed").WithCode("CONNECTION_FAILED")).WithMetadata(metadata), nil
}
