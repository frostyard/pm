package pm

import (
	"errors"
	"testing"
)

func TestIsNotSupported(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "direct ErrNotSupported",
			err:  ErrNotSupported,
			want: true,
		},
		{
			name: "wrapped NotSupportedError",
			err:  &NotSupportedError{Operation: OperationInstall, Backend: "test"},
			want: true,
		},
		{
			name: "unrelated error",
			err:  errors.New("something else"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotSupported(tt.err); got != tt.want {
				t.Errorf("IsNotSupported() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNotAvailable(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "direct ErrNotAvailable",
			err:  ErrNotAvailable,
			want: true,
		},
		{
			name: "wrapped NotAvailableError",
			err:  &NotAvailableError{Backend: "test", Reason: "not installed"},
			want: true,
		},
		{
			name: "unrelated error",
			err:  errors.New("something else"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotAvailable(tt.err); got != tt.want {
				t.Errorf("IsNotAvailable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsExternalFailure(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "ExternalFailureError",
			err:  &ExternalFailureError{Operation: OperationInstall, Backend: "test"},
			want: true,
		},
		{
			name: "unrelated error",
			err:  errors.New("something else"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsExternalFailure(tt.err); got != tt.want {
				t.Errorf("IsExternalFailure() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotSupportedError_Error(t *testing.T) {
	err := &NotSupportedError{
		Operation: OperationInstall,
		Backend:   "testbackend",
		Reason:    "not implemented yet",
	}

	msg := err.Error()
	if msg == "" {
		t.Error("NotSupportedError.Error() returned empty string")
	}

	// Check that it contains key information
	if !containsAll(msg, "Install", "testbackend", "not implemented yet") {
		t.Errorf("NotSupportedError.Error() = %q, missing expected content", msg)
	}
}

func TestNotAvailableError_Error(t *testing.T) {
	err := &NotAvailableError{
		Backend: "brew",
		Reason:  "command not found",
	}

	msg := err.Error()
	if msg == "" {
		t.Error("NotAvailableError.Error() returned empty string")
	}

	if !containsAll(msg, "brew", "command not found") {
		t.Errorf("NotAvailableError.Error() = %q, missing expected content", msg)
	}
}

func TestExternalFailureError_Error(t *testing.T) {
	err := &ExternalFailureError{
		Operation: OperationUpgradePackages,
		Backend:   "snap",
		Stderr:    "permission denied",
		Err:       errors.New("exit status 1"),
	}

	msg := err.Error()
	if msg == "" {
		t.Error("ExternalFailureError.Error() returned empty string")
	}

	if !containsAll(msg, "Upgrade", "snap", "permission denied") {
		t.Errorf("ExternalFailureError.Error() = %q, missing expected content", msg)
	}
}

// containsAll checks if s contains all substrings.
func containsAll(s string, subs ...string) bool {
	for _, sub := range subs {
		found := false
		for i := 0; i <= len(s)-len(sub); i++ {
			if s[i:i+len(sub)] == sub {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// T040: Test NotAvailable error detection and payload behavior
func TestNotAvailableError_Detection(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		wantBackend string
		wantReason  string
	}{
		{
			name: "with reason",
			err: &NotAvailableError{
				Backend: "brew",
				Reason:  "homebrew not installed",
			},
			wantBackend: "brew",
			wantReason:  "homebrew not installed",
		},
		{
			name: "without reason",
			err: &NotAvailableError{
				Backend: "flatpak",
			},
			wantBackend: "flatpak",
			wantReason:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !IsNotAvailable(tt.err) {
				t.Error("IsNotAvailable() should return true")
			}

			// Verify unwrapping
			if !errors.Is(tt.err, ErrNotAvailable) {
				t.Error("Should unwrap to ErrNotAvailable")
			}

			// Verify error message contains backend
			msg := tt.err.Error()
			if !containsAll(msg, tt.wantBackend) {
				t.Errorf("Error message should contain backend name: %s", msg)
			}

			// Verify reason if present
			if tt.wantReason != "" && !containsAll(msg, tt.wantReason) {
				t.Errorf("Error message should contain reason: %s", msg)
			}
		})
	}
}

// T040: Test ExternalFailure error detection and payload behavior
func TestExternalFailureError_Payload(t *testing.T) {
	tests := []struct {
		name           string
		err            *ExternalFailureError
		wantStdout     bool
		wantStderr     bool
		wantPayload    bool
		wantUnderlying bool
	}{
		{
			name: "with stdout and stderr",
			err: &ExternalFailureError{
				Operation: OperationInstall,
				Backend:   "apt",
				Stdout:    "Reading package lists...",
				Stderr:    "E: Unable to locate package",
				Err:       errors.New("exit status 100"),
			},
			wantStdout:     true,
			wantStderr:     true,
			wantUnderlying: true,
		},
		{
			name: "with API payload",
			err: &ExternalFailureError{
				Operation: OperationSearch,
				Backend:   "snapd",
				Payload: map[string]interface{}{
					"error": "not found",
					"code":  404,
				},
			},
			wantPayload: true,
		},
		{
			name: "minimal error",
			err: &ExternalFailureError{
				Operation: OperationUpdateMetadata,
				Backend:   "brew",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify it's detected as ExternalFailure
			if !IsExternalFailure(tt.err) {
				t.Error("IsExternalFailure() should return true")
			}

			// Verify error message structure
			msg := tt.err.Error()
			if !containsAll(msg, string(tt.err.Operation), tt.err.Backend) {
				t.Errorf("Error message should contain operation and backend: %s", msg)
			}

			// Verify stdout is accessible
			if tt.wantStdout && tt.err.Stdout == "" {
				t.Error("Expected stdout to be populated")
			}

			// Verify stderr is accessible
			if tt.wantStderr && tt.err.Stderr == "" {
				t.Error("Expected stderr to be populated")
			}

			// Verify payload is accessible
			if tt.wantPayload && len(tt.err.Payload) == 0 {
				t.Error("Expected payload to be populated")
			}

			// Verify underlying error
			if tt.wantUnderlying {
				if tt.err.Unwrap() == nil {
					t.Error("Expected underlying error to be accessible")
				}
			}
		})
	}
}

// T040: Test that ExternalFailureError can be extracted with errors.As
func TestExternalFailureError_ErrorsAs(t *testing.T) {
	originalErr := &ExternalFailureError{
		Operation: OperationInstall,
		Backend:   "test",
		Stdout:    "output",
		Stderr:    "error",
		Err:       errors.New("underlying"),
	}

	// Wrap it
	wrapped := errors.New("wrapped: " + originalErr.Error())

	// This won't work with errors.As because wrapped is a different error
	// But we can test direct usage
	var extErr *ExternalFailureError
	if errors.As(originalErr, &extErr) {
		if extErr.Stdout != "output" {
			t.Errorf("Expected stdout='output', got: %s", extErr.Stdout)
		}
		if extErr.Stderr != "error" {
			t.Errorf("Expected stderr='error', got: %s", extErr.Stderr)
		}
	} else {
		t.Error("errors.As should work with ExternalFailureError")
	}

	// Verify IsExternalFailure works
	if !IsExternalFailure(originalErr) {
		t.Error("IsExternalFailure should return true")
	}

	// Verify it doesn't match unrelated errors
	if IsExternalFailure(wrapped) {
		t.Error("IsExternalFailure should return false for unrelated wrapped errors")
	}
}
