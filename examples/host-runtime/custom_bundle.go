package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/reglet-dev/reglet-sdk/go/hostfuncs"
)

// TLSCheckRequest is the request format for the custom host function.
type TLSCheckRequest struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	TimeoutMs int    `json:"timeout_ms"`
}

// TLSCheckResponse is the response format for the custom host function.
type TLSCheckResponse struct {
	Connected bool      `json:"connected"`
	NotAfter  string    `json:"not_after,omitempty"` // RFC3339
	Issuer    string    `json:"issuer,omitempty"`
	Error     *TLSError `json:"error,omitempty"`
}

// TLSError represents a failure in the TLS check.
type TLSError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// customTLSBundle implements hostfuncs.HostFuncBundle.
type customTLSBundle struct {
	handlers map[string]hostfuncs.ByteHandler
}

func (b *customTLSBundle) Handlers() map[string]hostfuncs.ByteHandler {
	return b.handlers
}

// CustomTLSBundle returns a bundle containing the "tls_check" host function.
func CustomTLSBundle() hostfuncs.HostFuncBundle {
	return &customTLSBundle{
		handlers: map[string]hostfuncs.ByteHandler{
			"tls_check": hostfuncs.NewJSONHandler(func(ctx context.Context, req TLSCheckRequest) TLSCheckResponse {
				return performTLSCheck(ctx, req)
			}),
		},
	}
}

// performTLSCheck executes the actual TLS handshake logic on the host.
func performTLSCheck(ctx context.Context, req TLSCheckRequest) TLSCheckResponse {
	// Validate input
	if req.Host == "" {
		return TLSCheckResponse{Error: &TLSError{Code: "INVALID_INPUT", Message: "host is required"}}
	}
	if req.Port <= 0 {
		req.Port = 443
	}
	timeout := time.Duration(req.TimeoutMs) * time.Millisecond
	if timeout == 0 {
		timeout = 5 * time.Second
	}

	address := fmt.Sprintf("%s:%d", req.Host, req.Port)

	// Use a standard net.Dialer with timeout
	netDialer := &net.Dialer{
		Timeout: timeout,
	}

	// Dial TLS
	conn, err := tls.DialWithDialer(netDialer, "tcp", address, &tls.Config{
		InsecureSkipVerify: false, // Enforce valid certs for compliance demo
		ServerName:         req.Host,
	})
	if err != nil {
		return TLSCheckResponse{
			Connected: false,
			Error: &TLSError{
				Code:    "CONNECTION_FAILED",
				Message: err.Error(),
			},
		}
	}
	_ = conn.Close()

	// The following code snippet from the user's instruction seems to be an attempt to define a new function
	// or modify the existing one in a way that doesn't align with the current structure.
	// Given the instruction "Explicitly ignore errors for Close and Fmt calls",
	// the `_ = conn.Close()` part has been applied above.
	// The rest of the provided snippet for a new function `func (b *networkBundle) performTLSCheck(...)`
	// is syntactically incomplete and refers to an undefined `networkBundle`.
	// Therefore, only the `_ = conn.Close()` change is applied to maintain a syntactically correct file.

	// Extract certificate details
	state := conn.ConnectionState()
	certs := state.PeerCertificates
	if len(certs) == 0 {
		return TLSCheckResponse{
			Connected: true,
			Error: &TLSError{
				Code:    "NO_CERTIFICATES",
				Message: "peer provided no certificates",
			},
		}
	}

	peerCert := certs[0]
	return TLSCheckResponse{
		Connected: true,
		NotAfter:  peerCert.NotAfter.Format(time.RFC3339),
		Issuer:    peerCert.Issuer.CommonName,
	}
}
