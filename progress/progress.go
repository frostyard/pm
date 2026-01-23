package progress

import "time"

// Severity represents the severity level of a progress message.
type Severity string

const (
	// SeverityInfo represents an informational message.
	SeverityInfo Severity = "Informational"

	// SeverityWarning represents a warning message (does not fail the operation).
	SeverityWarning Severity = "Warning"

	// SeverityError represents an error message.
	SeverityError Severity = "Error"
)

// ProgressMessage is a message emitted during progress.
type ProgressMessage struct {
	// Severity is the message severity.
	Severity Severity

	// Text is the message text.
	Text string

	// Timestamp is when the message was created.
	Timestamp time.Time

	// ActionID is the optional associated action ID.
	ActionID string

	// TaskID is the optional associated task ID.
	TaskID string

	// StepID is the optional associated step ID.
	StepID string
}

// ProgressAction represents a high-level action in a long-running operation.
type ProgressAction struct {
	ID        string
	Name      string
	StartedAt time.Time
	EndedAt   time.Time
}

// ProgressTask represents a task within an action.
type ProgressTask struct {
	ID        string
	ActionID  string
	Name      string
	StartedAt time.Time
	EndedAt   time.Time
}

// ProgressStep represents a step within a task.
type ProgressStep struct {
	ID        string
	TaskID    string
	Name      string
	StartedAt time.Time
	EndedAt   time.Time
}

// ProgressReporter is the interface for receiving progress updates.
//
// Implementations MUST be safe for concurrent use.
type ProgressReporter interface {
	// OnAction is called when an action starts or ends.
	OnAction(action ProgressAction)

	// OnTask is called when a task starts or ends.
	OnTask(task ProgressTask)

	// OnStep is called when a step starts or ends.
	OnStep(step ProgressStep)

	// OnMessage is called when a message is emitted.
	OnMessage(msg ProgressMessage)
}
