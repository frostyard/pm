package runner

import (
	"context"
	"os/exec"
	"strings"

	"github.com/frostyard/pm/internal/types"
)

// realRunner implements Runner using os/exec.
type realRunner struct{}

// NewRealRunner creates a Runner that executes real commands using os/exec.
func NewRealRunner() Runner {
	return &realRunner{}
}

// Run executes a command using os/exec and returns stdout, stderr, and error.
func (r *realRunner) Run(ctx context.Context, name string, args ...string) (string, string, error) {
	cmd := exec.CommandContext(ctx, name, args...)

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// RunWithExternalError executes a command and wraps failures in ExternalFailureError.
// This provides structured error reporting with captured stdout/stderr for CLI-based backends.
//
// Parameters:
//   - ctx: Context for cancellation
//   - runner: Runner implementation (real or fake for testing)
//   - operation: The operation being performed (for error context)
//   - backend: The backend name (for error context)
//   - name: Command name to execute
//   - args: Command arguments
//
// Returns:
//   - stdout: Captured standard output
//   - stderr: Captured standard error
//   - error: nil on success, ExternalFailureError on failure
func RunWithExternalError(
	ctx context.Context,
	runner Runner,
	operation types.Operation,
	backend string,
	name string,
	args ...string,
) (stdout, stderr string, err error) {
	stdout, stderr, err = runner.Run(ctx, name, args...)

	if err != nil {
		return stdout, stderr, &types.ExternalFailureError{
			Operation: operation,
			Backend:   backend,
			Stdout:    sanitize(stdout),
			Stderr:    sanitize(stderr),
			Err:       err,
		}
	}

	return stdout, stderr, nil
}

// sanitize removes sensitive information from command output.
// For now, this is a simple length limiter to prevent huge error messages.
// In production, you might want to filter passwords, tokens, etc.
func sanitize(s string) string {
	const maxLen = 500
	if len(s) > maxLen {
		return s[:maxLen] + "... (truncated)"
	}
	return s
}
