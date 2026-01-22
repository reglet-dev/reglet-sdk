package ports

import (
	"context"
)

// TCPDialer defines the interface for TCP connection operations.
// Infrastructure adapters implement this to provide TCP functionality.
type TCPDialer interface {
	// Dial establishes a TCP connection to the given address.
	Dial(ctx context.Context, address string) (TCPConnection, error)

	// DialWithTimeout establishes a TCP connection with a timeout.
	DialWithTimeout(ctx context.Context, address string, timeoutMs int) (TCPConnection, error)
}

// TCPConnection represents an established TCP connection.
type TCPConnection interface {
	// Close closes the connection.
	Close() error

	// RemoteAddr returns the remote address.
	RemoteAddr() string

	// IsConnected returns true if the connection is established.
	IsConnected() bool
}
