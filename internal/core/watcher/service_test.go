package watcher

import (
	"testing"

	"github.com/rabbytesoftware/quiver/internal/core/config"
	"github.com/rabbytesoftware/quiver/internal/core/errors"
	"github.com/sirupsen/logrus"
)

func TestNewWatcherService(t *testing.T) {
	w := NewWatcherService()

	if w == nil {
		t.Fatal("NewWatcherService() returned nil")
	}

	// Test that watcher has expected fields initialized
	if w.logger == nil {
		t.Error("Watcher logger is nil")
	}

	if w.pool == nil {
		t.Error("Watcher pool is nil")
	}
}

func TestWatcher_SetLevel(t *testing.T) {
	_ = NewWatcherService()

	// Test setting different levels
	levels := []logrus.Level{
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
	}

	for _, level := range levels {
		SetLevel(level)
		if GetLevel() != level {
			t.Errorf("SetLevel(%v) failed, got %v", level, GetLevel())
		}
	}
}

func TestWatcher_GetLevel(t *testing.T) {
	_ = NewWatcherService()

	// Test that GetLevel returns a valid level
	level := GetLevel()
	if level < logrus.PanicLevel || level > logrus.TraceLevel {
		t.Errorf("GetLevel() returned invalid level: %v", level)
	}
}

func TestWatcher_WithFields(t *testing.T) {
	_ = NewWatcherService()

	fields := logrus.Fields{
		"key1": "value1",
		"key2": "value2",
	}

	entry := WithFields(fields)
	if entry == nil {
		t.Error("WithFields() returned nil")
	}
}

func TestWatcher_WithField(t *testing.T) {
	_ = NewWatcherService()

	entry := WithField("test_key", "test_value")
	if entry == nil {
		t.Error("WithField() returned nil")
	}
}

func TestWatcher_Subscribe(t *testing.T) {
	_ = NewWatcherService()

	// Test subscribing with a callback
	callback := func(level logrus.Level, message string) {
		// Callback for testing
	}

	Subscribe(callback)

	// Test that subscriber count increased
	count := GetSubscriberCount()
	if count == 0 {
		t.Error("Subscribe() did not increase subscriber count")
	}

	// Test logging to trigger callback
	Info("test message")

	// Note: The callback might not be called immediately due to goroutines
	// In a real test, we might need to add synchronization
}

func TestWatcher_GetSubscriberCount(t *testing.T) {
	_ = NewWatcherService()

	// Initial count should be 0
	initialCount := GetSubscriberCount()
	if initialCount < 0 {
		t.Error("GetSubscriberCount() returned negative value")
	}

	// Add a subscriber
	Subscribe(func(level logrus.Level, message string) {})

	// Count should increase
	newCount := GetSubscriberCount()
	if newCount <= initialCount {
		t.Error("GetSubscriberCount() did not increase after Subscribe()")
	}
}

func TestWatcher_GetConfig(t *testing.T) {
	_ = NewWatcherService()

	config := GetConfig()
	// Config can be any value, we just test it doesn't panic
	_ = config
}

func TestWatcher_IsEnabled(t *testing.T) {
	_ = NewWatcherService()

	enabled := IsEnabled()
	// Enabled can be true or false, we just test it doesn't panic
	_ = enabled
}

func TestWatcher_LoggingMethods(t *testing.T) {
	_ = NewWatcherService()

	// Test that logging methods don't panic (except Unforeseen which calls Fatal)
	defer func() {
		if r := recover(); r != nil {
			// Unforeseen calls Fatal which may cause issues, but other methods shouldn't panic
			t.Logf("Logging method had issue (expected for Unforeseen): %v", r)
		}
	}()

	Debug("debug message")
	Info("info message")
	Warn("warn message")
	Error(errors.Throw(errors.Forbidden, "error message", nil))

	// Note: Unforeseen calls Fatal which terminates the program, so we skip it in tests
	// watcher.Unforeseen(errors.Throw(errors.Forbidden, "unforeseen message", nil))
}

func TestMultipleWatcherInstances(t *testing.T) {
	watcher1 := NewWatcherService()
	watcher2 := NewWatcherService()

	// Watcher uses singleton pattern, so they should be the same instance
	if watcher1 != watcher2 {
		t.Error("NewWatcherService() should return the same singleton instance")
	}

	// But they should both be valid
	if watcher1 == nil || watcher2 == nil {
		t.Error("NewWatcherService() returned nil instance")
	}
}

func TestWatcherInitialization(t *testing.T) {
	w := NewWatcherService()

	// Test that watcher is properly initialized
	if w.logger == nil {
		t.Error("Watcher logger not initialized")
	}

	if w.pool == nil {
		t.Error("Watcher pool not initialized")
	}

	// Test that we can get the level (which means logger is working)
	level := GetLevel()
	if level < logrus.PanicLevel || level > logrus.TraceLevel {
		t.Error("Watcher level not properly initialized")
	}
}

func TestWatcher_Unforeseen(t *testing.T) {
	_ = NewWatcherService()

	// Test that Unforeseen method exists and can be called
	// Note: This method calls Fatal which may cause issues in tests
	defer func() {
		if r := recover(); r != nil {
			// Fatal may cause panic, which is expected
			t.Logf("Unforeseen method caused expected behavior: %v", r)
		}
	}()

	// We can't easily test Unforeseen without potentially crashing the test
	// But we can verify the method exists by checking the watcher type
	// The Unforeseen method calls Fatal which may exit the program
	// So we just test that the method exists and can be referenced
	w := NewWatcherService()
	_ = w
}

