package sdknet

import (
	"context"
	"fmt"
	"time"

	"github.com/reglet-dev/reglet-sdk/go/application/config"
	"github.com/reglet-dev/reglet-sdk/go/domain/entities"
	"github.com/reglet-dev/reglet-sdk/go/hostfuncs"
)

// RunHTTPCheck performs an HTTP request check.
// It parses configuration, executes the HTTP request, and returns a structured Result.
//
// Expected config fields:
//   - url (string, required): Target URL
//   - method (string, optional): HTTP method (default: GET)
//   - headers (map[string]string, optional): Request headers
//   - body (string, optional): Request body
//   - timeout_ms (int, optional): Request timeout in milliseconds (default: 30000)
//   - expected_status (int, optional): Expected HTTP status code for validation
//   - follow_redirects (bool, optional): Whether to follow redirects (default: true)
//   - max_redirects (int, optional): Maximum redirects to follow (default: 10)
//
// Returns a Result with:
//   - Status: "success" if request succeeded and matches expectations, "failure" if status mismatch, "error" if request failed
//   - Data: map containing "status_code", "headers", "body", "latency_ms", "body_truncated"
//   - Error: structured error details if request failed
func RunHTTPCheck(ctx context.Context, cfg config.Config) (entities.Result, error) {
	// Parse required fields
	url, err := config.MustGetString(cfg, "url")
	if err != nil {
		return entities.ResultError(entities.NewErrorDetail("config", err.Error()).WithCode("MISSING_URL")), nil
	}

	// Parse optional fields
	method := config.GetStringDefault(cfg, "method", "GET")
	timeoutMs := config.GetIntDefault(cfg, "timeout_ms", 30000)
	expectedStatus, hasExpectedStatus := config.GetInt(cfg, "expected_status")

	// Parse headers if provided
	var headers map[string]string
	if headersRaw, ok := cfg["headers"].(map[string]interface{}); ok {
		headers = make(map[string]string)
		for k, v := range headersRaw {
			if vStr, ok := v.(string); ok {
				headers[k] = vStr
			}
		}
	}

	// Parse body if provided
	body := config.GetStringDefault(cfg, "body", "")

	// Create request
	req := hostfuncs.HTTPRequest{
		Method:  method,
		URL:     url,
		Headers: headers,
		Timeout: timeoutMs,
	}

	if body != "" {
		req.Body = []byte(body)
	}

	// Handle redirect configuration
	if followRedirects, ok := config.GetBool(cfg, "follow_redirects"); ok {
		req.FollowRedirects = &followRedirects
	}
	if maxRedirects, ok := config.GetInt(cfg, "max_redirects"); ok {
		req.MaxRedirects = maxRedirects
	}

	// Execute HTTP request
	start := time.Now()
	resp := hostfuncs.PerformHTTPRequest(ctx, req)
	metadata := entities.NewRunMetadata(start, time.Now())

	// Check for request errors
	if resp.Error != nil {
		errDetail := entities.NewErrorDetail("network", resp.Error.Message).WithCode(resp.Error.Code)
		return entities.ResultError(errDetail).WithMetadata(metadata), errDetail
	}

	// Build result data
	resultData := map[string]any{
		"status_code":    resp.StatusCode,
		"latency_ms":     resp.LatencyMs,
		"body_truncated": resp.BodyTruncated,
	}

	if len(resp.Headers) > 0 {
		resultData["headers"] = resp.Headers
	}

	if len(resp.Body) > 0 {
		resultData["body"] = string(resp.Body)
		resultData["body_length"] = len(resp.Body)
	}

	// Check expected status if specified
	if hasExpectedStatus && resp.StatusCode != expectedStatus {
		message := fmt.Sprintf("HTTP status mismatch: expected %d, got %d", expectedStatus, resp.StatusCode)
		resultData["expected_status"] = expectedStatus
		resultData["actual_status"] = resp.StatusCode
		return entities.ResultFailure(message, resultData).WithMetadata(metadata), nil
	}

	// Success
	message := fmt.Sprintf("HTTP %s request successful: %d", method, resp.StatusCode)
	return entities.ResultSuccess(message, resultData).WithMetadata(metadata), nil
}
