// Package net provides network operations for Reglet WASM plugins.
package net //nolint:revive // intentional name to mirror Go stdlib net package

import (
	"github.com/reglet-dev/reglet-sdk/go/wireformat"
)

// Re-export wire format types from shared wireformat package
// This file has no build tags so tests can use these types
type (
	ContextWireFormat = wireformat.ContextWireFormat
	DNSRequestWire    = wireformat.DNSRequestWire
	DNSResponseWire   = wireformat.DNSResponseWire
	TCPRequestWire    = wireformat.TCPRequestWire
	TCPResponseWire   = wireformat.TCPResponseWire
	SMTPRequestWire   = wireformat.SMTPRequestWire
	SMTPResponseWire  = wireformat.SMTPResponseWire
)
