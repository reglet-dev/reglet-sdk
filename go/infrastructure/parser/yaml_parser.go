package parser

import (
	"github.com/reglet-dev/reglet-sdk/go/domain/entities"
	"github.com/reglet-dev/reglet-sdk/go/domain/ports"
	"gopkg.in/yaml.v3"
)

// YamlManifestParser implements ManifestParser for YAML.
type YamlManifestParser struct{}

// NewYamlManifestParser creates a new YamlManifestParser.
func NewYamlManifestParser() ports.ManifestParser {
	return &YamlManifestParser{}
}

// Parse unmarshals YAML bytes into a PluginManifest struct.
func (p *YamlManifestParser) Parse(data []byte) (*entities.PluginManifest, error) {
	var manifest entities.PluginManifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}
	return &manifest, nil
}
