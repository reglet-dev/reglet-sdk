//go:build wasip1

package sdknet

import (
	"time"
)

// resolverConfig holds the configuration for a WasmResolver.
// This struct is unexported to enforce the functional options pattern.
type resolverConfig struct {
	nameserver string        // DNS nameserver address (e.g., "8.8.8.8:53")
	timeout    time.Duration // Timeout for DNS queries (default: 5s)
	retries    int           // Number of retry attempts (default: 3)
}

// defaultResolverConfig returns secure defaults for DNS resolution.
// These defaults align with constitution requirements for secure-by-default.
func defaultResolverConfig() resolverConfig {
	return resolverConfig{
		nameserver: "", // Empty = use host's default resolver
		timeout:    5 * time.Second,
		retries:    3,
	}
}

// ResolverOption is a functional option for configuring a WasmResolver.
// Use With* functions to create options.
type ResolverOption func(*resolverConfig)

// WithNameserver sets a custom DNS nameserver address.
// The address should be in "host:port" format (e.g., "8.8.8.8:53").
// If port is omitted, ":53" is assumed.
// If not specified, the host's default resolver is used.
func WithNameserver(ns string) ResolverOption {
	return func(c *resolverConfig) {
		c.nameserver = ns
	}
}

// WithDNSTimeout sets the timeout for DNS queries.
// Default is 5 seconds per constitution specification.
// A zero or negative duration is ignored (uses default).
func WithDNSTimeout(d time.Duration) ResolverOption {
	return func(c *resolverConfig) {
		if d > 0 {
			c.timeout = d
		}
	}
}

// WithRetries sets the number of retry attempts for failed DNS queries.
// Default is 3 retries. A negative value is ignored (uses default).
// Setting to 0 disables retries (single attempt only).
func WithRetries(n int) ResolverOption {
	return func(c *resolverConfig) {
		if n >= 0 {
			c.retries = n
		}
	}
}

// NewResolver creates a new WasmResolver with the given options.
// Without any options, secure defaults are applied:
//   - timeout: 5 seconds
//   - retries: 3
//   - nameserver: host's default resolver
//
// Example:
//
//	// Use defaults
//	resolver := NewResolver()
//
//	// With custom nameserver and timeout
//	resolver := NewResolver(
//	    WithNameserver("1.1.1.1:53"),
//	    WithDNSTimeout(10*time.Second),
//	)
func NewResolver(opts ...ResolverOption) *WasmResolver {
	cfg := defaultResolverConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
	return &WasmResolver{
		Nameserver: cfg.nameserver,
		timeout:    cfg.timeout,
		retries:    cfg.retries,
	}
}
