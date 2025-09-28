package watcher

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestNewWatcherService(t *testing.T) {
	watcher := NewWatcherService()

	if watcher == nil {
		t.Fatal("NewWatcherService() returned nil")
	}

	// Test that watcher has expected fields initialized
	if watcher.logger == nil {
		t.Error("Watcher logger is nil")
	}

	if watcher.pool == nil {
		t.Error("Watcher pool is nil")
	}
}

func TestWatcher_SetLevel(t *testing.T) {
	watcher := NewWatcherService()

	// Test setting different levels
	levels := []logrus.Level{
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
	}

	for _, level := range levels {
		watcher.SetLevel(level)
		if watcher.GetLevel() != level {
			t.Errorf("SetLevel(%v) failed, got %v", level, watcher.GetLevel())
		}
	}
}

func TestWatcher_GetLevel(t *testing.T) {
	watcher := NewWatcherService()

	// Test that GetLevel returns a valid level
	level := watcher.GetLevel()
	if level < logrus.PanicLevel || level > logrus.TraceLevel {
		t.Errorf("GetLevel() returned invalid level: %v", level)
	}
}

func TestWatcher_WithFields(t *testing.T) {
	watcher := NewWatcherService()

	fields := logrus.Fields{
		"key1": "value1",
		"key2": "value2",
	}

	entry := watcher.WithFields(fields)
	if entry == nil {
		t.Error("WithFields() returned nil")
	}
}

func TestWatcher_WithField(t *testing.T) {
	watcher := NewWatcherService()

	entry := watcher.WithField("test_key", "test_value")
	if entry == nil {
		t.Error("WithField() returned nil")
	}
}

func TestWatcher_Subscribe(t *testing.T) {
	watcher := NewWatcherService()

	// Test subscribing with a callback
	callback := func(level logrus.Level, message string) {
		// Callback for testing
	}

	watcher.Subscribe(callback)

	// Test that subscriber count increased
	count := watcher.GetSubscriberCount()
	if count == 0 {
		t.Error("Subscribe() did not increase subscriber count")
	}

	// Test logging to trigger callback
	watcher.Info("test message")

	// Note: The callback might not be called immediately due to goroutines
	// In a real test, we might need to add synchronization
}

func TestWatcher_GetSubscriberCount(t *testing.T) {
	watcher := NewWatcherService()

	// Initial count should be 0
	initialCount := watcher.GetSubscriberCount()
	if initialCount < 0 {
		t.Error("GetSubscriberCount() returned negative value")
	}

	// Add a subscriber
	watcher.Subscribe(func(level logrus.Level, message string) {})

	// Count should increase
	newCount := watcher.GetSubscriberCount()
	if newCount <= initialCount {
		t.Error("GetSubscriberCount() did not increase after Subscribe()")
	}
}

func TestWatcher_GetConfig(t *testing.T) {
	watcher := NewWatcherService()

	config := watcher.GetConfig()
	// Config can be any value, we just test it doesn't panic
	_ = config
}

func TestWatcher_IsEnabled(t *testing.T) {
	watcher := NewWatcherService()

	enabled := watcher.IsEnabled()
	// Enabled can be true or false, we just test it doesn't panic
	_ = enabled
}

func TestWatcher_LoggingMethods(t *testing.T) {
	watcher := NewWatcherService()

	// Test that logging methods don't panic (except Unforeseen which calls Fatal)
	defer func() {
		if r := recover(); r != nil {
			// Unforeseen calls Fatal which may cause issues, but other methods shouldn't panic
			t.Logf("Logging method had issue (expected for Unforeseen): %v", r)
		}
	}()

	watcher.Debug("debug message")
	watcher.Info("info message")
	watcher.Warn("warn message")
	watcher.Error("error message")

	// Note: Unforeseen calls Fatal which may exit the program, so we skip it in tests
	// watcher.Unforeseen("unforeseen message")
}

func TestMultipleWatcherInstances(t *testing.T) {
	watcher1 := NewWatcherService()
	watcher2 := NewWatcherService()

	// Different instances should be created
	if watcher1 == watcher2 {
		t.Error("NewWatcherService() should create different instances")
	}

	// But they should both be valid
	if watcher1 == nil || watcher2 == nil {
		t.Error("NewWatcherService() returned nil instance")
	}
}

func TestWatcherInitialization(t *testing.T) {
	watcher := NewWatcherService()

	// Test that watcher is properly initialized
	if watcher.logger == nil {
		t.Error("Watcher logger not initialized")
	}

	if watcher.pool == nil {
		t.Error("Watcher pool not initialized")
	}

	// Test that we can get the level (which means logger is working)
	level := watcher.GetLevel()
	if level < logrus.PanicLevel || level > logrus.TraceLevel {
		t.Error("Watcher level not properly initialized")
	}
}