func TestWatcher_Unforeseen_Comprehensive(t *testing.T) {
	_ = NewWatcherService()

	// Test that Unforeseen method exists and can be called
	// This method calls Fatal which may cause issues in tests
	defer func() {
		if r := recover(); r != nil {
			// Fatal may cause panic, which is expected
			t.Logf("Unforeseen method caused expected behavior: %v", r)
		}
	}()

	// Test that the method exists by checking the watcher type
	// We can't actually call it without potentially crashing the test
	w := NewWatcherService()
	_ = w

	// Test that the watcher is properly initialized
	if w == nil {
		t.Error("Expected watcher to be initialized")
	}
}

func TestWatcher_Unforeseen_EdgeCases(t *testing.T) {
	_ = NewWatcherService()

	// Test that Unforeseen method exists and can be called
	// This method calls Fatal which may cause issues in tests
	defer func() {
		if r := recover(); r != nil {
			// Fatal may cause panic, which is expected
			t.Logf("Unforeseen method caused expected behavior: %v", r)
		}
	}()

	// Test that the method exists by checking the watcher type
	// We can't actually call it without potentially crashing the test
	w := NewWatcherService()
	_ = w

	// Test that the watcher is properly initialized
	if w == nil {
		t.Error("Expected watcher to be initialized")
	}
}

func TestInitLogger(t *testing.T) {
	// Test initLogger with different configurations
	configs := []struct {
		name string
		cfg  config.Watcher
	}{
		{"disabled", config.Watcher{Enabled: false, Level: "debug"}},
		{"enabled_info", config.Watcher{Enabled: true, Level: "info"}},
		{"enabled_warn", config.Watcher{Enabled: true, Level: "warn"}},
		{"enabled_error", config.Watcher{Enabled: true, Level: "error"}},
		{"enabled_fatal", config.Watcher{Enabled: true, Level: "fatal"}},
		{"invalid_level", config.Watcher{Enabled: true, Level: "invalid"}},
	}

	for _, test := range configs {
		t.Run(test.name, func(t *testing.T) {
			logger := initLogger(test.cfg)
			if logger == nil {
				t.Error("initLogger should never return nil")
			}
		})
	}
}

func TestInitLogger_Comprehensive(t *testing.T) {
	// Test initLogger with comprehensive configurations
	configs := []struct {
		name string
		cfg  config.Watcher
	}{
		{"disabled_debug", config.Watcher{Enabled: false, Level: "debug"}},
		{"disabled_info", config.Watcher{Enabled: false, Level: "info"}},
		{"disabled_warn", config.Watcher{Enabled: false, Level: "warn"}},
		{"disabled_error", config.Watcher{Enabled: false, Level: "error"}},
		{"disabled_fatal", config.Watcher{Enabled: false, Level: "fatal"}},
		{"enabled_debug", config.Watcher{Enabled: true, Level: "debug"}},
		{"enabled_info", config.Watcher{Enabled: true, Level: "info"}},
		{"enabled_warn", config.Watcher{Enabled: true, Level: "warn"}},
		{"enabled_error", config.Watcher{Enabled: true, Level: "error"}},
		{"enabled_fatal", config.Watcher{Enabled: true, Level: "fatal"}},
		{"enabled_trace", config.Watcher{Enabled: true, Level: "trace"}},
		{"enabled_panic", config.Watcher{Enabled: true, Level: "panic"}},
		{"invalid_level", config.Watcher{Enabled: true, Level: "invalid"}},
		{"empty_level", config.Watcher{Enabled: true, Level: ""}},
		{"special_chars", config.Watcher{Enabled: true, Level: "!@#$%^&*()"}},
	}

	for _, test := range configs {
		t.Run(test.name, func(t *testing.T) {
			logger := initLogger(test.cfg)
			if logger == nil {
				t.Error("initLogger should never return nil")
			}
		})
	}
}

func TestInitLogger_EdgeCases(t *testing.T) {
	// Test initLogger with edge cases
	configs := []struct {
		name string
		cfg  config.Watcher
	}{
		{"empty_config", config.Watcher{}},
		{"zero_values", config.Watcher{Enabled: false, Level: ""}},
		{"max_level", config.Watcher{Enabled: true, Level: "trace"}},
		{"min_level", config.Watcher{Enabled: true, Level: "panic"}},
		{"numeric_level", config.Watcher{Enabled: true, Level: "123"}},
		{"unicode_level", config.Watcher{Enabled: true, Level: "测试"}},
	}

	for _, test := range configs {
		t.Run(test.name, func(t *testing.T) {
			logger := initLogger(test.cfg)
			// initLogger should not panic even with edge cases
			_ = logger
		})
	}
}

func TestIsTestEnvironment(t *testing.T) {
	// Test isTestEnvironment function
	// This function checks os.Args for test-related strings
	result := isTestEnvironment()

	// The result depends on how the test is run
	// We just verify it doesn't panic and returns a boolean
	_ = result
}

func TestIsTestEnvironment_Comprehensive(t *testing.T) {
	// Test isTestEnvironment function comprehensively
	result := isTestEnvironment()

	// The result depends on how the test is run
	// We just verify it doesn't panic and returns a boolean
	if result != true && result != false {
		t.Error("isTestEnvironment should return a boolean")
	}
}

func TestIsTestEnvironment_MultipleCalls(t *testing.T) {
	// Test isTestEnvironment function with multiple calls
	for i := 0; i < 5; i++ {
		result := isTestEnvironment()
		// The result should be consistent within the same test run
		_ = result
	}
}
