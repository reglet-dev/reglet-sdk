//go:build wasip1

package wasm

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	"github.com/reglet-dev/reglet-sdk/go/domain/entities"
	"github.com/reglet-dev/reglet-sdk/go/domain/ports"
	"github.com/reglet-dev/reglet-sdk/go/internal/abi"
	wasmcontext "github.com/reglet-dev/reglet-sdk/go/internal/wasmcontext"
)

// Compile-time interface compliance check
var _ ports.TCPDialer = (*TCPAdapter)(nil)

// TCPAdapter implements ports.TCPDialer for the WASM environment.
type TCPAdapter struct{}

// NewTCPAdapter creates a new TCP adapter.
func NewTCPAdapter() *TCPAdapter {
	return &TCPAdapter{}
}

// Dial establishes a TCP connection to the given address.
func (a *TCPAdapter) Dial(ctx context.Context, address string) (ports.TCPConnection, error) {
	return a.DialWithTimeout(ctx, address, 5000) // Default 5s timeout
}

// DialWithTimeout establishes a TCP connection with a timeout.
func (a *TCPAdapter) DialWithTimeout(ctx context.Context, address string, timeoutMs int) (ports.TCPConnection, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, fmt.Errorf("invalid address: %w", err)
	}

	request := entities.TCPRequest{
		Context:   wasmcontext.ContextToWire(ctx),
		Host:      host,
		Port:      port,
		TimeoutMs: timeoutMs,
		TLS:       false, // Default to false as ports.TCPDialer interface implies raw TCP
	}

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal TCP request: %w", err)
	}

	requestPacked := abi.PtrFromBytes(requestBytes)
	defer abi.DeallocatePacked(requestPacked)

	responsePacked := host_tcp_connect(requestPacked)

	responseBytes := abi.BytesFromPtr(responsePacked)
	defer abi.DeallocatePacked(responsePacked)

	var response entities.TCPResponse
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal TCP response: %w", err)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("%s: %s", response.Error.Type, response.Error.Message)
	}

	return &WasmTCPConnection{
		response: response,
	}, nil
}

// WasmTCPConnection adapts the WASM response to the TCPConnection interface.
type WasmTCPConnection struct {
	response entities.TCPResponse
}

func (c *WasmTCPConnection) Close() error {
	// Connection is ephemeral in WASM check context, nothing to close
	return nil
}

func (c *WasmTCPConnection) RemoteAddr() string {
	return c.response.RemoteAddr
}

func (c *WasmTCPConnection) IsConnected() bool {
	return c.response.Connected
}
