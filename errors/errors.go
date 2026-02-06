// Package errors provides structured error types for plugins.
package errors

import "fmt"

// PluginError represents a structured error from a plugin.
type PluginError struct {
	ErrCode    string         `json:"code"`
	ErrMessage string         `json:"message"`
	ErrDetails map[string]any `json:"details,omitempty"`
}

func (e *PluginError) Error() string {
	return fmt.Sprintf("[%s] %s", e.ErrCode, e.ErrMessage)
}

// Code returns the error code.
func (e *PluginError) Code() string {
	return e.ErrCode
}

// Message returns the error message.
func (e *PluginError) Message() string {
	return e.ErrMessage
}

// Details returns the error details.
func (e *PluginError) Details() map[string]any {
	return e.ErrDetails
}

// ConfigError creates a config validation error.
func ConfigError(field, reason string) *PluginError {
	return &PluginError{
		ErrCode:    "config_invalid",
		ErrMessage: fmt.Sprintf("field %s: %s", field, reason),
	}
}

// NetworkError creates a network error.
func NetworkError(host, port, reason string) *PluginError {
	return &PluginError{
		ErrCode:    "network_error",
		ErrMessage: fmt.Sprintf("connect %s:%s: %s", host, port, reason),
	}
}
