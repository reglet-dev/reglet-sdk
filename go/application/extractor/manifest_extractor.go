package extractor

import (
	"fmt"

	"github.com/reglet-dev/reglet-sdk/go/domain/entities"
	"github.com/reglet-dev/reglet-sdk/go/domain/ports"
)

// ManifestExtractor extracts capabilities from a plugin manifest.
type ManifestExtractor struct {
	parser   ports.ManifestParser
	renderer ports.TemplateEngine
	manifest []byte
}

// ManifestExtractorOption configures the ManifestExtractor.
type ManifestExtractorOption func(*ManifestExtractor)

// WithParser sets the manifest parser.
func WithParser(p ports.ManifestParser) ManifestExtractorOption {
	return func(e *ManifestExtractor) {
		e.parser = p
	}
}

// WithTemplateEngine sets the template engine.
func WithTemplateEngine(t ports.TemplateEngine) ManifestExtractorOption {
	return func(e *ManifestExtractor) {
		e.renderer = t
	}
}

// NewManifestExtractor creates a new ManifestExtractor for the given manifest.
func NewManifestExtractor(manifest []byte, opts ...ManifestExtractorOption) *ManifestExtractor {
	e := &ManifestExtractor{
		manifest: manifest,
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// Extract analyzes the manifest and returns the required capabilities.
func (e *ManifestExtractor) Extract(config map[string]interface{}) (*entities.GrantSet, error) {
	if e.parser == nil {
		return nil, fmt.Errorf("manifest parser is required")
	}

	data := e.manifest
	if e.renderer != nil {
		var err error
		data, err = e.renderer.Render(data, config)
		if err != nil {
			return nil, fmt.Errorf("failed to render manifest: %w", err)
		}
	}

	manifest, err := e.parser.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	if manifest.Capabilities == nil {
		return &entities.GrantSet{}, nil
	}

	// TODO: look into this code
	// Convert []Capability to *GrantSet
	gs := &entities.GrantSet{}
	for _, cap := range manifest.Capabilities {
		switch cap.Category {
		case "network", "http": // Handle http as network for now
			// Simple parsing: resource as host. If host:port, split.
			host := cap.Resource
			port := "*" // default
			// This is a naive implementation, real parsing should be more robust
			// or delegate to specific parser.
			if gs.Network == nil {
				gs.Network = &entities.NetworkCapability{}
			}
			gs.Network.Rules = append(gs.Network.Rules, entities.NetworkRule{
				Hosts: []string{host},
				Ports: []string{port},
			})
		case "fs":
			if gs.FS == nil {
				gs.FS = &entities.FileSystemCapability{}
			}
			read := []string{}
			write := []string{}
			if cap.Action == "write" {
				write = append(write, cap.Resource)
			} else {
				read = append(read, cap.Resource)
			}
			gs.FS.Rules = append(gs.FS.Rules, entities.FileSystemRule{
				Read:  read,
				Write: write,
			})
		case "exec":
			if gs.Exec == nil {
				gs.Exec = &entities.ExecCapability{}
			}
			gs.Exec.Commands = append(gs.Exec.Commands, cap.Resource)
		case "env":
			if gs.Env == nil {
				gs.Env = &entities.EnvironmentCapability{}
			}
			gs.Env.Variables = append(gs.Env.Variables, cap.Resource)
		}
	}

	return gs, nil
}
