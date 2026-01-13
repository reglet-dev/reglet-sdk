// Package wireformat defines the JSON wire format structures for communication
// between the WASM host and guest (plugins). These types must remain stable
// and backward compatible as they define the ABI contract.
package wireformat

import (
	"fmt"
	"time"
)

// ContextWireFormat is the JSON wire format for context.Context propagation.
type ContextWireFormat struct {
	Deadline  *time.Time `json:"deadline,omitempty"`
	RequestID string     `json:"request_id,omitempty"`
	TimeoutMs int64      `json:"timeout_ms,omitempty"`
	Canceled  bool       `json:"Canceled,omitempty"`
}

// DNSRequestWire is the JSON wire format for a DNS lookup request from Guest to Host.
type DNSRequestWire struct {
	Hostname   string            `json:"hostname"`
	Type       string            `json:"type"`
	Nameserver string            `json:"nameserver,omitempty"`
	Context    ContextWireFormat `json:"context"`
}

// DNSResponseWire is the JSON wire format for a DNS lookup response from Host to Guest.
type DNSResponseWire struct {
	Error     *ErrorDetail   `json:"error,omitempty"`
	Records   []string       `json:"records,omitempty"`
	MXRecords []MXRecordWire `json:"mx_records,omitempty"`
}

// MXRecordWire represents a single MX record.
type MXRecordWire struct {
	Host string `json:"host"`
	Pref uint16 `json:"pref"`
}

// HTTPRequestWire is the JSON wire format for an HTTP request from Guest to Host.
type HTTPRequestWire struct {
	Headers map[string][]string `json:"headers,omitempty"`
	Method  string              `json:"method"`
	URL     string              `json:"url"`
	Body    string              `json:"body,omitempty"`
	Context ContextWireFormat   `json:"context"`
}

// HTTPResponseWire is the JSON wire format for an HTTP response from Host to Guest.
type HTTPResponseWire struct {
	Headers       map[string][]string `json:"headers,omitempty"`
	Error         *ErrorDetail        `json:"error,omitempty"`
	Body          string              `json:"body,omitempty"`
	StatusCode    int                 `json:"status_code"`
	BodyTruncated bool                `json:"body_truncated,omitempty"`
}

// TCPRequestWire is the JSON wire format for a TCP connection request from Guest to Host.
type TCPRequestWire struct {
	Host      string            `json:"host"`
	Port      string            `json:"port"`
	Context   ContextWireFormat `json:"context"`
	TimeoutMs int               `json:"timeout_ms,omitempty"`
	TLS       bool              `json:"tls"`
}

// TCPResponseWire is the JSON wire format for a TCP connection response from Host to Guest.
type TCPResponseWire struct {
	TLSCertNotAfter *time.Time   `json:"tls_cert_not_after,omitempty"`
	Error           *ErrorDetail `json:"error,omitempty"`
	TLSVersion      string       `json:"tls_version,omitempty"`
	LocalAddr       string       `json:"local_addr,omitempty"`
	TLSCipherSuite  string       `json:"tls_cipher_suite,omitempty"`
	TLSServerName   string       `json:"tls_server_name,omitempty"`
	TLSCertSubject  string       `json:"tls_cert_subject,omitempty"`
	TLSCertIssuer   string       `json:"tls_cert_issuer,omitempty"`
	RemoteAddr      string       `json:"remote_addr,omitempty"`
	Address         string       `json:"address,omitempty"`
	ResponseTimeMs  int64        `json:"response_time_ms,omitempty"`
	TLS             bool         `json:"tls,omitempty"`
	Connected       bool         `json:"connected"`
}

// SMTPRequestWire is the JSON wire format for an SMTP connection request from Guest to Host.
type SMTPRequestWire struct {
	Host      string            `json:"host"`
	Port      string            `json:"port"`
	Context   ContextWireFormat `json:"context"`
	TimeoutMs int               `json:"timeout_ms,omitempty"`
	TLS       bool              `json:"tls"`
	StartTLS  bool              `json:"starttls"`
}

// SMTPResponseWire is the JSON wire format for an SMTP connection response from Host to Guest.
type SMTPResponseWire struct {
	Error          *ErrorDetail `json:"error,omitempty"`
	Address        string       `json:"address,omitempty"`
	Banner         string       `json:"banner,omitempty"`
	TLSVersion     string       `json:"tls_version,omitempty"`
	TLSCipherSuite string       `json:"tls_cipher_suite,omitempty"`
	TLSServerName  string       `json:"tls_server_name,omitempty"`
	ResponseTimeMs int64        `json:"response_time_ms,omitempty"`
	Connected      bool         `json:"connected"`
	TLS            bool         `json:"tls,omitempty"`
}

// ExecRequestWire is the JSON wire format for an exec request from Guest to Host.
type ExecRequestWire struct {
	Args    []string          `json:"args"`
	Env     []string          `json:"env,omitempty"`
	Command string            `json:"command"`
	Dir     string            `json:"dir,omitempty"`
	Context ContextWireFormat `json:"context"`
}

// ExecResponseWire is the JSON wire format for an exec response from Host to Guest.
type ExecResponseWire struct {
	Error      *ErrorDetail `json:"error,omitempty"`
	Stdout     string       `json:"stdout"`
	Stderr     string       `json:"stderr"`
	ExitCode   int          `json:"exit_code"`
	DurationMs int64        `json:"duration_ms,omitempty"`
	IsTimeout  bool         `json:"is_timeout,omitempty"`
}

// ErrorDetail provides structured error information, consistent across host and SDK.
// Error Types: "network", "timeout", "config", "panic", "capability", "validation", "internal"
type ErrorDetail struct {
	Wrapped    *ErrorDetail `json:"wrapped,omitempty"`
	Message    string       `json:"message"`
	Type       string       `json:"type"`
	Code       string       `json:"code"`
	Stack      []byte       `json:"stack,omitempty"`
	IsTimeout  bool         `json:"is_timeout,omitempty"`
	IsNotFound bool         `json:"is_not_found,omitempty"`
}

// Error implements the error interface for ErrorDetail.
func (e *ErrorDetail) Error() string {
	if e == nil {
		return ""
	}
	msg := e.Message
	if e.Type != "" && e.Type != "internal" {
		msg = fmt.Sprintf("%s: %s", e.Type, msg)
	}
	if e.Code != "" {
		msg = fmt.Sprintf("%s [%s]", msg, e.Code)
	}
	if e.Wrapped != nil {
		msg = fmt.Sprintf("%s: %v", msg, e.Wrapped.Error())
	}
	return msg
}
