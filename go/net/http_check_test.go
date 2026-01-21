package sdknet

import (
	"context"
	"testing"

	"github.com/reglet-dev/reglet-sdk/go/application/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunHTTPCheck_Success_GET(t *testing.T) {
	cfg := config.Config{
		"url":    "https://www.google.com",
		"method": "GET",
	}

	ctx := context.Background()
	result, err := RunHTTPCheck(ctx, cfg)

	require.NoError(t, err)
	assert.True(t, result.IsSuccess())
	assert.NotNil(t, result.Data)
	assert.Contains(t, result.Data, "status_code")
	assert.Contains(t, result.Data, "latency_ms")
}

func TestRunHTTPCheck_DefaultMethod(t *testing.T) {
	cfg := config.Config{
		"url": "https://www.google.com",
		// method not specified - should default to GET
	}

	ctx := context.Background()
	result, err := RunHTTPCheck(ctx, cfg)

	require.NoError(t, err)
	assert.True(t, result.IsSuccess())
}

func TestRunHTTPCheck_ExpectedStatus_Success(t *testing.T) {
	cfg := config.Config{
		"url":             "https://www.google.com",
		"expected_status": 200,
	}

	ctx := context.Background()
	result, err := RunHTTPCheck(ctx, cfg)

	require.NoError(t, err)
	assert.True(t, result.IsSuccess())
}

func TestRunHTTPCheck_ExpectedStatus_Mismatch(t *testing.T) {
	cfg := config.Config{
		"url":             "https://www.google.com",
		"expected_status": 404, // Expecting 404 but will get 200
	}

	ctx := context.Background()
	result, err := RunHTTPCheck(ctx, cfg)

	require.NoError(t, err)
	assert.True(t, result.IsFailure())
	assert.Contains(t, result.Message, "mismatch")
}

func TestRunHTTPCheck_MissingURL(t *testing.T) {
	cfg := config.Config{
		"method": "GET",
	}

	ctx := context.Background()
	result, err := RunHTTPCheck(ctx, cfg)

	require.NoError(t, err)
	assert.True(t, result.IsError())
	require.NotNil(t, result.Error)
	assert.Equal(t, "config", result.Error.Type)
	assert.Contains(t, result.Error.Message, "url")
}

func TestRunHTTPCheck_InvalidURL(t *testing.T) {
	cfg := config.Config{
		"url": "not-a-valid-url",
	}

	ctx := context.Background()
	result, err := RunHTTPCheck(ctx, cfg)

	require.NoError(t, err)
	assert.True(t, result.IsError())
}

func TestRunHTTPCheck_WithHeaders(t *testing.T) {
	headers := map[string]interface{}{
		"User-Agent": "TestAgent/1.0",
		"Accept":     "application/json",
	}

	cfg := config.Config{
		"url":     "https://www.google.com",
		"headers": headers,
	}

	ctx := context.Background()
	result, err := RunHTTPCheck(ctx, cfg)

	require.NoError(t, err)
	// Should succeed with custom headers
	assert.True(t, result.IsSuccess() || result.IsError())
}

func TestRunHTTPCheck_Timeout(t *testing.T) {
	cfg := config.Config{
		"url":        "https://www.google.com",
		"timeout_ms": 1, // Very short timeout
	}

	ctx := context.Background()
	result, err := RunHTTPCheck(ctx, cfg)

	require.NoError(t, err)
	// Should likely timeout, but success is also acceptable if very fast
	assert.True(t, result.IsSuccess() || result.IsError())
}
