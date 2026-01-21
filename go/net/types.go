// Package sdknet provides network operations for Reglet WASM plugins.
package sdknet

import (
	"github.com/reglet-dev/reglet-sdk/go/domain/entities"
)

// Re-export wire format types from entities package for backward compatibility.
// These types serve as both domain entities and JSON wire protocol structures.
type (
	ContextWireFormat = entities.ContextWire
	DNSRequestWire    = entities.DNSRequest
	DNSResponseWire   = entities.DNSResponse
	MXRecordWire      = entities.MXRecord
	TCPRequestWire    = entities.TCPRequest
	TCPResponseWire   = entities.TCPResponse
	SMTPRequestWire   = entities.SMTPRequest
	SMTPResponseWire  = entities.SMTPResponse
)
