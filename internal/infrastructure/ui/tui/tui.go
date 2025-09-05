package tui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/rabbytesoftware/quiver/internal/infrastructure/ui/cmd/interpreter"
	"github.com/rabbytesoftware/quiver/internal/infrastructure/ui/stdout"
	"github.com/rabbytesoftware/quiver/internal/infrastructure/ui/stdout/models"
)

// Model represents the main TUI model
type Model struct {
	// Layout dimensions
	width  int
	height int

	// Components
	logViewport    *LogViewport
	inputPrompt    *InputPrompt
	footerBar      *FooterBar
	metricsCollector *SystemMetricsCollector
	commandInterpreter *interpreter.CommandInterpreter
	logService     *stdout.LogService
	
	// State
	ready        bool
	lastUpdate   time.Time
}

// Messages
type tickMsg time.Time
type resourceUpdateMsg ResourceMetrics

// NewModel creates a new TUI model
func NewModel(logService *stdout.LogService) Model {
	logViewport := NewLogViewport()
	inputPrompt := NewInputPrompt()
	footerBar := NewFooterBar()
	metricsCollector := NewSystemMetricsCollector()
	commandInterpreter := interpreter.NewCommandInterpreter()
	
	// Set up log service subscription to update the viewport
	logService.Subscribe(func(entry models.LogEntry) {
		logViewport.Add(entry)
	})
	
	// Add initial log entries through the log service
	logService.Info("Quiver TUI started")
	logService.Info("Initializing system...")
	
	return Model{
		logViewport:        logViewport,
		inputPrompt:        inputPrompt,
		footerBar:          footerBar,
		metricsCollector:   metricsCollector,
		commandInterpreter: commandInterpreter,
		logService:         logService,
		lastUpdate:         time.Now(),
	}
}

func (m Model) Init() tea.Cmd {
	m.metricsCollector.Start()
	
	return tea.Batch(
		tickCmd(),
		m.listenForMetricsUpdates(),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m = m.handleWindowSize(msg)

	case tea.KeyMsg:
		return m.handleKeyInput(msg)

	case tickMsg:
		cmds = append(cmds, tickCmd())

	case resourceUpdateMsg:
		m.footerBar.UpdateMetrics(ResourceMetrics(msg))
		cmds = append(cmds, m.listenForMetricsUpdates())
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

	// Render components
	logView := m.logViewport.View(m.width, logHeight)
	promptView := m.inputPrompt.View(m.width)
	footerView := m.footerBar.View(m.width)

	// Combine all components
	return lipgloss.JoinVertical(
		lipgloss.Left,
		logView,
		promptView,
		footerView,
	)
}

// listenForMetricsUpdates listens for metrics updates from the collector
func (m Model) listenForMetricsUpdates() tea.Cmd {
	return func() tea.Msg {
		metrics := <-m.metricsCollector.GetUpdateChannel()
		return resourceUpdateMsg(metrics)
	}
}

// handleWindowSize handles window size changes
func (m Model) handleWindowSize(msg tea.WindowSizeMsg) Model {
	m.width = msg.Width
	m.height = msg.Height
	m.ready = true
	return m
}

// handleKeyInput handles keyboard input
func (m Model) handleKeyInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		return m, tea.Quit
	case "enter":
		return m.handleEnterKey()
	case "up":
		m.logViewport.ScrollUp()
		return m, nil
	case "down":
		m.logViewport.ScrollDown()
		return m, nil
	default:
		m.inputPrompt.HandleKey(msg.String())
		return m, nil
	}
}

// handleEnterKey processes the enter key press
func (m Model) handleEnterKey() (tea.Model, tea.Cmd) {
	if m.inputPrompt.GetValue() == "" {
		return m, nil
	}
	
	m = m.processCommand()
	return m, nil
}

// processCommand processes the current command input
func (m Model) processCommand() Model {
	command := m.inputPrompt.GetValue()
	
	// Log the command
	m.logService.Info(fmt.Sprintf("> %s", command))
	
	// Process the command
	result, err := m.commandInterpreter.ProcessCommand(command)
	if err != nil {
		m.logService.Error(fmt.Sprintf("Error: %s", err.Error()))
	} else if result != "" {
		// Handle special commands
		if result == "CLEAR_SCREEN" {
			m = m.ClearScreen()
		} else {
			m.logService.Info(result)
		}
	}
	
	// Clear input
	m.inputPrompt.Clear()
	return m
}

func (m Model) ClearScreen() Model {
	m.logViewport = NewLogViewport()
	return m
}

// Commands
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func RunTUI() (*stdout.LogService, error) {
	logService := stdout.NewLogService()
	
	p := tea.NewProgram(
		NewModel(logService),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	_, err := p.Run()
	return logService, err
}
