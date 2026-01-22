//go:build !wasip1

package plugin

import (
	"context"

	"github.com/reglet-dev/reglet-sdk/go/domain/entities"
)

// Plugin is the interface every Reglet plugin must implement.
type Plugin interface {
	Describe(ctx context.Context) (entities.Metadata, error)
	Schema(ctx context.Context) ([]byte, error)
	Check(ctx context.Context, config map[string]any) (entities.Result, error)
}

// Register is a stub for non-WASM platforms.
func Register(p Plugin) {
	// No-op on non-WASM platforms
}
