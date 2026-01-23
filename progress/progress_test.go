package progress

import (
	"sync"
	"testing"
	"time"
)

func TestGetProgressReporter_NilReturnsNoOp(t *testing.T) {
	reporter := getProgressReporter(nil)
	if reporter == nil {
		t.Fatal("getProgressReporter(nil) returned nil, expected no-op")
	}

	// Should not panic
	reporter.OnAction(ProgressAction{ID: "test"})
	reporter.OnTask(ProgressTask{ID: "test"})
	reporter.OnStep(ProgressStep{ID: "test"})
	reporter.OnMessage(ProgressMessage{Text: "test"})
}

func TestGetProgressReporter_NonNilPassesThrough(t *testing.T) {
	called := false
	mock := &mockProgressReporter{
		onAction: func(a ProgressAction) { called = true },
	}

	reporter := getProgressReporter(mock)
	reporter.OnAction(ProgressAction{ID: "test"})

	if !called {
		t.Error("getProgressReporter did not pass through to provided reporter")
	}
}

func TestMakeThreadSafe_NilReturnsNoOp(t *testing.T) {
	reporter := MakeThreadSafe(nil)
	if reporter == nil {
		t.Fatal("MakeThreadSafe(nil) returned nil, expected no-op")
	}

	// Should not panic
	reporter.OnAction(ProgressAction{ID: "test"})
}

func TestMakeThreadSafe_ConcurrentCalls(t *testing.T) {
	calls := 0
	mock := &mockProgressReporter{
		onAction: func(a ProgressAction) {
			time.Sleep(1 * time.Millisecond) // Simulate work
			calls++
		},
	}

	reporter := MakeThreadSafe(mock)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			reporter.OnAction(ProgressAction{ID: "test"})
		}()
	}

	wg.Wait()

	if calls != 10 {
		t.Errorf("Expected 10 calls, got %d", calls)
	}
}

// mockProgressReporter is a test helper.
type mockProgressReporter struct {
	onAction  func(ProgressAction)
	onTask    func(ProgressTask)
	onStep    func(ProgressStep)
	onMessage func(ProgressMessage)
}

func (m *mockProgressReporter) OnAction(action ProgressAction) {
	if m.onAction != nil {
		m.onAction(action)
	}
}

func (m *mockProgressReporter) OnTask(task ProgressTask) {
	if m.onTask != nil {
		m.onTask(task)
	}
}

func (m *mockProgressReporter) OnStep(step ProgressStep) {
	if m.onStep != nil {
		m.onStep(step)
	}
}

func (m *mockProgressReporter) OnMessage(msg ProgressMessage) {
	if m.onMessage != nil {
		m.onMessage(msg)
	}
}
