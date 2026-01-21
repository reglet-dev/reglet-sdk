//go:build wasip1

package sdknet

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolverOption_DefaultConfig(t *testing.T) {
	cfg := defaultResolverConfig()

	assert.Empty(t, cfg.nameserver, "default nameserver should be empty (use host default)")
	assert.Equal(t, 5*time.Second, cfg.timeout, "default timeout should be 5s")
	assert.Equal(t, 3, cfg.retries, "default retries should be 3")
}

func TestNewResolver_WithDefaults(t *testing.T) {
	resolver := NewResolver()

	require.NotNil(t, resolver, "NewResolver should return a non-nil resolver")
	assert.Empty(t, resolver.Nameserver, "default nameserver should be empty")
	assert.Equal(t, 5*time.Second, resolver.timeout, "default timeout should be 5s")
	assert.Equal(t, 3, resolver.retries, "default retries should be 3")
}

func TestNewResolver_WithNameserver(t *testing.T) {
	resolver := NewResolver(WithNameserver("8.8.8.8:53"))

	assert.Equal(t, "8.8.8.8:53", resolver.Nameserver)
	// Other defaults should still apply
	assert.Equal(t, 5*time.Second, resolver.timeout)
	assert.Equal(t, 3, resolver.retries)
}

func TestNewResolver_WithDNSTimeout(t *testing.T) {
	resolver := NewResolver(WithDNSTimeout(10 * time.Second))

	assert.Equal(t, 10*time.Second, resolver.timeout)
	// Other defaults should still apply
	assert.Equal(t, 3, resolver.retries)
}

func TestNewResolver_WithDNSTimeout_IgnoresInvalid(t *testing.T) {
	tests := []struct {
		name     string
		timeout  time.Duration
		expected time.Duration
	}{
		{"zero should use default", 0, 5 * time.Second},
		{"negative should use default", -1 * time.Second, 5 * time.Second},
		{"positive should be applied", 15 * time.Second, 15 * time.Second},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resolver := NewResolver(WithDNSTimeout(tc.timeout))
			assert.Equal(t, tc.expected, resolver.timeout)
		})
	}
}

func TestNewResolver_WithRetries(t *testing.T) {
	resolver := NewResolver(WithRetries(5))

	assert.Equal(t, 5, resolver.retries)
}

func TestNewResolver_WithRetries_ZeroDisables(t *testing.T) {
	resolver := NewResolver(WithRetries(0))

	assert.Equal(t, 0, resolver.retries, "retries=0 should disable retries")
}

func TestNewResolver_WithRetries_IgnoresNegative(t *testing.T) {
	resolver := NewResolver(WithRetries(-1))

	assert.Equal(t, 3, resolver.retries, "negative retries should use default")
}

func TestNewResolver_MultipleOptions(t *testing.T) {
	resolver := NewResolver(
		WithNameserver("1.1.1.1:53"),
		WithDNSTimeout(8*time.Second),
		WithRetries(2),
	)

	assert.Equal(t, "1.1.1.1:53", resolver.Nameserver)
	assert.Equal(t, 8*time.Second, resolver.timeout)
	assert.Equal(t, 2, resolver.retries)
}

func TestNewResolver_OptionsApplyInOrder(t *testing.T) {
	// Last option wins for the same setting
	resolver := NewResolver(
		WithNameserver("8.8.8.8:53"),
		WithNameserver("1.1.1.1:53"), // This should override
	)

	assert.Equal(t, "1.1.1.1:53", resolver.Nameserver, "last option should win")
}
