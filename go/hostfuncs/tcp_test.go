package hostfuncs

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPerformTCPConnect_Success(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping network test in short mode")
	}

	req := TCPConnectRequest{
		Host:    "example.com",
		Port:    80,
		Timeout: 5000,
	}

	resp := PerformTCPConnect(context.Background(), req)

	assert.True(t, resp.Connected, "should connect to example.com:80")
	assert.NotEmpty(t, resp.RemoteAddr, "should have remote address")
	assert.Nil(t, resp.Error)
	assert.Greater(t, resp.LatencyMs, int64(0), "should have non-zero latency")
}

func TestPerformTCPConnect_InvalidHost(t *testing.T) {
	req := TCPConnectRequest{
		Host: "",
		Port: 80,
	}

	resp := PerformTCPConnect(context.Background(), req)

	require.NotNil(t, resp.Error)
	assert.Equal(t, "INVALID_REQUEST", resp.Error.Code)
	assert.False(t, resp.Connected)
}

func TestPerformTCPConnect_InvalidPort(t *testing.T) {
	tests := []struct {
		name string
		port int
	}{
		{"zero port", 0},
		{"negative port", -1},
		{"port too high", 65536},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := TCPConnectRequest{
				Host: "example.com",
				Port: tc.port,
			}

			resp := PerformTCPConnect(context.Background(), req)

			require.NotNil(t, resp.Error)
			assert.Equal(t, "INVALID_REQUEST", resp.Error.Code)
		})
	}
}

func TestPerformTCPConnect_ConnectionRefused(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping network test in short mode")
	}

	req := TCPConnectRequest{
		Host:    "127.0.0.1",
		Port:    59999, // Unlikely to have a service here
		Timeout: 1000,
	}

	resp := PerformTCPConnect(context.Background(), req)

	assert.False(t, resp.Connected)
	require.NotNil(t, resp.Error)
	// Error code depends on whether port is closed or firewalled
	assert.Contains(t, []string{"CONNECTION_REFUSED", "TIMEOUT", "CONNECTION_FAILED"}, resp.Error.Code)
}

func TestTCPConnectRequest_Fields(t *testing.T) {
	req := TCPConnectRequest{
		Host:    "test.example.com",
		Port:    443,
		Timeout: 10000,
	}

	assert.Equal(t, "test.example.com", req.Host)
	assert.Equal(t, 443, req.Port)
	assert.Equal(t, 10000, req.Timeout)
}

func TestTCPError_Error(t *testing.T) {
	err := &TCPError{
		Code:    "TEST_CODE",
		Message: "test error message",
	}

	assert.Equal(t, "test error message", err.Error())
}

func TestDefaultTCPConfig(t *testing.T) {
	cfg := defaultTCPConfig()

	assert.Equal(t, 5*time.Second, cfg.timeout)
}

func TestWithTCPTimeout(t *testing.T) {
	cfg := defaultTCPConfig()
	opt := WithTCPTimeout(10 * time.Second)
	opt(&cfg)

	assert.Equal(t, 10*time.Second, cfg.timeout)
}

func TestWithTCPTimeout_IgnoresInvalid(t *testing.T) {
	cfg := defaultTCPConfig()
	opt := WithTCPTimeout(-1 * time.Second)
	opt(&cfg)

	assert.Equal(t, 5*time.Second, cfg.timeout, "should keep default for negative duration")
}
