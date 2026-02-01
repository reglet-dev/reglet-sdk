package ports

import "github.com/reglet-dev/reglet-sdk/go/domain/entities"

// ManifestParser parses raw YAML bytes into a PluginManifest.
type ManifestParser interface {
	// Parse unmarshals YAML bytes into a PluginManifest struct.
	Parse(data []byte) (*entities.Manifest, error)
}
