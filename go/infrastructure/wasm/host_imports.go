//go:build wasip1

// Package wasm provides infrastructure adapters that interface with the WASM host environment.
package wasm

// Define the host function signature for HTTP requests.
//
//go:wasmimport reglet_host http_request
func host_http_request(requestPacked uint64) uint64

// Define the host function signature for DNS lookups.
//
//go:wasmimport reglet_host dns_lookup
func host_dns_lookup(requestPacked uint64) uint64
