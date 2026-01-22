package runner

import "context"

// Runner abstracts command execution for CLI-based backends.
// This enables deterministic unit testing by injecting fake/mock implementations.
type Runner interface {
	// Run executes a command and returns stdout, stderr, and any error.
	Run(ctx context.Context, name string, args ...string) (stdout, stderr string, err error)
}
