# Progress Reporting

A Go module for structured progress reporting with hierarchical actions, tasks, and steps.

## Features

- **Hierarchical Progress**: Track actions → tasks → steps
- **Thread-Safe**: Built-in concurrency support
- **Message Severity**: Info, Warning, and Error levels
- **Flexible Reporting**: Implement custom reporters for any output format

## Installation

```bash
go get github.com/frostyard/pm/progress
```

## Usage

```go
import "github.com/frostyard/pm/progress"

// Implement a reporter
type MyReporter struct{}

func (r *MyReporter) OnAction(action progress.ProgressAction) {
    if action.EndedAt.IsZero() {
        fmt.Printf("→ %s\n", action.Name)
    }
}

func (r *MyReporter) OnTask(task progress.ProgressTask) {
    if task.EndedAt.IsZero() {
        fmt.Printf("  • %s\n", task.Name)
    }
}

func (r *MyReporter) OnStep(step progress.ProgressStep) {
    if step.EndedAt.IsZero() {
        fmt.Printf("    - %s\n", step.Name)
    }
}

func (r *MyReporter) OnMessage(msg progress.ProgressMessage) {
    fmt.Printf("    %s: %s\n", msg.Severity, msg.Text)
}

// Use the helper
reporter := &MyReporter{}
helper := progress.NewProgressHelper(reporter, nil)

helper.BeginAction("Processing")
helper.BeginTask("Loading data")
helper.Info("Found 100 records")
helper.EndTask()
helper.EndAction()
```

## License

See the main repository LICENSE file.
