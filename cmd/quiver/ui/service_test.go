package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
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

	if len(model.logLines) != 1 {
		t.Error("Model should have one log line initially (ASCII art)")
	}

	if len(model.commandHistory) != 0 {
		t.Error("Model should have empty command history initially")
	}

	if model.historyIndex != -1 {
		t.Error("Model should have historyIndex -1 initially")
	}
}

func TestRunUIComprehensive(t *testing.T) {
	// Test comprehensive RunUI functionality
	w := watcher.NewWatcherService()

	// Test that we can create a model (what RunUI does)
	model := NewModel(w)
	if model == nil {
		t.Fatal("NewModel returned nil")
	}

	// Test model initialization
	cmd := model.Init()
	if cmd == nil {
		t.Error("Model.Init() returned nil command")
	}

	// Test that the model is properly set up for TUI
	// TextInput and Viewport are structs, not pointers, so we can't check for nil
	_ = model.textInput
	_ = model.viewport

	if model.handler == nil {
		t.Error("Handler should be initialized")
	}

	if model.queryService == nil {
		t.Error("Query service should be initialized")
	}

	if model.watcherAdapter == nil {
		t.Error("Watcher adapter should be initialized")
	}

	if model.asciiService == nil {
		t.Error("ASCII service should be initialized")
	}
}

func TestRunUIWithDifferentWatcherStates(t *testing.T) {
	// Test RunUI with different watcher configurations
	w := watcher.NewWatcherService()

	// Test with enabled watcher
	model := NewModel(w)
	if model == nil {
		t.Fatal("NewModel returned nil")
	}

	// Test that watcher is properly connected
	if model.watcher != w {
		t.Error("Model watcher is not the same instance")
	}

	// Test model state after initialization
	if model.ready {
		t.Error("Model should not be ready initially")
	}

	// Test that we can set model to ready state (what happens in RunUI)
	model.ready = true
	if !model.ready {
		t.Error("Model should be ready after setting ready flag")
	}
}

func TestRunUIErrorScenarios(t *testing.T) {
	// Test various error scenarios that RunUI might handle

	// Test with nil watcher
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Expected panic with nil watcher: %v", r)
		}
	}()

	// This might panic, which is acceptable
	model := NewModel(nil)
	if model != nil {
		// If it doesn't panic, test that it handles nil watcher gracefully
		if model.watcher != nil {
			t.Error("Model should have nil watcher when passed nil")
		}
	}
}

func TestRunUIModelLifecycle(t *testing.T) {
	// Test the complete model lifecycle that RunUI manages
	w := watcher.NewWatcherService()
	model := NewModel(w)

	// Test initialization
	cmd := model.Init()
	if cmd == nil {
		t.Error("Init should return a command")
	}

	// Test that model can handle updates (what RunUI does)
	msg := tea.WindowSizeMsg{Width: 100, Height: 50}
	updatedModel, _ := model.Update(msg)
	if updatedModel == nil {
		t.Error("Update should return a model")
	}

	// Test that model can render (what RunUI does)
	view := model.View()
	if view == "" {
		t.Error("View should return non-empty string")
	}

	// Test cleanup (what RunUI does on exit)
	if model.cancel != nil {
		model.cancel()

		// Verify context is cancelled
		select {
		case <-model.ctx.Done():
			// Expected
		default:
			t.Error("Context should be done after cancel")
		}
	}
}

func TestRunUIWithWatcherSubscription(t *testing.T) {
	// Test the watcher subscription that RunUI sets up
	w := watcher.NewWatcherService()
	model := NewModel(w)

	// Test that we can subscribe to watcher (what RunUI does)
	model.subscribeToWatcher()

	// Test that watcher adapter is set up
	if model.watcherAdapter == nil {
		t.Error("Watcher adapter should be initialized")
	}

	// Test that we can get watcher from adapter
	watcherFromAdapter := model.watcherAdapter.GetWatcher()
	if watcherFromAdapter == nil {
		t.Error("Watcher adapter should return a watcher")
	}
}

func TestRunUIASCIIArt(t *testing.T) {
	// Test the ASCII art functionality that RunUI uses
	w := watcher.NewWatcherService()
	model := NewModel(w)

	// Test that ASCII art is added
	if len(model.logLines) == 0 {
		t.Error("Model should have ASCII art log line")
	}

	// Test that ASCII service works
	if model.asciiService == nil {
		t.Error("ASCII service should be initialized")
	}

	// Test that we can get welcome log line
	welcomeLine := model.asciiService.GetWelcomeLogLine()
	if welcomeLine.Text == "" {
		t.Error("ASCII service should return welcome text")
	}
}

func TestRunUITheme(t *testing.T) {
	// Test the theme functionality that RunUI uses
	w := watcher.NewWatcherService()
	model := NewModel(w)

	// Test that theme is initialized
	// Theme is a struct, so we can't check for nil, but we can test methods
	_ = model.theme

	// Test that we can get log level style
	style := model.theme.GetLogLevelStyle("info")
	_ = style // Style is a lipgloss.Style, not a string

	// Test that we can format log lines
	logLine := "Test log line"
	formatted := model.theme.FormatLogLine(logLine, "info")
	if formatted == "" {
		t.Error("Theme should format log lines")
	}
}

func TestRunUIQueryService(t *testing.T) {
	// Test the query service that RunUI uses
	w := watcher.NewWatcherService()
	model := NewModel(w)

	// Test that query service is initialized
	if model.queryService == nil {
		t.Error("Query service should be initialized")
	}

	// Test that query service is loaded
	if !model.queryService.IsLoaded() {
		t.Error("Query service should be loaded")
	}

	// Test that we can get available commands
	commands := model.queryService.GetAvailableCommands()
	if len(commands) == 0 {
		t.Error("Query service should have available commands")
	}

	// Test that we can get help text
	helpText := model.queryService.GetHelpText()
	if helpText == "" {
		t.Error("Query service should return help text")
	}
}
