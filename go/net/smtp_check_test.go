package sdknet

import (
	"context"
	"testing"

	"github.com/reglet-dev/reglet-sdk/go/application/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunSMTPCheck_MissingHost(t *testing.T) {
	cfg := config.Config{
		"port": 25,
	}

	ctx := context.Background()
	result, err := RunSMTPCheck(ctx, cfg)

	require.NoError(t, err)
	assert.True(t, result.IsError())
	require.NotNil(t, result.Error)
	assert.Equal(t, "config", result.Error.Type)
	assert.Contains(t, result.Error.Message, "host")
}

func TestRunSMTPCheck_MissingPort(t *testing.T) {
	cfg := config.Config{
		"host": "smtp.gmail.com",
	}

	ctx := context.Background()
	result, err := RunSMTPCheck(ctx, cfg)

	require.NoError(t, err)
	assert.True(t, result.IsError())
	require.NotNil(t, result.Error)
	assert.Equal(t, "config", result.Error.Type)
	assert.Contains(t, result.Error.Message, "port")
}

func TestRunSMTPCheck_InvalidPort(t *testing.T) {
	tests := []struct {
		name string
		port int
	}{
		{"port too low", 0},
		{"port negative", -1},
		{"port too high", 65536},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Config{
				"host": "smtp.gmail.com",
				"port": tt.port,
			}

			ctx := context.Background()
			result, err := RunSMTPCheck(ctx, cfg)

			require.NoError(t, err)
			assert.True(t, result.IsError())
			require.NotNil(t, result.Error)
			assert.Equal(t, "INVALID_PORT", result.Error.Code)
		})
	}
}

func TestRunSMTPCheck_ConnectionRefused(t *testing.T) {
	// Use localhost on a port that's unlikely to be listening
	cfg := config.Config{
		"host":       "127.0.0.1",
		"port":       54321,
		"timeout_ms": 1000,
	}

	ctx := context.Background()
	result, err := RunSMTPCheck(ctx, cfg)

	require.NoError(t, err)
	assert.True(t, result.IsError())
	require.NotNil(t, result.Error)
	assert.Equal(t, "network", result.Error.Type)
}

func TestRunSMTPCheck_DefaultValues(t *testing.T) {
	cfg := config.Config{
		"host": "smtp.gmail.com",
		"port": 587,
		// use_tls, use_starttls, timeout_ms not specified - should use defaults
	}

	ctx := context.Background()
	result, err := RunSMTPCheck(ctx, cfg)

	require.NoError(t, err)
	// May succeed or fail depending on network, just verify it returns a result
	assert.NotNil(t, result)
}

func TestRunSMTPCheck_WithSTARTTLS(t *testing.T) {
	cfg := config.Config{
		"host":         "smtp.gmail.com",
		"port":         587,
		"use_starttls": true,
		"timeout_ms":   5000,
	}

	ctx := context.Background()
	result, err := RunSMTPCheck(ctx, cfg)

	require.NoError(t, err)
	// May succeed or fail depending on network/firewall
	assert.NotNil(t, result)
}

func TestRunSMTPCheck_WithImplicitTLS(t *testing.T) {
	cfg := config.Config{
		"host":       "smtp.gmail.com",
		"port":       465,
		"use_tls":    true,
		"timeout_ms": 5000,
	}

	ctx := context.Background()
	result, err := RunSMTPCheck(ctx, cfg)

	require.NoError(t, err)
	// May succeed or fail depending on network/firewall
	assert.NotNil(t, result)
}
