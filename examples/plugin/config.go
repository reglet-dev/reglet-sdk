package main

import (
	"github.com/reglet-dev/reglet-sdk/go/application/config"
)

// PluginConfig represents the configuration for the example compliance plugin.
type PluginConfig struct {
	TargetHost string `json:"target_host" jsonschema:"default=example.com"`
	TargetPort int    `json:"target_port" jsonschema:"default=443"`
	MinDays    int    `json:"min_days" jsonschema:"default=30"`
}

// LoadConfig extracts and validates the configuration from the input map.
func LoadConfig(cfg map[string]any) (PluginConfig, error) {
	return PluginConfig{
		TargetHost: config.GetStringDefault(cfg, "target_host", "example.com"),
		TargetPort: config.GetIntDefault(cfg, "target_port", 443),
		MinDays:    config.GetIntDefault(cfg, "min_days", 30),
	}, nil
}
