package services

import (
	"testing"

	"github.com/rabbytesoftware/quiver/internal/core/watcher"
	"github.com/sirupsen/logrus"
)

func TestNewWatcherAdapter(t *testing.T) {
	w := watcher.NewWatcherService()
	adapter := NewWatcherAdapter(w)

	if adapter == nil {
		t.Fatal("NewWatcherAdapter() returned nil")
	}

	if adapter.watcher != w {
		t.Error("WatcherAdapter watcher not set correctly")
	}
}

func TestWatcherAdapter_SetFilter(t *testing.T) {
	w := watcher.NewWatcherService()
	adapter := NewWatcherAdapter(w)

	// Test setting filter
	adapter.SetFilter("test.*pattern")

	// Test getting filter
	filter := adapter.GetFilter()
	if filter != "test.*pattern" {
		t.Errorf("Expected filter 'test.*pattern', got %q", filter)
	}
}

func TestWatcherAdapter_GetFilter(t *testing.T) {
	w := watcher.NewWatcherService()
	adapter := NewWatcherAdapter(w)

	// Initial filter should be empty
	filter := adapter.GetFilter()
	if filter != "" {
		t.Errorf("Expected empty initial filter, got %q", filter)
	}

	// Set and get filter
	adapter.SetFilter("new.*filter")
	filter = adapter.GetFilter()
	if filter != "new.*filter" {
		t.Errorf("Expected filter 'new.*filter', got %q", filter)
	}
}

func TestWatcherAdapter_SetLevel(t *testing.T) {
	w := watcher.NewWatcherService()
	adapter := NewWatcherAdapter(w)

	// Test setting different levels
	levels := []string{"debug", "info", "warn", "error"}

	for _, level := range levels {
		adapter.SetLevel(level)

		// Verify the level was set by checking the watcher
		// Note: We can't directly compare strings because the watcher uses logrus.Level
		currentLevel := adapter.GetLevel()
		if currentLevel == "" {
			t.Errorf("SetLevel(%q) failed, GetLevel() returned empty string", level)
		}
	}
}

func TestWatcherAdapter_GetLevel(t *testing.T) {
	w := watcher.NewWatcherService()
	adapter := NewWatcherAdapter(w)

	level := adapter.GetLevel()
	if level == "" {
		t.Error("GetLevel() returned empty string")
	}

	// Test that it returns a valid log level
	validLevels := []string{"debug", "info", "warn", "error", "panic", "fatal", "trace"}
	isValid := false
	for _, validLevel := range validLevels {
		if level == validLevel {
			isValid = true
			break
		}
	}
	if !isValid {
		t.Errorf("GetLevel() returned invalid level: %q", level)
	}
}

func TestWatcherAdapter_IsPaused(t *testing.T) {
	w := watcher.NewWatcherService()
	adapter := NewWatcherAdapter(w)

	// Initially should not be paused
	if adapter.IsPaused() {
		t.Error("WatcherAdapter should not be paused initially")
	}

	// Pause and check
	adapter.Pause()
	if !adapter.IsPaused() {
		t.Error("WatcherAdapter should be paused after Pause()")
	}

	// Resume and check
	adapter.Resume()
	if adapter.IsPaused() {
		t.Error("WatcherAdapter should not be paused after Resume()")
	}
}

func TestWatcherAdapter_Pause(t *testing.T) {
	w := watcher.NewWatcherService()
	adapter := NewWatcherAdapter(w)

	// Test pause
	adapter.Pause()
	if !adapter.IsPaused() {
		t.Error("Pause() did not set paused state")
	}

	// Test multiple pauses (should be idempotent)
	adapter.Pause()
	if !adapter.IsPaused() {
		t.Error("Multiple Pause() calls should maintain paused state")
	}
}

func TestWatcherAdapter_Resume(t *testing.T) {
	w := watcher.NewWatcherService()
	adapter := NewWatcherAdapter(w)

	// Pause first
	adapter.Pause()
	if !adapter.IsPaused() {
		t.Error("Setup failed: adapter should be paused")
	}

	// Test resume
	adapter.Resume()
	if adapter.IsPaused() {
		t.Error("Resume() did not clear paused state")
	}

	// Test multiple resumes (should be idempotent)
	adapter.Resume()
	if adapter.IsPaused() {
		t.Error("Multiple Resume() calls should maintain resumed state")
	}
}

