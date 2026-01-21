//go:build wasip1

package sdknet

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransportOption_DefaultConfig(t *testing.T) {
	cfg := defaultTransportConfig()

	assert.Equal(t, 30*time.Second, cfg.timeout, "default timeout should be 30s")
	assert.Equal(t, 10, cfg.maxRedirects, "default maxRedirects should be 10")
	assert.Nil(t, cfg.tlsConfig, "default tlsConfig should be nil")
}

func TestNewTransport_WithDefaults(t *testing.T) {
	transport := NewTransport()

	require.NotNil(t, transport, "NewTransport should return a non-nil transport")
	assert.Equal(t, 30*time.Second, transport.timeout, "default timeout should be 30s")
	assert.Equal(t, 10, transport.maxRedirects, "default maxRedirects should be 10")
	assert.Nil(t, transport.tlsConfig, "default tlsConfig should be nil")
}

func TestNewTransport_WithHTTPTimeout(t *testing.T) {
	transport := NewTransport(WithHTTPTimeout(60 * time.Second))

	assert.Equal(t, 60*time.Second, transport.timeout)
	// Other defaults should still apply
	assert.Equal(t, 10, transport.maxRedirects)
}

func TestNewTransport_WithHTTPTimeout_IgnoresInvalid(t *testing.T) {
	tests := []struct {
		name     string
		timeout  time.Duration
		expected time.Duration
	}{
		{"zero should use default", 0, 30 * time.Second},
		{"negative should use default", -1 * time.Second, 30 * time.Second},
		{"positive should be applied", 45 * time.Second, 45 * time.Second},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			transport := NewTransport(WithHTTPTimeout(tc.timeout))
			assert.Equal(t, tc.expected, transport.timeout)
		})
	}
}

func TestNewTransport_WithMaxRedirects(t *testing.T) {
	transport := NewTransport(WithMaxRedirects(5))

	assert.Equal(t, 5, transport.maxRedirects)
}

func TestNewTransport_WithMaxRedirects_ZeroDisables(t *testing.T) {
	transport := NewTransport(WithMaxRedirects(0))

	assert.Equal(t, 0, transport.maxRedirects, "maxRedirects=0 should disable redirects")
}

func TestNewTransport_WithMaxRedirects_IgnoresNegative(t *testing.T) {
	transport := NewTransport(WithMaxRedirects(-1))

	assert.Equal(t, 10, transport.maxRedirects, "negative maxRedirects should use default")
}

func TestNewTransport_MultipleOptions(t *testing.T) {
	transport := NewTransport(
		WithHTTPTimeout(60*time.Second),
		WithMaxRedirects(3),
	)

	assert.Equal(t, 60*time.Second, transport.timeout)
	assert.Equal(t, 3, transport.maxRedirects)
}

func TestNewTransport_OptionsApplyInOrder(t *testing.T) {
	// Last option wins for the same setting
	transport := NewTransport(
		WithMaxRedirects(5),
		WithMaxRedirects(2), // This should override
	)

	assert.Equal(t, 2, transport.maxRedirects, "last option should win")
}
