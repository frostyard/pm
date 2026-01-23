package pm

import (
	"sync"
	"testing"
	"time"
)

// capturingReporter captures all progress events for testing.
type capturingReporter struct {
	mu       sync.Mutex
	actions  []ProgressAction
	tasks    []ProgressTask
	steps    []ProgressStep
	messages []ProgressMessage
}

func (r *capturingReporter) OnAction(action ProgressAction) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.actions = append(r.actions, action)
}

func (r *capturingReporter) OnTask(task ProgressTask) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks = append(r.tasks, task)
}

func (r *capturingReporter) OnStep(step ProgressStep) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.steps = append(r.steps, step)
}

func (r *capturingReporter) OnMessage(msg ProgressMessage) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.messages = append(r.messages, msg)
}

// T032: Test progress sequences
func TestProgressHelper_Sequences(t *testing.T) {
	t.Run("Action -> Task -> Step -> Message sequence", func(t *testing.T) {
		reporter := &capturingReporter{}
		helper := NewProgressHelper(nil, reporter)

		// Begin action
		actionID := helper.BeginAction("TestAction")
		if actionID == "" {
			t.Error("Expected non-empty action ID")
		}

		// Begin task
		taskID := helper.BeginTask("TestTask")
		if taskID == "" {
			t.Error("Expected non-empty task ID")
		}

		// Begin step
		stepID := helper.BeginStep("TestStep")
		if stepID == "" {
			t.Error("Expected non-empty step ID")
		}

		// Emit message
		helper.Info("Test message")

		// End step
		helper.EndStep()

		// End task
		helper.EndTask()

		// End action
		helper.EndAction()

		// Verify captured events
		if len(reporter.actions) != 2 {
			t.Errorf("Expected 2 actions (start+end), got %d", len(reporter.actions))
		}
		if len(reporter.tasks) != 2 {
			t.Errorf("Expected 2 tasks (start+end), got %d", len(reporter.tasks))
		}
		if len(reporter.steps) != 2 {
			t.Errorf("Expected 2 steps (start+end), got %d", len(reporter.steps))
		}
		if len(reporter.messages) != 1 {
			t.Errorf("Expected 1 message, got %d", len(reporter.messages))
		}

		// Verify IDs are linked
		if reporter.tasks[0].ActionID != actionID {
			t.Errorf("Task ActionID=%s doesn't match action ID=%s", reporter.tasks[0].ActionID, actionID)
		}
		if reporter.steps[0].TaskID != taskID {
			t.Errorf("Step TaskID=%s doesn't match task ID=%s", reporter.steps[0].TaskID, taskID)
		}
		if reporter.messages[0].ActionID != actionID {
			t.Errorf("Message ActionID=%s doesn't match action ID=%s", reporter.messages[0].ActionID, actionID)
		}
		if reporter.messages[0].TaskID != taskID {
			t.Errorf("Message TaskID=%s doesn't match task ID=%s", reporter.messages[0].TaskID, taskID)
		}
		if reporter.messages[0].StepID != stepID {
			t.Errorf("Message StepID=%s doesn't match step ID=%s", reporter.messages[0].StepID, stepID)
		}
	})

	t.Run("Multiple tasks in one action", func(t *testing.T) {
		reporter := &capturingReporter{}
		helper := NewProgressHelper(nil, reporter)

		actionID := helper.BeginAction("MultiTaskAction")

		helper.BeginTask("Task1")
		helper.Info("Task1 message")
		helper.EndTask()

		helper.BeginTask("Task2")
		helper.Info("Task2 message")
		helper.EndTask()

		helper.EndAction()

		// 2 actions (start+end), 4 tasks (2 start + 2 end)
		if len(reporter.actions) != 2 {
			t.Errorf("Expected 2 actions, got %d", len(reporter.actions))
		}
		if len(reporter.tasks) != 4 {
			t.Errorf("Expected 4 tasks (2 start + 2 end), got %d", len(reporter.tasks))
		}

		// All tasks should have same action ID
		for _, task := range reporter.tasks {
			if task.ActionID != actionID {
				t.Errorf("Task ActionID=%s doesn't match action ID=%s", task.ActionID, actionID)
			}
		}
	})

	t.Run("Timestamps are set", func(t *testing.T) {
		reporter := &capturingReporter{}
		helper := NewProgressHelper(nil, reporter)

		before := time.Now()
		helper.BeginAction("TimedAction")
		time.Sleep(10 * time.Millisecond)
		helper.EndAction()
		after := time.Now()

		if len(reporter.actions) != 2 {
			t.Fatalf("Expected 2 actions, got %d", len(reporter.actions))
		}

		start := reporter.actions[0]
		end := reporter.actions[1]

		if start.StartedAt.IsZero() {
			t.Error("Start action should have non-zero StartedAt")
		}
		if start.StartedAt.Before(before) || start.StartedAt.After(after) {
			t.Error("Start action timestamp out of expected range")
		}

		if end.EndedAt.IsZero() {
			t.Error("End action should have non-zero EndedAt")
		}
		if end.EndedAt.Before(start.StartedAt) {
			t.Error("End time should be after start time")
		}
	})
}

