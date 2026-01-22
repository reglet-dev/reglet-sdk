package hostfuncs

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"time"
)

// ExecCommandRequest contains parameters for a command execution.
type ExecCommandRequest struct {
	// Command is the command to execute.
	Command string `json:"command"`

	// Args contains command arguments.
	Args []string `json:"args"`

	// Dir is the working directory.
	Dir string `json:"dir,omitempty"`

	// Env contains environment variables (KEY=VALUE).
	Env []string `json:"env,omitempty"`

	// Timeout is the execution timeout in milliseconds. Default is 30000 (30s).
	Timeout int `json:"timeout_ms,omitempty"`
}

// ExecCommandResponse contains the result of a command execution.
type ExecCommandResponse struct {
	// Error contains error information if execution failed to start.
	Error *ExecError `json:"error,omitempty"`

	// Stdout is the standard output.
	Stdout string `json:"stdout"`

	// Stderr is the standard error.
	Stderr string `json:"stderr"`

	// DurationMs is the execution duration in milliseconds.
	DurationMs int64 `json:"duration_ms,omitempty"`

	// ExitCode is the exit code.
	ExitCode int `json:"exit_code"`

	// IsTimeout indicates if the command timed out.
	IsTimeout bool `json:"is_timeout,omitempty"`
}

// ExecError represents an execution error.
type ExecError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error implements the error interface.
func (e *ExecError) Error() string {
	return e.Message
}

// ExecOption is a functional option for configuring execution behavior.
type ExecOption func(*execConfig)

type execConfig struct {
	timeout time.Duration
}

func defaultExecConfig() execConfig {
	return execConfig{
		timeout: 30 * time.Second,
	}
}

// WithExecTimeout sets the execution timeout.
func WithExecTimeout(d time.Duration) ExecOption {
	return func(c *execConfig) {
		if d > 0 {
			c.timeout = d
		}
	}
}

// PerformExecCommand executes a command on the host.
// This is a pure Go implementation with no WASM runtime dependencies.
func PerformExecCommand(ctx context.Context, req ExecCommandRequest, opts ...ExecOption) ExecCommandResponse {
	cfg := defaultExecConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	if req.Timeout > 0 {
		cfg.timeout = time.Duration(req.Timeout) * time.Millisecond
	}

	// Validate request
	if req.Command == "" {
		return ExecCommandResponse{
			Error: &ExecError{
				Code:    "INVALID_REQUEST",
				Message: "command is required",
			},
		}
	}

	// Apply timeout to context
	ctx, cancel := context.WithTimeout(ctx, cfg.timeout)
	defer cancel()

	//nolint:gosec // G204: Command execution is the purpose of this function
	cmd := exec.CommandContext(ctx, req.Command, req.Args...)
	if req.Dir != "" {
		cmd.Dir = req.Dir
	}
	if len(req.Env) > 0 {
		cmd.Env = req.Env
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err := cmd.Run()
	duration := time.Since(start)

	resp := ExecCommandResponse{
		Stdout:     stdout.String(),
		Stderr:     stderr.String(),
		DurationMs: duration.Milliseconds(),
	}

	if err != nil {
		// Check for timeout
		if ctx.Err() == context.DeadlineExceeded {
			resp.IsTimeout = true
			resp.ExitCode = -1 // Conventional timeout code
			return resp
		}

		// Check for exit code
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			resp.ExitCode = exitErr.ExitCode()
			return resp
		}

		// Failed to start or other error
		resp.Error = &ExecError{
			Code:    "EXECUTION_FAILED",
			Message: err.Error(),
		}
		return resp
	}

	resp.ExitCode = 0
	return resp
}
