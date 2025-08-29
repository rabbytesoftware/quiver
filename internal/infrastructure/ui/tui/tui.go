package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp time.Time
	Level     string
	Message   string
}

// ResourceMetrics holds system resource information
type ResourceMetrics struct {
	CPUPercent     float64
	MemoryPercent  float64
	MemoryUsed     uint64
	MemoryTotal    uint64
	ActivePorts    int
}

// Model represents the main TUI model
type Model struct {
	// Layout dimensions
	width  int
	height int

	// Components
	logs         []LogEntry
	logViewport  logViewport
	inputPrompt  inputPrompt
	footerBar    footerBar
	
	// State
	ready        bool
	metrics      ResourceMetrics
	lastUpdate   time.Time
}

// logViewport handles the scrollable log area
type logViewport struct {
	scrollOffset int
	maxLines     int
}

// inputPrompt handles the command input
type inputPrompt struct {
	value  string
	cursor int
	focus  bool
}

// footerBar handles the resource metrics display
type footerBar struct {
	style lipgloss.Style
}

// Styles
var (
	logStyle = lipgloss.NewStyle().
		Padding(0, 1).
		Margin(0, 1)

	promptStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#333333")).
		Padding(0, 1).
		Margin(0, 1)

	footerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#0F4C75")).
		Padding(0, 1).
		Bold(true)

	cpuStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF6B6B")).
		Bold(true)

	memoryStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4ECDC4")).
		Bold(true)

	networkStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#45B7D1")).
		Bold(true)

	timestampStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888"))

	levelStyle = lipgloss.NewStyle().
		Bold(true)

	messageStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))
)

// Messages
type tickMsg time.Time
type resourceUpdateMsg ResourceMetrics

// NewModel creates a new TUI model
func NewModel() Model {
	return Model{
		logs: []LogEntry{
			{Timestamp: time.Now(), Level: "INFO", Message: "Quiver TUI started"},
			{Timestamp: time.Now(), Level: "INFO", Message: "Initializing system..."},
		},
		logViewport: logViewport{},
		inputPrompt: inputPrompt{focus: true},
		footerBar:   footerBar{style: footerStyle},
		lastUpdate:  time.Now(),
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
		updateResourcesCmd(),
	)
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		m.updateLayout()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "enter":
			if m.inputPrompt.value != "" {
				// Add command to logs
				m.addLogEntry("CMD", fmt.Sprintf("> %s", m.inputPrompt.value))
				
				// Process command (placeholder)
				m.addLogEntry("INFO", fmt.Sprintf("Executed: %s", m.inputPrompt.value))
				
				// Clear input
				m.inputPrompt.value = ""
				m.inputPrompt.cursor = 0
			}

		case "up":
			if m.logViewport.scrollOffset > 0 {
				m.logViewport.scrollOffset--
			}

		case "down":
			maxScroll := len(m.logs) - m.logViewport.maxLines
			if maxScroll > 0 && m.logViewport.scrollOffset < maxScroll {
				m.logViewport.scrollOffset++
			}

		case "backspace":
			if m.inputPrompt.cursor > 0 {
				m.inputPrompt.value = m.inputPrompt.value[:m.inputPrompt.cursor-1] + m.inputPrompt.value[m.inputPrompt.cursor:]
				m.inputPrompt.cursor--
			}

		case "left":
			if m.inputPrompt.cursor > 0 {
				m.inputPrompt.cursor--
			}

		case "right":
			if m.inputPrompt.cursor < len(m.inputPrompt.value) {
				m.inputPrompt.cursor++
			}

		default:
			// Add character to input
			if len(msg.String()) == 1 {
				m.inputPrompt.value = m.inputPrompt.value[:m.inputPrompt.cursor] + msg.String() + m.inputPrompt.value[m.inputPrompt.cursor:]
				m.inputPrompt.cursor++
			}
		}

	case tickMsg:
		cmds = append(cmds, tickCmd())

	case resourceUpdateMsg:
		m.metrics = ResourceMetrics(msg)
		cmds = append(cmds, updateResourcesCmd())
	}

	return m, tea.Batch(cmds...)
}

// View renders the TUI
func (m Model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	// Calculate layout dimensions
	footerHeight := 1
	promptHeight := 1
	logHeight := m.height - footerHeight - promptHeight - 2 // -2 for borders/spacing

	// Render log viewport
	logView := m.renderLogView(logHeight)

	// Render input prompt
	promptView := m.renderInputPrompt()

	// Render footer
	footerView := m.renderFooter()

	// Combine all components
	return lipgloss.JoinVertical(
		lipgloss.Left,
		logView,
		promptView,
		footerView,
	)
}

