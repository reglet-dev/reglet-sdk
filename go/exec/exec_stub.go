//go:build !wasip1

package exec

import (
	"time"
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
