package hostfuncs

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPerformExecCommand_Success(t *testing.T) {
	req := ExecCommandRequest{
		Command: "echo",
		Args:    []string{"hello", "world"},
	}

	resp := PerformExecCommand(context.Background(), req)

	assert.Nil(t, resp.Error)
	assert.Equal(t, 0, resp.ExitCode)
	assert.Equal(t, "hello world\n", resp.Stdout)
	assert.Empty(t, resp.Stderr)
	assert.False(t, resp.IsTimeout)
	assert.GreaterOrEqual(t, resp.DurationMs, int64(0))
}

func TestPerformExecCommand_InvalidCommand(t *testing.T) {
	req := ExecCommandRequest{
		Command: "nonexistentcommand12345",
	}

	resp := PerformExecCommand(context.Background(), req)

	require.NotNil(t, resp.Error)
	assert.Equal(t, "EXECUTION_FAILED", resp.Error.Code)
}

func TestPerformExecCommand_Timeout(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping timeout test in short mode")
	}

	req := ExecCommandRequest{
		Command: "sleep",
		Args:    []string{"2"},
		Timeout: 100, // 100ms
	}

	resp := PerformExecCommand(context.Background(), req)

	assert.True(t, resp.IsTimeout)
	assert.Nil(t, resp.Error) // Timeout is a valid result state, not an execution error
}

func TestPerformExecCommand_ExitCode(t *testing.T) {
	// false command returns exit code 1
	req := ExecCommandRequest{
		Command: "false",
	}

	resp := PerformExecCommand(context.Background(), req)

	assert.Nil(t, resp.Error)
	assert.Equal(t, 1, resp.ExitCode)
}

func TestPerformExecCommand_Env(t *testing.T) {
	req := ExecCommandRequest{
		Command: "env",
		Env:     []string{"TEST_VAR=test_value"},
	}

	resp := PerformExecCommand(context.Background(), req)

	assert.Nil(t, resp.Error)
	assert.Contains(t, resp.Stdout, "TEST_VAR=test_value")
}

func TestPerformExecCommand_EmptyCommand(t *testing.T) {
	req := ExecCommandRequest{
		Command: "",
	}

	resp := PerformExecCommand(context.Background(), req)

	require.NotNil(t, resp.Error)
	assert.Equal(t, "INVALID_REQUEST", resp.Error.Code)
}
