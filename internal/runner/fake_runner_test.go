package runner

import "context"

// FakeRunner is a deterministic fake runner for unit tests.
type FakeRunner struct {
	// StdoutResponse is the stdout to return.
	StdoutResponse string

	// StderrResponse is the stderr to return.
	StderrResponse string

	// ErrResponse is the error to return.
	ErrResponse error

	// LastCommand captures the last command executed for assertions.
	LastCommand string

	// LastArgs captures the last args for assertions.
	LastArgs []string
}

// Run executes the fake command.
func (f *FakeRunner) Run(ctx context.Context, name string, args ...string) (string, string, error) {
	f.LastCommand = name
	f.LastArgs = args
	return f.StdoutResponse, f.StderrResponse, f.ErrResponse
}
