package entities

// ValidationResult represents the outcome of a capability validation.
type ValidationResult struct {
	Valid  bool
	Errors []ValidationError
}

// ValidationError represents a specific validation error.
type ValidationError struct {
	Field   string
	Message string
}
