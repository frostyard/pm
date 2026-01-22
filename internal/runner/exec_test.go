package runner

import (
	"context"
	"testing"

	"github.com/frostyard/pm/internal/types"
)

func TestRunWithExternalError_Success(t *testing.T) {
	runner := &FakeRunner{
		StdoutResponse: "success output",
		StderrResponse: "",
		ErrResponse:    nil,
	}

	stdout, stderr, err := RunWithExternalError(
		context.Background(),
		runner,
		types.OperationSearch,
		"test-backend",
		"test-command",
		"arg1", "arg2",
	)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if stdout != "success output" {
		t.Errorf("Expected stdout='success output', got: %s", stdout)
	}
	if stderr != "" {
		t.Errorf("Expected empty stderr, got: %s", stderr)
	}
}

func TestRunWithExternalError_Failure(t *testing.T) {
	runner := &FakeRunner{
		StdoutResponse: "some output",
		StderrResponse: "error details",
		ErrResponse:    &fakeError{msg: "command failed"},
	}

	stdout, stderr, err := RunWithExternalError(
		context.Background(),
		runner,
		types.OperationInstall,
		"test-backend",
		"test-command",
		"arg1",
	)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Should return stdout/stderr even on error
	if stdout != "some output" {
		t.Errorf("Expected stdout='some output', got: %s", stdout)
	}
	if stderr != "error details" {
		t.Errorf("Expected stderr='error details', got: %s", stderr)
	}

	// Check that error is wrapped as ExternalFailureError
	if !types.IsExternalFailure(err) {
		t.Errorf("Expected ExternalFailureError, got: %T", err)
	}

	extErr, ok := err.(*types.ExternalFailureError)
	if !ok {
		t.Fatalf("Expected *ExternalFailureError, got: %T", err)
	}

	if extErr.Operation != types.OperationInstall {
		t.Errorf("Expected operation=Install, got: %s", extErr.Operation)
	}
	if extErr.Backend != "test-backend" {
		t.Errorf("Expected backend='test-backend', got: %s", extErr.Backend)
	}
	if extErr.Stdout != "some output" {
		t.Errorf("Expected stdout='some output', got: %s", extErr.Stdout)
	}
	if extErr.Stderr != "error details" {
		t.Errorf("Expected stderr='error details', got: %s", extErr.Stderr)
	}
}

func TestRunWithExternalError_Sanitization(t *testing.T) {
	// Create a very long output to test truncation
	longOutput := make([]byte, 1000)
	for i := range longOutput {
		longOutput[i] = 'x'
	}

	runner := &FakeRunner{
		StdoutResponse: string(longOutput),
		StderrResponse: string(longOutput),
		ErrResponse:    &fakeError{msg: "failed"},
	}

	_, _, err := RunWithExternalError(
		context.Background(),
		runner,
		types.OperationUpgradePackages,
		"test-backend",
		"test-command",
	)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	extErr := err.(*types.ExternalFailureError)

	// Check that output was truncated
	if len(extErr.Stdout) > 520 { // 500 + "... (truncated)" length tolerance
		t.Errorf("Expected stdout to be truncated, got length: %d", len(extErr.Stdout))
	}
	if len(extErr.Stderr) > 520 {
		t.Errorf("Expected stderr to be truncated, got length: %d", len(extErr.Stderr))
	}
}

func TestSanitize(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		maxLen    int
		wantLen   int
		wantTrunc bool
	}{
		{
			name:      "short string",
			input:     "short",
			wantLen:   5,
			wantTrunc: false,
		},
		{
			name:      "exactly at limit",
			input:     string(make([]byte, 500)),
			wantLen:   500,
			wantTrunc: false,
		},
		{
			name:      "over limit",
			input:     string(make([]byte, 600)),
			wantLen:   500 + len("... (truncated)"),
			wantTrunc: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitize(tt.input)
			if len(result) > 520 && tt.wantTrunc {
				// Allow some tolerance for truncation suffix
			} else if len(result) != tt.wantLen && !tt.wantTrunc {
				t.Errorf("Expected length %d, got %d", tt.wantLen, len(result))
			}
		})
	}
}

// fakeError is a simple error for testing.
type fakeError struct {
	msg string
}

func (e *fakeError) Error() string {
	return e.msg
}