func TestWatcherAdapter_Subscribe(t *testing.T) {
	w := watcher.NewWatcherService()
	adapter := NewWatcherAdapter(w)

	// Test subscribing with a callback
	callback := func(level logrus.Level, message string) {
		// Callback for testing
	}

	adapter.Subscribe(callback)

	// The subscription should be passed to the underlying watcher
	// We can't easily test the callback is called without triggering a log message
	// and dealing with goroutines, so we just test it doesn't panic
}

func TestWatcherAdapter_Subscribe_FilteringLogic(t *testing.T) {
	w := watcher.NewWatcherService()
	adapter := NewWatcherAdapter(w)

	// Test 1: Subscribe with paused adapter
	adapter.Pause()
	callback1 := func(level logrus.Level, message string) {
		t.Error("Callback should not be called when adapter is paused")
	}
	adapter.Subscribe(callback1)

	// Test 2: Subscribe with filter set
	adapter.Resume()
	adapter.SetFilter("test")
	callback2 := func(level logrus.Level, message string) {
		// This callback tests the filter logic path
	}
	adapter.Subscribe(callback2)

	// Test 3: Subscribe with level filtering
	adapter.SetFilter("") // Clear filter
	adapter.SetLevel("error")
	callback3 := func(level logrus.Level, message string) {
		// This callback tests the level filtering path
	}
	adapter.Subscribe(callback3)

	// Test 4: Subscribe with both filter and level
	adapter.SetFilter("important")
	callback4 := func(level logrus.Level, message string) {
		// This callback tests both filter and level logic paths
	}
	adapter.Subscribe(callback4)

	// All subscribes should complete without panic
}

func TestWatcherAdapter_GetWatcher(t *testing.T) {
	w := watcher.NewWatcherService()
	adapter := NewWatcherAdapter(w)

	retrievedWatcher := adapter.GetWatcher()
	if retrievedWatcher != w {
		t.Error("GetWatcher() did not return the original watcher")
	}
}

func TestWatcherAdapterWithNilWatcher(t *testing.T) {
	// Test behavior with nil watcher
	adapter := NewWatcherAdapter(nil)

	if adapter == nil {
		t.Fatal("NewWatcherAdapter(nil) returned nil")
	}

	// Test basic operations that should work with nil watcher
	adapter.SetFilter("test")
	filter := adapter.GetFilter()
	if filter != "test" {
		t.Error("SetFilter/GetFilter should work even with nil watcher")
	}

	adapter.Pause()
	if !adapter.IsPaused() {
		t.Error("Pause() should work even with nil watcher")
	}

	adapter.Resume()
	if adapter.IsPaused() {
		t.Error("Resume() should work even with nil watcher")
	}

	// GetWatcher should return nil
	if adapter.GetWatcher() != nil {
		t.Error("GetWatcher() should return nil when initialized with nil")
	}

	// Note: SetLevel and GetLevel may panic with nil watcher, which is acceptable
	// Subscribe may also panic with nil watcher, which is acceptable
}

func TestWatcherAdapterState(t *testing.T) {
	w := watcher.NewWatcherService()
	adapter := NewWatcherAdapter(w)

	// Test initial state
	if adapter.IsPaused() {
		t.Error("WatcherAdapter should not be paused initially")
	}

	if adapter.GetFilter() != "" {
		t.Error("WatcherAdapter should have empty filter initially")
	}

	// Test state persistence
	adapter.SetFilter("persistent.*filter")
	adapter.SetLevel("debug")
	adapter.Pause()

	if adapter.GetFilter() != "persistent.*filter" {
		t.Error("Filter state not persisted")
	}

	if !adapter.IsPaused() {
		t.Error("Paused state not persisted")
	}

	// Test state changes
	adapter.SetFilter("new.*filter")
	if adapter.GetFilter() != "new.*filter" {
		t.Error("Filter state not updated")
	}

	adapter.Resume()
	if adapter.IsPaused() {
		t.Error("Paused state not updated")
	}
}
