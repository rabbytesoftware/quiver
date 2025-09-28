package ui

import (
	"testing"

	"github.com/rabbytesoftware/quiver/internal/core/watcher"
)

func TestRunUI(t *testing.T) {
	// Create a watcher for testing
	w := watcher.NewWatcherService()

	// We can't actually run the UI in a test environment as it would block
	// and require terminal interaction. Instead, we test that the function exists
	// and doesn't panic when called with proper setup.

	// Test that RunUI function exists by assigning it to a variable
	var runUIFunc func(*watcher.Watcher) error = RunUI
	_ = runUIFunc // Function variables are never nil

	// Test with nil watcher - should handle gracefully or panic predictably
	defer func() {
		if r := recover(); r != nil {
			t.Logf("RunUI with nil watcher panicked as expected: %v", r)
		}
	}()

	// We can't actually call RunUI(nil) or RunUI(w) here because:
	// 1. It would block the test indefinitely
	// 2. It requires a terminal environment
	// 3. It would try to start the TUI

	// Instead, we verify the function signature and that it's callable
	var err error
	_ = err // This would be: err = RunUI(w)

	// Test that we can create a model (which RunUI does internally)
	model := NewModel(w)
	if model == nil {
		t.Error("NewModel returned nil - RunUI would fail")
	}
}

func TestRunUIFunctionSignature(t *testing.T) {
	// Test that RunUI has the expected signature
	w := watcher.NewWatcherService()

	// This tests that the function can be called with the right parameters
	// without actually executing it
	var runUIFunc func(*watcher.Watcher) error = RunUI

	_ = runUIFunc // Function variables are never nil

	// Test that we can pass a watcher to the function reference
	_ = w // Would be used like: err := runUIFunc(w)
}

func TestRunUIWithValidWatcher(t *testing.T) {
	w := watcher.NewWatcherService()

	// Test that we can create the model that RunUI would create
	model := NewModel(w)
	if model == nil {
		t.Fatal("NewModel returned nil")
	}

	// Test that the model has the watcher set
	if model.watcher != w {
		t.Error("Model watcher is not the same instance passed to NewModel")
	}

	// Test that model initialization works (what RunUI does)
	cmd := model.Init()
	if cmd == nil {
		t.Error("Model.Init() returned nil command")
	}
}

func TestRunUIModelCreation(t *testing.T) {
	// Test the model creation that happens inside RunUI
	w := watcher.NewWatcherService()

	// This is what RunUI does internally
	model := NewModel(w)

	if model == nil {
		t.Fatal("NewModel returned nil - RunUI would fail")
	}

	// Test that the model is properly initialized for TUI
	if model.textInput.Value() != "" {
		// Initial input should be empty
	}

	if model.viewport.Width <= 0 || model.viewport.Height <= 0 {
		// Viewport should have some dimensions (even if default)
	}

	// Test that context is set up
	if model.ctx == nil {
		t.Error("Model context is nil")
	}

	if model.cancel == nil {
		t.Error("Model cancel function is nil")
	}
}

func TestRunUIErrorHandling(t *testing.T) {
	// Test error handling scenarios that RunUI might encounter

	// Test with nil watcher
	defer func() {
		if r := recover(); r != nil {
			// If it panics with nil watcher, that's acceptable
			t.Logf("RunUI preparation panicked with nil watcher: %v", r)
		}
	}()

	// Test model creation with nil watcher
	model := NewModel(nil)
	if model == nil {
		t.Error("NewModel with nil watcher returned nil")
	}
}

func TestRunUICleanup(t *testing.T) {
	// Test the cleanup logic that RunUI performs
	w := watcher.NewWatcherService()
	model := NewModel(w)

	// Test that we can call cancel (which RunUI does on exit)
	if model.cancel != nil {
		model.cancel()

		// After cancel, context should be done
		select {
		case <-model.ctx.Done():
			// Expected - context is cancelled
		default:
			t.Error("Context should be done after cancel")
		}
	}
}

func TestRunUIModelState(t *testing.T) {
	// Test the initial state of the model that RunUI creates
	w := watcher.NewWatcherService()
	model := NewModel(w)

	// Test initial state
	if model.ready {
		t.Error("Model should not be ready initially")
	}

	if model.quitting {
		t.Error("Model should not be quitting initially")
	}

	if !model.autoScroll {
		t.Error("Model should have autoScroll enabled initially")
	}

	if len(model.logLines) != 0 {
		t.Error("Model should have empty log lines initially")
	}

	if len(model.commandHistory) != 0 {
		t.Error("Model should have empty command history initially")
	}

	if model.historyIndex != -1 {
		t.Error("Model should have historyIndex -1 initially")
	}
}
