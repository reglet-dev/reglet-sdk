// Package sdk provides core types and functions for building Reglet plugins.
package sdk

import (
	"errors"
	"time" // Added for Timestamp

	"github.com/reglet-dev/reglet-sdk/go/wireformat"
)

// Config represents the configuration passed to a plugin observation.
type Config map[string]interface{}

// Evidence represents the structured data returned by a plugin observation.
// This struct directly mirrors the WIT 'evidence' record for direct mapping
// across the WebAssembly boundary.
type Evidence struct {
	Timestamp time.Time
	Error     *ErrorDetail
	Data      map[string]interface{}
	Raw       *string
	Status    bool
}

// ErrorDetail is re-exported from wireformat for backward compatibility.
// Error Types: "network", "timeout", "config", "panic", "capability", "validation", "internal"
type ErrorDetail = wireformat.ErrorDetail

// Metadata contains information about the plugin.
type Metadata struct {
	Name           string       `json:"name"`
	Version        string       `json:"version"`
	Description    string       `json:"description"`
	SDKVersion     string       `json:"sdk_version"`      // Auto-populated
	MinHostVersion string       `json:"min_host_version"` // Minimum compatible host
	Capabilities   []Capability `json:"capabilities"`
}

// Capability describes a permission required by the plugin.
type Capability struct {
	Kind    string `json:"kind"`
	Pattern string `json:"pattern"`
}

// ToErrorDetail converts a Go error to our structured ErrorDetail.
// This function recognizes custom error types and categorizes them appropriately.
// Error types implementing DetailedError will use their own ToErrorDetail method (OCP compliant).
func ToErrorDetail(err error) *ErrorDetail {
	if err == nil {
		return nil
	}

	// If the error is already a *wireformat.ErrorDetail, use it directly.
	var wfError *wireformat.ErrorDetail
	if errors.As(err, &wfError) {
		return wfError
	}

	// Check if error implements DetailedError interface (OCP compliant)
	var detailedErr DetailedError
	if errors.As(err, &detailedErr) {
		return detailedErr.ToErrorDetail()
	}

	// Generic error - categorize as internal
	return &ErrorDetail{
		Message: err.Error(),
		Type:    "internal",
		Code:    "",
	}
}

// Success creates a successful Evidence with data.
func Success(data map[string]interface{}) Evidence {
	return Evidence{Status: true, Data: data, Timestamp: time.Now()}
}

// Failure creates a failed Evidence with an error.
func Failure(errType, message string) Evidence {
	return Evidence{
		Status:    false,
		Error:     &ErrorDetail{Message: message, Type: errType},
		Timestamp: time.Now(),
	}
}

const (
	// Version of the SDK
	Version = "0.1.0-alpha"
	// MinHostVersion is the minimum compatible Reglet host version.
	MinHostVersion = "0.2.0" // Placeholder, will be determined by host capabilities
)
