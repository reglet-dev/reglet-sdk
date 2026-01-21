//go:build wasip1

package exec

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/reglet-dev/reglet-sdk/go/domain/entities"
	"github.com/reglet-dev/reglet-sdk/go/internal/abi"
	wasmcontext "github.com/reglet-dev/reglet-sdk/go/internal/wasmcontext"
)

// runConfig holds the configuration for command execution.
// This struct is unexported to enforce the functional options pattern.
type runConfig struct {
	workdir string        // Working directory for command (default: inherit)
	env     []string      // Environment variables (default: inherit)
	timeout time.Duration // Execution timeout (default: 30s)
}

// defaultRunConfig returns secure defaults for command execution.
// These defaults align with constitution requirements for secure-by-default.
func defaultRunConfig() runConfig {
	return runConfig{
		workdir: "",               // Empty = inherit from host
		env:     nil,              // nil = inherit from host
		timeout: 30 * time.Second, // 30s default per spec
	}
}

// RunOption is a functional option for configuring command execution.
// Use With* functions to create options.
type RunOption func(*runConfig)

// WithWorkdir sets the working directory for the command.
// If not specified, the host's current working directory is used.
func WithWorkdir(dir string) RunOption {
	return func(c *runConfig) {
		c.workdir = dir
	}
}

// WithEnv sets the environment variables for the command.
// Each entry should be in KEY=VALUE format.
// If not specified, the host's environment is inherited.
func WithEnv(env []string) RunOption {
	return func(c *runConfig) {
		c.env = env
	}
}

// WithExecTimeout sets the execution timeout for the command.
// Default is 30 seconds per constitution specification.
// A zero or negative duration is ignored (uses default).
func WithExecTimeout(d time.Duration) RunOption {
	return func(c *runConfig) {
		if d > 0 {
			c.timeout = d
		}
	}
}

// applyRunOptions applies functional options and returns the configuration.
// This is used by the Run function to process variadic options.
func applyRunOptions(opts ...RunOption) runConfig {
	cfg := defaultRunConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
	return cfg
}

//go:wasmimport reglet_host exec_command
func host_exec_command(reqPacked uint64) uint64

// CommandRequest defines the parameters for executing a command.
type CommandRequest struct {
	Command string
	Args    []string
	Dir     string
	Env     []string
	Timeout int // seconds
}

// CommandResponse contains the result of the command execution.
type CommandResponse struct {
	Stdout     string
	Stderr     string
	ExitCode   int
	DurationMs int64 // Execution duration in milliseconds
	IsTimeout  bool  // True if command timed out
}

// Run executes a command on the host system.
// Requires "exec:<command>" capability.
//
// Options:
//   - WithWorkdir(dir): Set working directory (default: inherit from host)
//   - WithEnv(env): Set environment variables (default: inherit from host)
//   - WithExecTimeout(d): Set execution timeout (default: 30s)
//
// Example:
//
//	resp, err := exec.Run(ctx, exec.CommandRequest{
//	    Command: "ls",
//	    Args:    []string{"-la"},
//	}, exec.WithWorkdir("/tmp"), exec.WithExecTimeout(10*time.Second))
func Run(ctx context.Context, req CommandRequest, opts ...RunOption) (*CommandResponse, error) {
	// Apply functional options
	cfg := applyRunOptions(opts...)

	// Override request fields from options if not already set
	if req.Dir == "" && cfg.workdir != "" {
		req.Dir = cfg.workdir
	}
	if req.Env == nil && cfg.env != nil {
		req.Env = cfg.env
	}
	if req.Timeout == 0 {
		req.Timeout = int(cfg.timeout.Seconds())
	}

	// 1. Prepare wire request with context
	wireReq := entities.ExecRequest{
		Context: wasmcontext.ContextToWire(ctx),
		Command: req.Command,
		Args:    req.Args,
		Dir:     req.Dir,
		Env:     req.Env,
	}

	reqData, err := json.Marshal(wireReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// 2. Send to host
	reqPacked := abi.PtrFromBytes(reqData)
	defer abi.DeallocatePacked(reqPacked)

	resPacked := host_exec_command(reqPacked)

	// 3. Read response
	resBytes := abi.BytesFromPtr(resPacked)
	if resBytes == nil {
		// This might happen if host returns 0 packed (null/empty), which indicates error usually handled inside hostWriteResponse but if 0 returned...
		// Wait, hostWriteResponse returns 0 on alloc failure.
		return nil, fmt.Errorf("host returned null response")
	}
	defer abi.DeallocatePacked(resPacked) // Free host-allocated response memory

	var wireRes entities.ExecResponse
	if err := json.Unmarshal(resBytes, &wireRes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// 4. Handle errors
	if wireRes.Error != nil {
		return nil, wireRes.Error
	}

	return &CommandResponse{
		Stdout:     wireRes.Stdout,
		Stderr:     wireRes.Stderr,
		ExitCode:   wireRes.ExitCode,
		DurationMs: wireRes.DurationMs,
		IsTimeout:  wireRes.IsTimeout,
	}, nil
}
