package pm

import (
	"errors"
	"fmt"
)

var (
	// ErrNotSupported is returned when an operation is not supported by the backend.
	ErrNotSupported = errors.New("operation not supported")

	// ErrNotAvailable is returned when a backend is not available (not installed/reachable).
	ErrNotAvailable = errors.New("backend not available")
)

// NotSupportedError wraps ErrNotSupported with additional context.
type NotSupportedError struct {
	Operation Operation
	Backend   string
	Reason    string
}

func (e *NotSupportedError) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("%s: %s operation not supported by %s: %s", ErrNotSupported, e.Operation, e.Backend, e.Reason)
	}
	return fmt.Sprintf("%s: %s operation not supported by %s", ErrNotSupported, e.Operation, e.Backend)
}

func (e *NotSupportedError) Unwrap() error {
	return ErrNotSupported
}

// IsNotSupported checks if an error is a NotSupported error.
func IsNotSupported(err error) bool {
	return errors.Is(err, ErrNotSupported)
}

// NotAvailableError wraps ErrNotAvailable with additional context.
type NotAvailableError struct {
	Backend string
	Reason  string
}

func (e *NotAvailableError) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("%s: %s: %s", ErrNotAvailable, e.Backend, e.Reason)
	}
	return fmt.Sprintf("%s: %s", ErrNotAvailable, e.Backend)
}

func (e *NotAvailableError) Unwrap() error {
	return ErrNotAvailable
}

// IsNotAvailable checks if an error is a NotAvailable error.
func IsNotAvailable(err error) bool {
	return errors.Is(err, ErrNotAvailable)
}

// ExternalFailureError represents a failure from an external command or API.
type ExternalFailureError struct {
	Operation Operation
	Backend   string
	// Stdout captured from command (if applicable, sanitized).
	Stdout string
	// Stderr captured from command (if applicable, sanitized).
	Stderr string
	// Payload is structured error data from an API (if applicable).
	Payload map[string]interface{}
	// Underlying error.
	Err error
}

func (e *ExternalFailureError) Error() string {
	msg := fmt.Sprintf("external failure: %s operation on %s", e.Operation, e.Backend)
	if e.Err != nil {
		msg = fmt.Sprintf("%s: %v", msg, e.Err)
	}
	if e.Stderr != "" {
		msg = fmt.Sprintf("%s (stderr: %s)", msg, e.Stderr)
	}
	return msg
}

func (e *ExternalFailureError) Unwrap() error {
	return e.Err
}

// IsExternalFailure checks if an error is an ExternalFailure error.
func IsExternalFailure(err error) bool {
	var extErr *ExternalFailureError
	return errors.As(err, &extErr)
}
