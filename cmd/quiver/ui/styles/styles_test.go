package styles

import (
	"testing"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func TestNewDefaultTheme(t *testing.T) {
	theme := NewDefaultTheme()

	if theme.PrimaryColor == "" {
		t.Error("NewDefaultTheme() PrimaryColor is empty")
	}

	if theme.SecondaryColor == "" {
		t.Error("NewDefaultTheme() SecondaryColor is empty")
	}

	if theme.AccentColor == "" {
		t.Error("NewDefaultTheme() AccentColor is empty")
	}

	if theme.ErrorColor == "" {
		t.Error("NewDefaultTheme() ErrorColor is empty")
	}

	if theme.SuccessColor == "" {
		t.Error("NewDefaultTheme() SuccessColor is empty")
	}

	if theme.WarningColor == "" {
		t.Error("NewDefaultTheme() WarningColor is empty")
	}

	// InfoColor doesn't exist, skip this check

	if theme.MutedColor == "" {
		t.Error("NewDefaultTheme() MutedColor is empty")
	}
}

func TestTheme_GetLogLevelStyle(t *testing.T) {
	theme := NewDefaultTheme()

	// Test different log levels
	levels := []string{"debug", "info", "warn", "error", "unknown"}

	for _, level := range levels {
		style := theme.GetLogLevelStyle(level)
		// Can't compare lipgloss.Style directly, just check it doesn't panic
		_ = style
	}
}

func TestTheme_FormatLogLine(t *testing.T) {
	theme := NewDefaultTheme()

	// Test formatting different log levels
	testCases := []struct {
		message string
		level   string
	}{
		{"Debug message", "debug"},
		{"Info message", "info"},
		{"Warning message", "warn"},
		{"Error message", "error"},
		{"Unknown level message", "unknown"},
	}

	for _, tc := range testCases {
		formatted := theme.FormatLogLine(tc.message, tc.level)
		if formatted == "" {
			t.Errorf("FormatLogLine(%q, %q) returned empty string", tc.message, tc.level)
		}

		// Should contain the original message
		if !contains(formatted, tc.message) {
			t.Errorf("FormatLogLine(%q, %q) does not contain original message", tc.message, tc.level)
		}
	}
}

func TestTheme_FormatLogLineWithTime(t *testing.T) {
	theme := NewDefaultTheme()

	now := time.Now()
	message := "Test message"
	level := "info"

	formatted := theme.FormatLogLineWithTime(message, level, now)

	if formatted == "" {
		t.Error("FormatLogLineWithTime() returned empty string")
	}

	// Should contain the original message
	if !contains(formatted, message) {
		t.Error("FormatLogLineWithTime() does not contain original message")
	}

	// Should contain some time representation (at least the hour)
	timeStr := now.Format("15:04")
	if !contains(formatted, timeStr) {
		t.Errorf("FormatLogLineWithTime() does not contain time %q", timeStr)
	}
}

func TestTheme_FormatStatus(t *testing.T) {
	theme := NewDefaultTheme()

	status := "Test status message"
	formatted := theme.FormatStatus(status)

	if formatted == "" {
		t.Error("FormatStatus() returned empty string")
	}

	// Should contain the original status
	if !contains(formatted, status) {
		t.Error("FormatStatus() does not contain original status")
	}
}

func TestTheme_FormatError(t *testing.T) {
	theme := NewDefaultTheme()

	errorMsg := "Test error message"
	formatted := theme.FormatError(errorMsg)

	if formatted == "" {
		t.Error("FormatError() returned empty string")
	}

	// Should contain the original error message
	if !contains(formatted, errorMsg) {
		t.Error("FormatError() does not contain original error message")
	}
}

func TestTheme_FormatSuccess(t *testing.T) {
	theme := NewDefaultTheme()

	successMsg := "Test success message"
	formatted := theme.FormatSuccess(successMsg)

	if formatted == "" {
		t.Error("FormatSuccess() returned empty string")
	}

	// Should contain the original success message
	if !contains(formatted, successMsg) {
		t.Error("FormatSuccess() does not contain original success message")
	}
}

func TestTheme_FormatHelp(t *testing.T) {
	theme := NewDefaultTheme()

	helpText := "Test help text\nWith multiple lines\nAnd commands"
	formatted := theme.FormatHelp(helpText)

	if formatted == "" {
		t.Error("FormatHelp() returned empty string")
	}

	// Should contain the original help text
	if !contains(formatted, "Test help text") {
		t.Error("FormatHelp() does not contain original help text")
	}
}

func TestThemeStyles(t *testing.T) {
	theme := NewDefaultTheme()

	// Test that all styles exist (can't compare lipgloss.Style directly)
	_ = theme.ViewportStyle
	_ = theme.InputStyle
	_ = theme.InputPromptStyle
	_ = theme.StatusStyle
	_ = theme.ErrorStyle
	_ = theme.HelpStyle
}

func TestThemeColors(t *testing.T) {
	theme := NewDefaultTheme()

	// Test that all colors are valid (non-empty)
	colors := []lipgloss.Color{
		theme.PrimaryColor,
		theme.SecondaryColor,
		theme.AccentColor,
		theme.ErrorColor,
		theme.SuccessColor,
		theme.WarningColor,
		theme.MutedColor,
	}

	for i, color := range colors {
		if string(color) == "" {
			t.Errorf("Color %d is empty", i)
		}
	}
}

func TestLogLevelStyleConsistency(t *testing.T) {
	theme := NewDefaultTheme()

	// Test that the same level always returns the same style
	level := "info"
	style1 := theme.GetLogLevelStyle(level)
	style2 := theme.GetLogLevelStyle(level)

	// Compare style properties (we can't directly compare styles)
	if style1.GetForeground() != style2.GetForeground() {
		t.Error("GetLogLevelStyle() returns inconsistent styles for same level")
	}
}

func TestFormatMethodsWithEmptyInput(t *testing.T) {
	theme := NewDefaultTheme()

	// Test formatting methods with empty input - they may return empty strings
	// which is acceptable behavior
	_ = theme.FormatLogLine("", "info")
	_ = theme.FormatStatus("")
	_ = theme.FormatError("")
	_ = theme.FormatSuccess("")
	_ = theme.FormatHelp("")

	// Just test they don't panic
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
