//go:build wasip1

package sdknet

import (
	"context"

	wasmcontext "github.com/reglet-dev/reglet-sdk/go/internal/wasmcontext"
)

// createContextWireFormat extracts relevant info from a Go context into the wire format.
// This is now a wrapper around wasmcontext.ContextToWire for backwards compatibility.
func createContextWireFormat(ctx context.Context) ContextWireFormat {
	return wasmcontext.ContextToWire(ctx)
}
