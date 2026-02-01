package entities

import "encoding/json"

// Manifest contains complete plugin metadata for introspection.
// Used by both reglet and CLI tools to understand plugin capabilities.
type Manifest struct {
	// Identity
	Name        string `json:"name" yaml:"name"`
	Version     string `json:"version" yaml:"version"`
	Description string `json:"description" yaml:"description"`

	// Compatibility
	SDKVersion     string `json:"sdk_version" yaml:"sdk_version"`
	MinHostVersion string `json:"min_host_version,omitempty" yaml:"min_host_version,omitempty"`

	// Capabilities (http, dns, file, exec, etc.)
	Capabilities []Capability `json:"capabilities" yaml:"capabilities"`

	// Config schema (JSON Schema)
	ConfigSchema json.RawMessage `json:"config_schema" yaml:"config_schema"`

	// Registered services and operations
	Services map[string]ServiceManifest `json:"services" yaml:"services"`
}

// ServiceManifest describes a service and its operations.
type ServiceManifest struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Operations  []OperationManifest `json:"operations"`
}

// OperationManifest describes a single operation.
type OperationManifest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
