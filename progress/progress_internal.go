package progress

import "sync"

// noOpProgressReporter is a safe no-op progress reporter for when nil is provided.
type noOpProgressReporter struct{}

func (n *noOpProgressReporter) OnAction(action ProgressAction) {}
func (n *noOpProgressReporter) OnTask(task ProgressTask)       {}
func (n *noOpProgressReporter) OnStep(step ProgressStep)       {}
func (n *noOpProgressReporter) OnMessage(msg ProgressMessage)  {}

var noOpReporter = &noOpProgressReporter{}

// getProgressReporter returns the provided reporter or a no-op if nil.
func getProgressReporter(p ProgressReporter) ProgressReporter {
	if p == nil {
		return noOpReporter
	}
	return p
}

// threadSafeProgressReporter wraps a ProgressReporter with a mutex for thread safety.
type threadSafeProgressReporter struct {
	mu       sync.Mutex
	reporter ProgressReporter
}

func (t *threadSafeProgressReporter) OnAction(action ProgressAction) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.reporter.OnAction(action)
}

func (t *threadSafeProgressReporter) OnTask(task ProgressTask) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.reporter.OnTask(task)
}

func (t *threadSafeProgressReporter) OnStep(step ProgressStep) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.reporter.OnStep(step)
}

func (t *threadSafeProgressReporter) OnMessage(msg ProgressMessage) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.reporter.OnMessage(msg)
}

// MakeThreadSafe wraps a ProgressReporter to make it safe for concurrent use.
// If the reporter is already known to be thread-safe, this is unnecessary.
func MakeThreadSafe(p ProgressReporter) ProgressReporter {
	if p == nil {
		return noOpReporter
	}
	return &threadSafeProgressReporter{reporter: p}
}
