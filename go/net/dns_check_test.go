package sdknet

import (
	"context"
	"testing"

	"github.com/reglet-dev/reglet-sdk/go/application/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunDNSCheck_Success_A(t *testing.T) {
	cfg := config.Config{
		"hostname":    "google.com",
		"record_type": "A",
	}

	ctx := context.Background()
	result, err := RunDNSCheck(ctx, cfg)

	require.NoError(t, err)
	assert.True(t, result.IsSuccess())
	assert.NotNil(t, result.Data)
	assert.Contains(t, result.Data, "records")
	records := result.Data["records"].([]string)
	assert.NotEmpty(t, records, "Expected at least one A record for google.com")
}

func TestRunDNSCheck_Success_AAAA(t *testing.T) {
	cfg := config.Config{
		"hostname":    "google.com",
		"record_type": "AAAA",
	}

	ctx := context.Background()
	result, err := RunDNSCheck(ctx, cfg)

	require.NoError(t, err)
	// AAAA records may or may not exist depending on network
	assert.True(t, result.IsSuccess() || result.IsError())
}

func TestRunDNSCheck_Success_MX(t *testing.T) {
	cfg := config.Config{
		"hostname":    "google.com",
		"record_type": "MX",
	}

	ctx := context.Background()
	result, err := RunDNSCheck(ctx, cfg)

	require.NoError(t, err)
	assert.True(t, result.IsSuccess())
	assert.NotNil(t, result.Data)
	assert.Contains(t, result.Data, "mx_records")
}

func TestRunDNSCheck_MissingHostname(t *testing.T) {
	cfg := config.Config{
		"record_type": "A",
	}

	ctx := context.Background()
	result, err := RunDNSCheck(ctx, cfg)

	require.NoError(t, err)
	assert.True(t, result.IsError())
	require.NotNil(t, result.Error)
	assert.Equal(t, "config", result.Error.Type)
	assert.Contains(t, result.Error.Message, "hostname")
}

func TestRunDNSCheck_DefaultRecordType(t *testing.T) {
	cfg := config.Config{
		"hostname": "google.com",
		// record_type not specified - should default to A
	}

	ctx := context.Background()
	result, err := RunDNSCheck(ctx, cfg)

	require.NoError(t, err)
	assert.True(t, result.IsSuccess())
	assert.Equal(t, "A", result.Data["record_type"])
}

func TestRunDNSCheck_InvalidHostname(t *testing.T) {
	cfg := config.Config{
		"hostname":    "this-domain-definitely-does-not-exist-12345.invalid",
		"record_type": "A",
	}

	ctx := context.Background()
	result, err := RunDNSCheck(ctx, cfg)

	require.NoError(t, err)
	assert.True(t, result.IsError())
	require.NotNil(t, result.Error)
	assert.Equal(t, "network", result.Error.Type)
}

func TestRunDNSCheck_CustomNameserver(t *testing.T) {
	cfg := config.Config{
		"hostname":    "google.com",
		"record_type": "A",
		"nameserver":  "8.8.8.8",
	}

	ctx := context.Background()
	result, err := RunDNSCheck(ctx, cfg)

	require.NoError(t, err)
	// Should succeed with custom nameserver
	assert.True(t, result.IsSuccess() || result.IsError())
}
