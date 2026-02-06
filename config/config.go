// Package config provides configuration validation helpers.
package config

import "fmt"

// RequireString returns the string value of a key in the config map.
// It returns an error if the key is missing, not a string, or empty.
func RequireString(cfg map[string]any, key string) (string, error) {
	val, ok := cfg[key]
	if !ok {
		return "", fmt.Errorf("missing required field: %s", key)
	}
	str, ok := val.(string)
	if !ok || str == "" {
		return "", fmt.Errorf("field %s must be non-empty string", key)
	}
	return str, nil
}

// OptionalString returns the string value of a key in the config map.
// It returns the default value if the key is missing, not a string, or empty.
func OptionalString(cfg map[string]any, key, defaultVal string) string {
	if val, ok := cfg[key].(string); ok && val != "" {
		return val
	}
	return defaultVal
}

// OptionalInt returns the int value of a key in the config map.
// It returns the default value if the key is missing or not a number.
// Handles both float64 (JSON default) and int types.
func OptionalInt(cfg map[string]any, key string, defaultVal int) int {
	if val, ok := cfg[key].(float64); ok {
		return int(val)
	}
	if val, ok := cfg[key].(int); ok {
		return val
	}
	return defaultVal
}
