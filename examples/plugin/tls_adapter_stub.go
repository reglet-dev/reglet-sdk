//go:build !wasip1

package main

import (
	"context"
	"fmt"
)

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

func PerformTLSCheck(ctx context.Context, host string, port int) (TLSCheckResponse, error) {
	return TLSCheckResponse{}, fmt.Errorf("tls_check not available on this platform")
}