// T033: Test that Warning messages do not fail operations
func TestProgressHelper_WarningsDoNotFail(t *testing.T) {
	reporter := &capturingReporter{}
	helper := NewProgressHelper(nil, reporter)

	helper.BeginAction("ActionWithWarnings")
	helper.Warning("This is a warning")
	helper.Warning("Another warning")
	helper.Info("Normal operation continues")
	helper.EndAction()

	// Verify warnings were captured
	warnings := 0
	infos := 0
	for _, msg := range reporter.messages {
		switch msg.Severity {
		case SeverityWarning:
			warnings++
		case SeverityInfo:
			infos++
		}
	}

	if warnings != 2 {
		t.Errorf("Expected 2 warnings, got %d", warnings)
	}
	if infos != 1 {
		t.Errorf("Expected 1 info message, got %d", infos)
	}

	// The operation completed successfully (action was ended)
	if len(reporter.actions) != 2 {
		t.Error("Action should have completed despite warnings")
	}
}

// T034: Test that nil ProgressReporter does not panic
func TestProgressHelper_NilReporterSafe(t *testing.T) {
	helper := NewProgressHelper(nil, nil)

	// All these should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Operations panicked with nil reporter: %v", r)
		}
	}()

	actionID := helper.BeginAction("TestAction")
	if actionID != "" {
		t.Error("Expected empty action ID with nil reporter")
	}

	taskID := helper.BeginTask("TestTask")
	if taskID != "" {
		t.Error("Expected empty task ID with nil reporter")
	}

	stepID := helper.BeginStep("TestStep")
	if stepID != "" {
		t.Error("Expected empty step ID with nil reporter")
	}

	helper.Info("Info message")
	helper.Warning("Warning message")
	helper.Error("Error message")

	helper.EndStep()
	helper.EndTask()
	helper.EndAction()

	// If we get here, no panic occurred
}

// T047: Test concurrency safety
func TestProgressHelper_Concurrency(t *testing.T) {
	reporter := &capturingReporter{}

	// Run multiple goroutines emitting progress concurrently
	// Each goroutine uses its own ProgressHelper (realistic use case)
	var wg sync.WaitGroup
	concurrency := 10
	messagesPerGoroutine := 100

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// Each goroutine gets its own helper
			helper := NewProgressHelper(nil, reporter)
			helper.BeginAction("ConcurrentAction")
			helper.BeginTask("Task")
			for j := 0; j < messagesPerGoroutine; j++ {
				helper.Info("Message")
			}
			helper.EndTask()
			helper.EndAction()
		}(i)
	}

	wg.Wait()

	// Verify all messages were captured (should be safe due to mutex in capturingReporter)
	expectedMessages := concurrency * messagesPerGoroutine
	if len(reporter.messages) != expectedMessages {
		t.Errorf("Expected %d messages, got %d", expectedMessages, len(reporter.messages))
	}

	// Verify tasks were captured (concurrency tasks * 2 for start+end)
	expectedTasks := concurrency * 2
	if len(reporter.tasks) != expectedTasks {
		t.Errorf("Expected %d tasks, got %d", expectedTasks, len(reporter.tasks))
	}

	// Verify actions were captured (concurrency actions * 2 for start+end)
	expectedActions := concurrency * 2
	if len(reporter.actions) != expectedActions {
		t.Errorf("Expected %d actions, got %d", expectedActions, len(reporter.actions))
	}
}

func TestProgressHelper_MessageSeverities(t *testing.T) {
	reporter := &capturingReporter{}
	helper := NewProgressHelper(nil, reporter)

	helper.BeginAction("SeverityTest")
	helper.Info("Info message")
	helper.Warning("Warning message")
	helper.Error("Error message")
	helper.EndAction()

	if len(reporter.messages) != 3 {
		t.Fatalf("Expected 3 messages, got %d", len(reporter.messages))
	}

	severities := make(map[Severity]int)
	for _, msg := range reporter.messages {
		severities[msg.Severity]++
	}

	if severities[SeverityInfo] != 1 {
		t.Errorf("Expected 1 Info message, got %d", severities[SeverityInfo])
	}
	if severities[SeverityWarning] != 1 {
		t.Errorf("Expected 1 Warning message, got %d", severities[SeverityWarning])
	}
	if severities[SeverityError] != 1 {
		t.Errorf("Expected 1 Error message, got %d", severities[SeverityError])
	}
}

func TestProgressHelper_OrphanedEvents(t *testing.T) {
	t.Run("Task without action", func(t *testing.T) {
		reporter := &capturingReporter{}
		helper := NewProgressHelper(nil, reporter)

		// Start a task without an action
		taskID := helper.BeginTask("OrphanTask")
		helper.EndTask()

		if len(reporter.tasks) != 2 {
			t.Errorf("Expected 2 tasks, got %d", len(reporter.tasks))
		}

		// ActionID should be empty for orphaned task
		if reporter.tasks[0].ActionID != "" {
			t.Errorf("Expected empty ActionID for orphaned task, got %s", reporter.tasks[0].ActionID)
		}

		if taskID == "" {
			t.Error("Task ID should still be generated")
		}
	})

	t.Run("Step without task", func(t *testing.T) {
		reporter := &capturingReporter{}
		helper := NewProgressHelper(nil, reporter)

		helper.BeginAction("Action")
		stepID := helper.BeginStep("OrphanStep")
		helper.EndStep()
		helper.EndAction()

		if len(reporter.steps) != 2 {
			t.Errorf("Expected 2 steps, got %d", len(reporter.steps))
		}

		// TaskID should be empty for step without task
		if reporter.steps[0].TaskID != "" {
			t.Errorf("Expected empty TaskID for step without task, got %s", reporter.steps[0].TaskID)
		}

		if stepID == "" {
			t.Error("Step ID should still be generated")
		}
	})
}
