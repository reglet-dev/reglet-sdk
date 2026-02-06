// Package result provides a builder pattern for creating SDK results.
package result

import (
	abi "github.com/reglet-dev/reglet-abi"
	"github.com/reglet-dev/reglet-abi/hostfunc"
)

// Success creates a successful Result with the given message.
func Success(message string) *abi.Result {
	return &abi.Result{
		Status:  abi.ResultStatusSuccess,
		Message: message,
		Data:    make(map[string]any),
	}
}

// Failure creates a failure Result with the given message.
func Failure(message string) *abi.Result {
	return &abi.Result{
		Status:  abi.ResultStatusFailure,
		Message: message,
		Data:    make(map[string]any),
	}
}

// Errorer defines an interface for structured errors that provide codes and details.
type Errorer interface {
	Code() string
	Message() string
	Details() map[string]any
}

// Error creates an error Result from an error.
// It attempts to extract structured error details if the error implements the Errorer interface.
func Error(err error) *abi.Result {
	if err == nil {
		return Success("no error")
	}

	res := &abi.Result{
		Status:  abi.ResultStatusError,
		Message: err.Error(),
	}

	// helper to get details
	var details map[string]any
	code := "unknown"
	msg := err.Error()

	if e, ok := err.(Errorer); ok {
		code = e.Code()
		msg = e.Message()
		details = e.Details()
	}

	res.Error = &hostfunc.ErrorDetail{
		Type:    code,
		Message: msg,
		Details: details,
	}

	return res
}
