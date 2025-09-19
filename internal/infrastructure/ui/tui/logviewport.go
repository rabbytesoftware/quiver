package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type LogEntry struct {
	Timestamp time.Time
	Message   string
}

type LogViewport struct {
	entries      []LogEntry
	scrollOffset int
	maxLines     int
	width        int
	height       int
}

func NewLogViewport() *LogViewport {
	return &LogViewport{
		entries:      make([]LogEntry, 0),
		scrollOffset: 0,
		maxLines:     0,
	}
}

func (lv *LogViewport) Add(message string) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Message:   message,
	}
	lv.entries = append(lv.entries, entry)
	
	// ? Auto-scroll to bottom when new entries are added
	if lv.maxLines > 0 {
		maxScroll := len(lv.entries) - lv.maxLines
		if maxScroll > 0 {
			lv.scrollOffset = maxScroll
		}
	}
}

func (lv *LogViewport) View(width, height int) string {
	lv.width = width
	lv.height = height
	lv.maxLines = height

	var lines []string
	totalLogs := len(lv.entries)
	
	// ? Calculate which logs to show based on scroll offset
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

	// ? Fill remaining lines with empty space
	for len(lines) < height {
		lines = append(lines, "")
	}

	// ? Add scroll indicator if needed
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

func (lv *LogViewport) ScrollUp() {
	if lv.scrollOffset > 0 {
		lv.scrollOffset--
	}
}

func (lv *LogViewport) ScrollDown() {
	maxScroll := len(lv.entries) - lv.maxLines
	if maxScroll > 0 && lv.scrollOffset < maxScroll {
		lv.scrollOffset++
	}
}

func (lv *LogViewport) ScrollToBottom() {
	maxScroll := len(lv.entries) - lv.maxLines
	if maxScroll > 0 {
		lv.scrollOffset = maxScroll
	} else {
		lv.scrollOffset = 0
	}
}

func (lv *LogViewport) GetEntryCount() int {
	return len(lv.entries)
}

func (lv *LogViewport) formatLogEntry(entry LogEntry) string {
	timestamp := lv.getTimestampStyle().Render(entry.Timestamp.Format("15:04:05"))
	
	levelStyled := lv.getLevelStyle().Foreground(lipgloss.Color("#4ECDC4")).Render("[INFO]")
	
	message := lv.getMessageStyle().Render(entry.Message)
	return fmt.Sprintf("%s %s %s", timestamp, levelStyled, message)
}

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
