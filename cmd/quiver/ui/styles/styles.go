package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme contains all the styling for the application
type Theme struct {
	// Log level colors
	DebugStyle lipgloss.Style
	InfoStyle  lipgloss.Style
	WarnStyle  lipgloss.Style
	ErrorStyle lipgloss.Style
	
	// UI components
	ViewportStyle    lipgloss.Style
	InputPromptStyle lipgloss.Style
	InputStyle       lipgloss.Style
	StatusStyle      lipgloss.Style
	HelpStyle        lipgloss.Style
	
	// Colors
	PrimaryColor   lipgloss.Color
	SecondaryColor lipgloss.Color
	AccentColor    lipgloss.Color
	ErrorColor     lipgloss.Color
	WarningColor   lipgloss.Color
	SuccessColor   lipgloss.Color
	MutedColor     lipgloss.Color
}

// NewDefaultTheme creates a new default theme
func NewDefaultTheme() Theme {
	// Define color palette
	primaryColor := lipgloss.Color("#7C3AED")    // Purple
	secondaryColor := lipgloss.Color("#64748B")  // Slate
	accentColor := lipgloss.Color("#06B6D4")     // Cyan
	errorColor := lipgloss.Color("#EF4444")      // Red
	warningColor := lipgloss.Color("#F59E0B")    // Amber
	successColor := lipgloss.Color("#10B981")    // Emerald
	mutedColor := lipgloss.Color("#94A3B8")      // Slate light

	return Theme{
		// Log level styles
		DebugStyle: lipgloss.NewStyle().
			Foreground(mutedColor).
			Faint(true),
		InfoStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E2E8F0")), // Light slate
		WarnStyle: lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true),
		ErrorStyle: lipgloss.NewStyle().
			Foreground(errorColor).
			Bold(true),
		
		// UI component styles
		ViewportStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Padding(0, 1),
		InputPromptStyle: lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true),
		InputStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F1F5F9")), // Very light slate
		StatusStyle: lipgloss.NewStyle().
			Foreground(accentColor).
			Italic(true).
			Padding(0, 1),
		HelpStyle: lipgloss.NewStyle().
			Foreground(mutedColor).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor),
		
		// Colors for direct use
		PrimaryColor:   primaryColor,
		SecondaryColor: secondaryColor,
		AccentColor:    accentColor,
		ErrorColor:     errorColor,
		WarningColor:   warningColor,
		SuccessColor:   successColor,
		MutedColor:     mutedColor,
	}
}

// GetLogLevelStyle returns the appropriate style for a log level
func (t Theme) GetLogLevelStyle(level string) lipgloss.Style {
	switch level {
	case "debug":
		return t.DebugStyle
	case "info":
		return t.InfoStyle
	case "warn":
		return t.WarnStyle
	case "error":
		return t.ErrorStyle
	default:
		return t.InfoStyle
	}
}

// FormatLogLine formats a log line with the appropriate style
func (t Theme) FormatLogLine(text, level string) string {
	style := t.GetLogLevelStyle(level)
	return style.Render(text)
}

// FormatStatus formats a status message
func (t Theme) FormatStatus(text string) string {
	return t.StatusStyle.Render(text)
}

// FormatError formats an error message
func (t Theme) FormatError(text string) string {
	return lipgloss.NewStyle().
		Foreground(t.ErrorColor).
		Bold(true).
		Render("Error: " + text)
}

// FormatSuccess formats a success message
func (t Theme) FormatSuccess(text string) string {
	return lipgloss.NewStyle().
		Foreground(t.SuccessColor).
		Bold(true).
		Render(text)
}

// FormatHelp formats help text
func (t Theme) FormatHelp(text string) string {
	return t.HelpStyle.Render(text)
}
