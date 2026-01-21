//go:build wasip1

package sdknet

import (
	"crypto/tls"
	"time"
)

// transportConfig holds the configuration for a WasmTransport.
// This struct is unexported to enforce the functional options pattern.
type transportConfig struct {
	timeout      time.Duration // HTTP request timeout (default: 30s)
	maxRedirects int           // Maximum number of redirects to follow (default: 10)
	tlsConfig    *tls.Config   // Optional TLS configuration
}

// defaultTransportConfig returns secure defaults for HTTP transport.
// These defaults align with constitution requirements for secure-by-default.
func defaultTransportConfig() transportConfig {
	return transportConfig{
		timeout:      30 * time.Second,
		maxRedirects: 10,
		tlsConfig:    nil, // Use system defaults
	}
}

// TransportOption is a functional option for configuring a WasmTransport.
// Use With* functions to create options.
type TransportOption func(*transportConfig)

// WithHTTPTimeout sets the timeout for HTTP requests.
// Default is 30 seconds per constitution specification.
// A zero or negative duration is ignored (uses default).
func WithHTTPTimeout(d time.Duration) TransportOption {
	return func(c *transportConfig) {
		if d > 0 {
			c.timeout = d
		}
	}
}

// WithMaxRedirects sets the maximum number of redirects to follow.
// Default is 10 redirects. A negative value is ignored (uses default).
// Setting to 0 disables following redirects.
func WithMaxRedirects(n int) TransportOption {
	return func(c *transportConfig) {
		if n >= 0 {
			c.maxRedirects = n
		}
	}
}

// WithTLSConfig sets a custom TLS configuration.
// If nil is passed, the system default TLS configuration is used.
func WithTLSConfig(cfg *tls.Config) TransportOption {
	return func(c *transportConfig) {
		c.tlsConfig = cfg
	}
}

// NewTransport creates a new WasmTransport with the given options.
// Without any options, secure defaults are applied:
//   - timeout: 30 seconds
//   - maxRedirects: 10
//   - tlsConfig: system defaults
//
// Example:
//
//	// Use defaults
//	transport := NewTransport()
//
//	// With custom timeout
//	transport := NewTransport(
//	    WithHTTPTimeout(60*time.Second),
//	    WithMaxRedirects(5),
//	)
func NewTransport(opts ...TransportOption) *WasmTransport {
	cfg := defaultTransportConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
	return &WasmTransport{
		timeout:      cfg.timeout,
		maxRedirects: cfg.maxRedirects,
		tlsConfig:    cfg.tlsConfig,
	}
}
