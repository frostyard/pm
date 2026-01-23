package pm

import "github.com/frostyard/pm/progress"

// Re-export progress types for backward compatibility
type (
	// ProgressReporter is the interface for receiving progress updates.
	ProgressReporter = progress.ProgressReporter

	// ProgressAction represents a high-level action in a long-running operation.
	ProgressAction = progress.ProgressAction

	// ProgressTask represents a task within an action.
	ProgressTask = progress.ProgressTask

	// ProgressStep represents a step within a task.
	ProgressStep = progress.ProgressStep

	// ProgressMessage is a message emitted during progress.
	ProgressMessage = progress.ProgressMessage

	// ProgressHelper provides a convenient API for backends to emit progress updates.
	ProgressHelper = progress.ProgressHelper

	// Severity represents the severity level of a progress message.
	Severity = progress.Severity
)

// Re-export severity constants
const (
	SeverityInfo    = progress.SeverityInfo
	SeverityWarning = progress.SeverityWarning
	SeverityError   = progress.SeverityError
)

// NewProgressHelper creates a new progress helper with progress reporting.
func NewProgressHelper(defaultReporter, overrideReporter ProgressReporter) *ProgressHelper {
	return progress.NewProgressHelper(defaultReporter, overrideReporter)
}
