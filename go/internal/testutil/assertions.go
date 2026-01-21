// Package testutil provides common test utilities and assertions for SDK tests
package testutil

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertNoError is a convenience wrapper for require.NoError with a descriptive message
func AssertNoError(t *testing.T, err error, msgAndArgs ...interface{}) {
	t.Helper()
	require.NoError(t, err, msgAndArgs...)
}

// AssertError is a convenience wrapper for require.Error with a descriptive message
func AssertError(t *testing.T, err error, msgAndArgs ...interface{}) {
	t.Helper()
	require.Error(t, err, msgAndArgs...)
}

// AssertEqual is a convenience wrapper for assert.Equal
func AssertEqual(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	assert.Equal(t, expected, actual, msgAndArgs...)
}

// AssertNotEqual is a convenience wrapper for assert.NotEqual
func AssertNotEqual(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	assert.NotEqual(t, expected, actual, msgAndArgs...)
}

// AssertTrue is a convenience wrapper for assert.True
func AssertTrue(t *testing.T, value bool, msgAndArgs ...interface{}) {
	t.Helper()
	assert.True(t, value, msgAndArgs...)
}

// AssertFalse is a convenience wrapper for assert.False
func AssertFalse(t *testing.T, value bool, msgAndArgs ...interface{}) {
	t.Helper()
	assert.False(t, value, msgAndArgs...)
}

// AssertJSONEqual compares two JSON strings for equality, ignoring formatting
func AssertJSONEqual(t *testing.T, expected, actual string, msgAndArgs ...interface{}) {
	t.Helper()

	var expectedJSON, actualJSON interface{}
	require.NoError(t, json.Unmarshal([]byte(expected), &expectedJSON), "expected JSON is invalid")
	require.NoError(t, json.Unmarshal([]byte(actual), &actualJSON), "actual JSON is invalid")

	assert.Equal(t, expectedJSON, actualJSON, msgAndArgs...)
}

// AssertDurationWithin asserts that a duration is within a tolerance of an expected value
func AssertDurationWithin(t *testing.T, expected, actual, tolerance time.Duration, msgAndArgs ...interface{}) {
	t.Helper()

	diff := expected - actual
	if diff < 0 {
		diff = -diff
	}

	assert.LessOrEqual(t, diff, tolerance, msgAndArgs...)
}

// AssertMapContains asserts that a map contains all expected key-value pairs
func AssertMapContains(t *testing.T, expectedMap, actualMap map[string]interface{}, msgAndArgs ...interface{}) {
	t.Helper()

	for key, expectedValue := range expectedMap {
		actualValue, ok := actualMap[key]
		assert.True(t, ok, "map should contain key %q", key)
		assert.Equal(t, expectedValue, actualValue, msgAndArgs...)
	}
}

// AssertPanics asserts that the function panics
func AssertPanics(t *testing.T, f func(), msgAndArgs ...interface{}) {
	t.Helper()
	assert.Panics(t, f, msgAndArgs...)
}

// AssertNotPanics asserts that the function does not panic
func AssertNotPanics(t *testing.T, f func(), msgAndArgs ...interface{}) {
	t.Helper()
	assert.NotPanics(t, f, msgAndArgs...)
}

// RequireNoError is a convenience wrapper for require.NoError
func RequireNoError(t *testing.T, err error, msgAndArgs ...interface{}) {
	t.Helper()
	require.NoError(t, err, msgAndArgs...)
}

// RequireError is a convenience wrapper for require.Error
func RequireError(t *testing.T, err error, msgAndArgs ...interface{}) {
	t.Helper()
	require.Error(t, err, msgAndArgs...)
}

// RequireEqual is a convenience wrapper for require.Equal
func RequireEqual(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	require.Equal(t, expected, actual, msgAndArgs...)
}

// RequireNotNil is a convenience wrapper for require.NotNil
func RequireNotNil(t *testing.T, object interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	require.NotNil(t, object, msgAndArgs...)
}

// RequireNil is a convenience wrapper for require.Nil
func RequireNil(t *testing.T, object interface{}, msgAndArgs ...interface{}) {
	t.Helper()
	require.Nil(t, object, msgAndArgs...)
}
