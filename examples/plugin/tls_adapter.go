//go:build wasip1

package main

import (
	"context"

	"github.com/reglet-dev/reglet-sdk/go/application/plugin"
)

// Define the host function signature for our custom TLS check.
//
//go:wasmimport reglet_host tls_check
func host_tls_check(requestPacked uint64) uint64

// TLSCheckRequest matches the host's expectation.
type TLSCheckRequest struct {
	Host      string `json:"host"`
	Port      int    `json:"port"`
	TimeoutMs int    `json:"timeout_ms"`
}

// TLSCheckResponse matches the host's response.
type TLSCheckResponse struct {
	Connected bool      `json:"connected"`
	NotAfter  string    `json:"not_after,omitempty"`
	Issuer    string    `json:"issuer,omitempty"`
	Error     *TLSError `json:"error,omitempty"`
}

type TLSError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// PerformTLSCheck calls the host function to verify a certificate.
func PerformTLSCheck(ctx context.Context, host string, port int) (TLSCheckResponse, error) {
	req := TLSCheckRequest{
		Host:      host,
		Port:      port,
		TimeoutMs: 5000,
	}

	return plugin.CallHost[TLSCheckRequest, TLSCheckResponse](host_tls_check, req)
}
