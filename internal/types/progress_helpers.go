package types

import (
	"time"

	"github.com/google/uuid"
)

// ProgressHelper provides a convenient API for backends to emit progress updates.
// It tracks the current action/task/step context and handles ID generation.
type ProgressHelper struct {
	reporter      ProgressReporter
	currentAction *ProgressAction
	currentTask   *ProgressTask
	currentStep   *ProgressStep
}

// NewProgressHelper creates a new progress helper with progress reporting.
// It uses the override reporter if non-nil, otherwise falls back to the default reporter.
// If both are nil, returns a helper that no-ops all operations.
func NewProgressHelper(defaultReporter, overrideReporter ProgressReporter) *ProgressHelper {
	reporter := overrideReporter
	if reporter == nil {
		reporter = defaultReporter
	}
	return &ProgressHelper{
		reporter: reporter,
	}
}

// BeginAction starts a new action and returns its ID.
func (h *ProgressHelper) BeginAction(name string) string {
	if h.reporter == nil {
		return ""
	}

	action := ProgressAction{
		ID:        uuid.New().String(),
		Name:      name,
		StartedAt: time.Now(),
	}
	h.currentAction = &action
	h.reporter.OnAction(action)
	return action.ID
}

// EndAction ends the current action.
func (h *ProgressHelper) EndAction() {
	if h.reporter == nil || h.currentAction == nil {
		return
	}

	action := *h.currentAction
	action.EndedAt = time.Now()
	h.reporter.OnAction(action)
	h.currentAction = nil
}

// BeginTask starts a new task and returns its ID.
func (h *ProgressHelper) BeginTask(name string) string {
	if h.reporter == nil {
		return ""
	}

	task := ProgressTask{
		ID:        uuid.New().String(),
		Name:      name,
		StartedAt: time.Now(),
	}
	if h.currentAction != nil {
		task.ActionID = h.currentAction.ID
	}
	h.currentTask = &task
	h.reporter.OnTask(task)
	return task.ID
}

// EndTask ends the current task.
func (h *ProgressHelper) EndTask() {
	if h.reporter == nil || h.currentTask == nil {
		return
	}

	task := *h.currentTask
	task.EndedAt = time.Now()
	h.reporter.OnTask(task)
	h.currentTask = nil
}

// BeginStep starts a new step and returns its ID.
func (h *ProgressHelper) BeginStep(name string) string {
	if h.reporter == nil {
		return ""
	}

	step := ProgressStep{
		ID:        uuid.New().String(),
		Name:      name,
		StartedAt: time.Now(),
	}
	if h.currentTask != nil {
		step.TaskID = h.currentTask.ID
	}
	h.currentStep = &step
	h.reporter.OnStep(step)
	return step.ID
}

// EndStep ends the current step.
func (h *ProgressHelper) EndStep() {
	if h.reporter == nil || h.currentStep == nil {
		return
	}

	step := *h.currentStep
	step.EndedAt = time.Now()
	h.reporter.OnStep(step)
	h.currentStep = nil
}

// Info emits an informational message.
func (h *ProgressHelper) Info(text string) {
	h.message(SeverityInfo, text)
}

// Warning emits a warning message.
func (h *ProgressHelper) Warning(text string) {
	h.message(SeverityWarning, text)
}

// Error emits an error message.
func (h *ProgressHelper) Error(text string) {
	h.message(SeverityError, text)
}

// message emits a message with the given severity.
func (h *ProgressHelper) message(severity Severity, text string) {
	if h.reporter == nil {
		return
	}

	msg := ProgressMessage{
		Text:      text,
		Severity:  severity,
		Timestamp: time.Now(),
	}

	if h.currentAction != nil {
		msg.ActionID = h.currentAction.ID
	}
	if h.currentTask != nil {
		msg.TaskID = h.currentTask.ID
	}
	if h.currentStep != nil {
		msg.StepID = h.currentStep.ID
	}

	h.reporter.OnMessage(msg)
}
