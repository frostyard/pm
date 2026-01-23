package progress

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

// EndAction marks the current action as ended.
func (h *ProgressHelper) EndAction() {
	if h.reporter == nil || h.currentAction == nil {
		return
	}

	h.currentAction.EndedAt = time.Now()
	h.reporter.OnAction(*h.currentAction)
	h.currentAction = nil
	h.currentTask = nil
	h.currentStep = nil
}

// BeginTask starts a new task within the current action and returns its ID.
func (h *ProgressHelper) BeginTask(name string) string {
	if h.reporter == nil {
		return ""
	}

	actionID := ""
	if h.currentAction != nil {
		actionID = h.currentAction.ID
	}

	task := ProgressTask{
		ID:        uuid.New().String(),
		ActionID:  actionID,
		Name:      name,
		StartedAt: time.Now(),
	}
	h.currentTask = &task
	h.reporter.OnTask(task)
	return task.ID
}

// EndTask marks the current task as ended.
func (h *ProgressHelper) EndTask() {
	if h.reporter == nil || h.currentTask == nil {
		return
	}

	h.currentTask.EndedAt = time.Now()
	h.reporter.OnTask(*h.currentTask)
	h.currentTask = nil
	h.currentStep = nil
}

// BeginStep starts a new step within the current task and returns its ID.
func (h *ProgressHelper) BeginStep(name string) string {
	if h.reporter == nil {
		return ""
	}

	taskID := ""
	if h.currentTask != nil {
		taskID = h.currentTask.ID
	}

	step := ProgressStep{
		ID:        uuid.New().String(),
		TaskID:    taskID,
		Name:      name,
		StartedAt: time.Now(),
	}
	h.currentStep = &step
	h.reporter.OnStep(step)
	return step.ID
}

// EndStep marks the current step as ended.
func (h *ProgressHelper) EndStep() {
	if h.reporter == nil || h.currentStep == nil {
		return
	}

	h.currentStep.EndedAt = time.Now()
	h.reporter.OnStep(*h.currentStep)
	h.currentStep = nil
}

// Info emits an informational message.
func (h *ProgressHelper) Info(text string) {
	h.message(SeverityInfo, text)
}

// Warning emits a warning message (does not fail the operation).
func (h *ProgressHelper) Warning(text string) {
	h.message(SeverityWarning, text)
}

// Error emits an error message.
func (h *ProgressHelper) Error(text string) {
	h.message(SeverityError, text)
}

// message emits a progress message with the specified severity.
func (h *ProgressHelper) message(severity Severity, text string) {
	if h.reporter == nil {
		return
	}

	msg := ProgressMessage{
		Severity:  severity,
		Text:      text,
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
