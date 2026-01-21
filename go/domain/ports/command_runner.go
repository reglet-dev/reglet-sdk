package ports

import (
	"context"

	"github.com/reglet-dev/reglet-sdk/go/domain/entities"
)

// CommandRunner defines the interface for command execution.
// Infrastructure adapters implement this to provide exec functionality.
type CommandRunner interface {
	// Run executes a command and returns the result.
	Run(ctx context.Context, command string, args []string, opts ...RunOption) (*entities.Result, error)
}

// RunOption is a functional option for configuring command execution.
// This mirrors the exec package's RunOption for port compatibility.
type RunOption func(*RunConfig)

// RunConfig holds configuration for command execution.
type RunConfig struct {
	Workdir string
	Env     []string
	Timeout int // milliseconds
}
