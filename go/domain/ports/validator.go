package ports

import "github.com/reglet-dev/reglet-sdk/go/domain/entities"

// CapabilityValidator validates capability configurations against schemas.
type CapabilityValidator interface {
	// Validate checks the manifest capabilities against registered schemas.
	Validate(manifest *entities.PluginManifest) (*entities.ValidationResult, error)
}
