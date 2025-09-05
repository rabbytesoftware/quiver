package stdout

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/rabbytesoftware/quiver/internal/infrastructure/ui/stdout/models"
)

var (
	timestampStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888"))
		
	levelStyle = lipgloss.NewStyle().
		Bold(true)

	messageStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))

	logStyle = lipgloss.NewStyle().
		Padding(0, 1).
		Margin(0, 1)
)

type logView struct {
	entries []models.LogEntry
}

func NewLog() *logView {
	return &logView{}
}

func (l *logView) Add(entry models.LogEntry) {
	l.entries = append(l.entries, entry)
}

func (l *logView) View(width int, height int) string {
	var lines []string
	totalLogs := len(l.entries)
	
	// Calculate which logs to show based on scroll offset
	start := 0
	end := start + height
	if end > totalLogs {
		end = totalLogs
	}

	for i := start; i < end; i++ {
		log := l.entries[i]
		timestamp := timestampStyle.Render(log.Timestamp.Format("15:04:05"))
		
		var levelStyled string
		switch log.Level {
		case "ERROR":
			levelStyled = levelStyle.Foreground(lipgloss.Color("#FF6B6B")).Render(fmt.Sprintf("[%s]", log.Level))
		case "WARN":
			levelStyled = levelStyle.Foreground(lipgloss.Color("#FFE66D")).Render(fmt.Sprintf("[%s]", log.Level))
		case "INFO":
			levelStyled = levelStyle.Foreground(lipgloss.Color("#4ECDC4")).Render(fmt.Sprintf("[%s]", log.Level))
		case "CMD":
			levelStyled = levelStyle.Foreground(lipgloss.Color("#45B7D1")).Render(fmt.Sprintf("[%s]", log.Level))
		default:
			levelStyled = levelStyle.Render(fmt.Sprintf("[%s]", log.Level))
		}
		
		message := messageStyle.Render(log.Message)
		line := fmt.Sprintf("%s %s %s", timestamp, levelStyled, message)
		lines = append(lines, line)
	}

	// Fill remaining lines with empty space
	for len(lines) < height {
		lines = append(lines, "")
	}

	// Add scroll indicator if needed
	scrollInfo := ""
	if len(l.entries) > height {
		scrollInfo = fmt.Sprintf(" [%d/%d]", start+1, totalLogs)
	}

	logContent := strings.Join(lines, "\n")
	
	return logStyle.
		Width(width - 4).
		Height(height).
		Render(logContent + scrollInfo)
}