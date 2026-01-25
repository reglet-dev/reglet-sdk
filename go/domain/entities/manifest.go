package entities

// PluginManifest represents the root configuration of a plugin.
type PluginManifest struct {
	Name         string    `json:"name" yaml:"name"`
	Version      string    `json:"version" yaml:"version"`
	Description  string    `json:"description,omitempty" yaml:"description,omitempty"`
	Capabilities *GrantSet `json:"capabilities,omitempty" yaml:"capabilities,omitempty"`
}
