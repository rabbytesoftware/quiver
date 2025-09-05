package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rabbytesoftware/quiver/internal/infrastructure/ui/stdout/models"
)

// LogViewport handles the scrollable log area
type LogViewport struct {
	entries      []models.LogEntry
	scrollOffset int
	maxLines     int
	width        int
	height       int
}

// NewLogViewport creates a new log viewport
func NewLogViewport() *LogViewport {
	return &LogViewport{
		entries:      make([]models.LogEntry, 0),
		scrollOffset: 0,
		maxLines:     0,
	}
}

// Add adds a new log entry
func (lv *LogViewport) Add(entry models.LogEntry) {
	lv.entries = append(lv.entries, entry)
	
	// Auto-scroll to bottom when new entries are added
	if lv.maxLines > 0 {
		maxScroll := len(lv.entries) - lv.maxLines
		if maxScroll > 0 {
			lv.scrollOffset = maxScroll
		}
	}
}

// View renders the log viewport
func (lv *LogViewport) View(width, height int) string {
	lv.width = width
	lv.height = height
	lv.maxLines = height

	var lines []string
	totalLogs := len(lv.entries)
	
	// Calculate which logs to show based on scroll offset
	start := lv.scrollOffset
	end := start + height
	if end > totalLogs {
		end = totalLogs
	}

	for i := start; i < end; i++ {
		log := lv.entries[i]
		line := lv.formatLogEntry(log)
		lines = append(lines, line)
	}

	// Fill remaining lines with empty space
	for len(lines) < height {
		lines = append(lines, "")
	}

	// Add scroll indicator if needed
	scrollInfo := ""
	if len(lv.entries) > height {
		scrollInfo = fmt.Sprintf(" [%d/%d]", start+1, totalLogs)
	}

	logContent := strings.Join(lines, "\n")
	
	return lv.getLogStyle().
		Width(width - 4).
		Height(height).
		Render(logContent + scrollInfo)
}

// ScrollUp scrolls the viewport up
func (lv *LogViewport) ScrollUp() {
	if lv.scrollOffset > 0 {
		lv.scrollOffset--
	}
}

// ScrollDown scrolls the viewport down
func (lv *LogViewport) ScrollDown() {
	maxScroll := len(lv.entries) - lv.maxLines
	if maxScroll > 0 && lv.scrollOffset < maxScroll {
		lv.scrollOffset++
	}
}

// ScrollToBottom scrolls to the bottom
func (lv *LogViewport) ScrollToBottom() {
	maxScroll := len(lv.entries) - lv.maxLines
	if maxScroll > 0 {
		lv.scrollOffset = maxScroll
	} else {
		lv.scrollOffset = 0
	}
}

// GetEntryCount returns the number of log entries
func (lv *LogViewport) GetEntryCount() int {
	return len(lv.entries)
}

// formatLogEntry formats a single log entry for display
func (lv *LogViewport) formatLogEntry(entry models.LogEntry) string {
	timestamp := lv.getTimestampStyle().Render(entry.Timestamp.Format("15:04:05"))
	
	var levelStyled string
	switch entry.Level {
	case models.LogLevelError:
		levelStyled = lv.getLevelStyle().Foreground(lipgloss.Color("#FF6B6B")).Render(fmt.Sprintf("[%s]", entry.Level))
	case models.LogLevelWarning:
		levelStyled = lv.getLevelStyle().Foreground(lipgloss.Color("#FFE66D")).Render(fmt.Sprintf("[%s]", entry.Level))
	case models.LogLevelInfo:
		levelStyled = lv.getLevelStyle().Foreground(lipgloss.Color("#4ECDC4")).Render(fmt.Sprintf("[%s]", entry.Level))
	case models.LogLevelDebug:
		levelStyled = lv.getLevelStyle().Foreground(lipgloss.Color("#45B7D1")).Render(fmt.Sprintf("[%s]", entry.Level))
	default:
		levelStyled = lv.getLevelStyle().Render(fmt.Sprintf("[%s]", entry.Level))
	}
	
	message := lv.getMessageStyle().Render(entry.Message)
	return fmt.Sprintf("%s %s %s", timestamp, levelStyled, message)
}

// Styles
func (lv *LogViewport) getLogStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Padding(0, 1).
		Margin(0, 1)
}

func (lv *LogViewport) getTimestampStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888"))
}

func (lv *LogViewport) getLevelStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true)
}

func (lv *LogViewport) getMessageStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))
}
