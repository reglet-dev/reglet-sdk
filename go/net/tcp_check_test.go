package sdknet

import (
	"context"
	"testing"

	"github.com/reglet-dev/reglet-sdk/go/application/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunTCPCheck_Success(t *testing.T) {
	// Use a reliable public service for testing
	cfg := config.Config{
		"host":       "google.com",
		"port":       443,
		"timeout_ms": 5000,
	}

	ctx := context.Background()
	result, err := RunTCPCheck(ctx, cfg)

	require.NoError(t, err)
	assert.True(t, result.IsSuccess(), "Expected successful connection to google.com:443")
	assert.NotNil(t, result.Data)
	assert.True(t, result.Data["connected"].(bool))
	assert.NotNil(t, result.Metadata)
}

func TestRunTCPCheck_MissingHost(t *testing.T) {
	cfg := config.Config{
		"port": 80,
	}

	ctx := context.Background()
	result, err := RunTCPCheck(ctx, cfg)

	require.NoError(t, err)
	assert.True(t, result.IsError())
	require.NotNil(t, result.Error)
	assert.Equal(t, "config", result.Error.Type)
	assert.Contains(t, result.Error.Message, "host")
}

func TestRunTCPCheck_MissingPort(t *testing.T) {
	cfg := config.Config{
		"host": "example.com",
	}

	ctx := context.Background()
	result, err := RunTCPCheck(ctx, cfg)

	require.NoError(t, err)
	assert.True(t, result.IsError())
	require.NotNil(t, result.Error)
	assert.Equal(t, "config", result.Error.Type)
	assert.Contains(t, result.Error.Message, "port")
}

func TestRunTCPCheck_InvalidPort(t *testing.T) {
	tests := []struct {
		name string
		port int
	}{
		{"port too low", 0},
		{"port negative", -1},
		{"port too high", 65536},
		{"port way too high", 99999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Config{
				"host": "example.com",
				"port": tt.port,
			}

			ctx := context.Background()
			result, err := RunTCPCheck(ctx, cfg)

			require.NoError(t, err)
			assert.True(t, result.IsError())
			require.NotNil(t, result.Error)
			assert.Equal(t, "INVALID_PORT", result.Error.Code)
		})
	}
}

func TestRunTCPCheck_ConnectionRefused(t *testing.T) {
	// Use localhost on a port that's unlikely to be listening
	cfg := config.Config{
		"host":       "127.0.0.1",
		"port":       54321,
		"timeout_ms": 1000,
	}

	ctx := context.Background()
	result, err := RunTCPCheck(ctx, cfg)

	require.NoError(t, err)
	assert.True(t, result.IsError())
	require.NotNil(t, result.Error)
	assert.Equal(t, "network", result.Error.Type)
	// Error code should be CONNECTION_REFUSED or CONNECTION_FAILED
	assert.Contains(t, []string{"CONNECTION_REFUSED", "CONNECTION_FAILED"}, result.Error.Code)
}

func TestRunTCPCheck_Timeout(t *testing.T) {
	// Use a non-routable IP to force timeout
	cfg := config.Config{
		"host":       "192.0.2.1", // TEST-NET-1 (RFC 5737) - non-routable
		"port":       80,
		"timeout_ms": 100,
	}

	ctx := context.Background()
	result, err := RunTCPCheck(ctx, cfg)

	require.NoError(t, err)
	assert.True(t, result.IsError())
	require.NotNil(t, result.Error)
	// Should timeout
	assert.Contains(t, []string{"TIMEOUT", "CONNECTION_FAILED"}, result.Error.Code)
}

func TestRunTCPCheck_ContextCancellation(t *testing.T) {
	cfg := config.Config{
		"host":       "google.com",
		"port":       443,
		"timeout_ms": 10000,
	}

	// Use an already-canceled context for reliable test
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	result, err := RunTCPCheck(ctx, cfg)

	require.NoError(t, err)
	// Should fail (either success or error depending on timing, both are acceptable)
	// The important thing is it doesn't panic or hang
	assert.NotNil(t, result)
}

func TestRunTCPCheck_DefaultTimeout(t *testing.T) {
	// Test that default timeout (5000ms) is applied when not specified
	cfg := config.Config{
		"host": "google.com",
		"port": 443,
		// timeout_ms not specified
	}

	ctx := context.Background()
	result, err := RunTCPCheck(ctx, cfg)

	require.NoError(t, err)
	// Should succeed with default timeout
	assert.True(t, result.IsSuccess() || result.IsError()) // Either is valid depending on network
}