// renderLogView renders the scrollable log area
func (m Model) renderLogView(height int) string {
	m.logViewport.maxLines = height

	var lines []string
	totalLogs := len(m.logs)
	
	// Calculate which logs to show based on scroll offset
	start := m.logViewport.scrollOffset
	end := start + height
	if end > totalLogs {
		end = totalLogs
	}

	for i := start; i < end; i++ {
		log := m.logs[i]
		timestamp := timestampStyle.Render(log.Timestamp.Format("15:04:05"))
		
		var levelStyled string
		switch log.Level {
		case "ERROR":
			levelStyled = levelStyle.Foreground(lipgloss.Color("#FF6B6B")).Render("[" + log.Level + "]")
		case "WARN":
			levelStyled = levelStyle.Foreground(lipgloss.Color("#FFE66D")).Render("[" + log.Level + "]")
		case "INFO":
			levelStyled = levelStyle.Foreground(lipgloss.Color("#4ECDC4")).Render("[" + log.Level + "]")
		case "CMD":
			levelStyled = levelStyle.Foreground(lipgloss.Color("#45B7D1")).Render("[" + log.Level + "]")
		default:
			levelStyled = levelStyle.Render("[" + log.Level + "]")
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
	if len(m.logs) > height {
		scrollInfo = fmt.Sprintf(" [%d/%d]", start+1, totalLogs)
	}

	logContent := strings.Join(lines, "\n")
	
	return logStyle.
		Width(m.width - 4).
		Height(height).
		Render(logContent + scrollInfo)
}

// renderInputPrompt renders the bash-style input prompt
func (m Model) renderInputPrompt() string {
	prompt := "quiver> "
	
	// Add cursor
	value := m.inputPrompt.value
	if m.inputPrompt.focus {
		if m.inputPrompt.cursor >= len(value) {
			value += "█"
		} else {
			value = value[:m.inputPrompt.cursor] + "█" + value[m.inputPrompt.cursor+1:]
		}
	}

	promptText := prompt + value
	
	return promptStyle.
		Width(m.width - 4).
		Render(promptText)
}

// renderFooter renders the resource usage footer
func (m Model) renderFooter() string {
	cpu := cpuStyle.Render(fmt.Sprintf("CPU: %.1f%%", m.metrics.CPUPercent))
	memory := memoryStyle.Render(fmt.Sprintf("RAM: %.1f%% (%.1fGB/%.1fGB)", 
		m.metrics.MemoryPercent, 
		float64(m.metrics.MemoryUsed)/1024/1024/1024,
		float64(m.metrics.MemoryTotal)/1024/1024/1024))
	network := networkStyle.Render(fmt.Sprintf("Ports: %d", m.metrics.ActivePorts))
	
	footerContent := fmt.Sprintf("%s  %s  %s", cpu, memory, network)
	
	return footerStyle.
		Width(m.width).
		Render(footerContent)
}

// updateLayout updates component layouts based on window size
func (m *Model) updateLayout() {
	// Layout is handled dynamically in View()
}

// addLogEntry adds a new log entry
func (m *Model) addLogEntry(level, message string) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
	}
	m.logs = append(m.logs, entry)

	// Auto-scroll to bottom when new entries are added
	if m.logViewport.maxLines > 0 {
		maxScroll := len(m.logs) - m.logViewport.maxLines
		if maxScroll > 0 {
			m.logViewport.scrollOffset = maxScroll
		}
	}
}

// Commands
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func updateResourcesCmd() tea.Cmd {
	return func() tea.Msg {
		// Get CPU usage
		cpuPercents, err := cpu.Percent(0, false)
		var cpuUsage float64
		if err == nil && len(cpuPercents) > 0 {
			cpuUsage = cpuPercents[0]
		}

		// Get memory usage
		memInfo, err := mem.VirtualMemory()
		var memUsage, memUsed, memTotal float64
		if err == nil {
			memUsage = memInfo.UsedPercent
			memUsed = float64(memInfo.Used)
			memTotal = float64(memInfo.Total)
		}

		// Get network connections (as proxy for active ports)
		connections, err := net.Connections("all")
		activePorts := 0
		if err == nil {
			activePorts = len(connections)
		}

		return resourceUpdateMsg(ResourceMetrics{
			CPUPercent:    cpuUsage,
			MemoryPercent: memUsage,
			MemoryUsed:    uint64(memUsed),
			MemoryTotal:   uint64(memTotal),
			ActivePorts:   activePorts,
		})
	}
}

// RunTUI starts the TUI application
func RunTUI() error {
	p := tea.NewProgram(
		NewModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	_, err := p.Run()
	return err
}