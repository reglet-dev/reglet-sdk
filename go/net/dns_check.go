package sdknet

import (
	"context"
	"fmt"
	"time"

	"github.com/reglet-dev/reglet-sdk/go/application/config"
	"github.com/reglet-dev/reglet-sdk/go/domain/entities"
	"github.com/reglet-dev/reglet-sdk/go/hostfuncs"
)

// RunDNSCheck performs a DNS lookup check.
// It parses configuration, executes the DNS lookup, and returns a structured Result.
//
// Expected config fields:
//   - hostname (string, required): Domain name to resolve
//   - record_type (string, optional): DNS record type (A, AAAA, CNAME, MX, TXT, NS). Default: A
//   - nameserver (string, optional): Custom nameserver (e.g., "8.8.8.8")
//   - timeout_ms (int, optional): Lookup timeout in milliseconds (default: 5000)
//
// Returns a Result with:
//   - Status: "success" if lookup succeeded, "error" if failed
//   - Data: map containing "records" ([]string) or "mx_records" (for MX queries)
//   - Error: structured error details if lookup failed
func RunDNSCheck(ctx context.Context, cfg config.Config) (entities.Result, error) {
	// Parse required fields
	hostname, err := config.MustGetString(cfg, "hostname")
	if err != nil {
		return entities.ResultError(entities.NewErrorDetail("config", err.Error()).WithCode("MISSING_HOSTNAME")), nil
	}

	// Parse optional fields
	recordType := config.GetStringDefault(cfg, "record_type", "A")
	nameserver := config.GetStringDefault(cfg, "nameserver", "")
	timeoutMs := config.GetIntDefault(cfg, "timeout_ms", 5000)

	// Create request
	req := hostfuncs.DNSLookupRequest{
		Hostname:   hostname,
		RecordType: recordType,
		Nameserver: nameserver,
		Timeout:    timeoutMs,
	}

	// Execute DNS lookup
	start := time.Now()
	resp := hostfuncs.PerformDNSLookup(ctx, req)
	metadata := entities.NewRunMetadata(start, time.Now())

	// Build result data
	resultData := make(map[string]any)

	if len(resp.Records) > 0 {
		resultData["records"] = resp.Records
	}
	if len(resp.MXRecords) > 0 {
		// Convert MXRecords to map format for result
		mxRecords := make([]map[string]any, len(resp.MXRecords))
		for i, mx := range resp.MXRecords {
			mxRecords[i] = map[string]any{
				"host": mx.Host,
				"pref": mx.Pref,
			}
		}
		resultData["mx_records"] = mxRecords
	}
	resultData["record_type"] = recordType
	resultData["hostname"] = hostname

	// Return result based on lookup status
	if resp.Error == nil {
		message := fmt.Sprintf("DNS lookup successful for %s (%s)", hostname, recordType)
		return entities.ResultSuccess(message, resultData).WithMetadata(metadata), nil
	}

	// Lookup failed - convert DNS error to ErrorDetail
	errDetail := entities.NewErrorDetail("network", resp.Error.Message).WithCode(resp.Error.Code)
	return entities.ResultError(errDetail).WithMetadata(metadata), nil
}
