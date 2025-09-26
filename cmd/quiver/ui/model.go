package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sirupsen/logrus"

	"github.com/rabbytesoftware/quiver/cmd/quiver/ui/domain/commands"
	"github.com/rabbytesoftware/quiver/cmd/quiver/ui/domain/events"
	"github.com/rabbytesoftware/quiver/cmd/quiver/ui/domain/handlers"
	"github.com/rabbytesoftware/quiver/cmd/quiver/ui/queries"
	"github.com/rabbytesoftware/quiver/cmd/quiver/ui/services"
	"github.com/rabbytesoftware/quiver/cmd/quiver/ui/styles"

	"github.com/rabbytesoftware/quiver/internal/core/config"
	"github.com/rabbytesoftware/quiver/internal/core/metadata"
	"github.com/rabbytesoftware/quiver/internal/core/watcher"
)

const (
	maxLogLines = 10000 // Ring buffer size for log lines
)

// Model represents the Bubble Tea model for the TUI
type Model struct {
	// UI components
	viewport  viewport.Model
	textInput textinput.Model
	
	// Application state
	logLines       []string
	autoScroll     bool
	status         string
	statusExpiry   time.Time
	ready          bool
	quitting       bool
	
	// Command history
	commandHistory    []string
	historyIndex      int
	currentInput      string
	navigatingHistory bool
	
	// Services and handlers
	watcher        *watcher.Watcher
	watcherAdapter *services.WatcherAdapter
	queryService   *queries.QueryService
	handler        *handlers.Handler
	theme          styles.Theme
	
	// Context and cancellation
	ctx    context.Context
	cancel context.CancelFunc
	
	// Dimensions
	width  int
	height int
}

// NewModel creates a new TUI model with the provided watcher
func NewModel(w *watcher.Watcher) *Model {
	watcherAdapter := services.NewWatcherAdapter(w)
	
	queryService := queries.NewService(fmt.Sprintf("http://%s:%d", config.GetAPI().Host, config.GetAPI().Port))
	
	handler := handlers.NewHandler(watcherAdapter, queryService)
	theme := styles.NewDefaultTheme()
	
	ctx, cancel := context.WithCancel(context.Background())
	
	ti := textinput.New()
	ti.Placeholder = "Enter a command (e.g., help)"
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 50
	
	vp := viewport.New(80, 20)
	vp.SetContent("")
	
	model := &Model{
		viewport:          vp,
		textInput:         ti,
		logLines:          make([]string, 0, maxLogLines),
		autoScroll:        true,
		theme:             theme,
		watcher:           w,
		watcherAdapter:    watcherAdapter,
		queryService:      queryService,
		handler:           handler,
		ctx:               ctx,
		cancel:            cancel,
		commandHistory:    make([]string, 0),
		historyIndex:      -1,
		currentInput:      "",
		navigatingHistory: false,
	}
	
	// Subscribe to watcher for log messages
	model.subscribeToWatcher()
	
	return model
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		m.tickStatus(),
	)
}

// Update handles messages and updates the model
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)
	
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		// Update viewport size (leave space for input and status)
		viewportHeight := msg.Height - 3 // 1 for input, 1 for status, 1 for border
		m.viewport.Width = msg.Width - 2  // Account for border
		m.viewport.Height = viewportHeight
		
		// Update text input width
		m.textInput.Width = msg.Width - 4 // Account for prompt "> "
		
		if !m.ready {
			m.ready = true
		}
		
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			m.cancel()
			return m, tea.Quit
			
		case "enter":
			return m.handleCommand()
			
		case "esc":
			m.textInput.SetValue("")
			m.resetHistoryNavigation()
			
		case "up":
			// If text input is focused, navigate command history
			if m.textInput.Focused() {
				return m.navigateHistory(-1), nil
			}
			// Otherwise, handle viewport scrolling
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
			m.autoScroll = m.viewport.AtTop()
			
		case "down":
			// If text input is focused, navigate command history
			if m.textInput.Focused() {
				return m.navigateHistory(1), nil
			}
			// Otherwise, handle viewport scrolling
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
			m.autoScroll = m.viewport.AtTop()
			
		case "pgup", "pgdown":
			// Handle viewport scrolling
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
			m.autoScroll = m.viewport.AtTop()
		}
		
	case events.LogLineReceivedMsg:
		m.addLogLine(msg.Event.LogLine)
		
	case events.FilterAppliedMsg:
		if msg.Event.Pattern == "" {
			m.setStatus("Filter cleared", 3*time.Second)
		} else {
			m.setStatus(fmt.Sprintf("Filter applied: %s", msg.Event.Pattern), 3*time.Second)
		}
		
	case events.LevelChangedMsg:
		m.setStatus(fmt.Sprintf("Log level set to: %s", msg.Event.Level), 3*time.Second)
		
	case events.StreamPausedMsg:
		m.setStatus("Log streaming paused", 3*time.Second)
		
	case events.StreamResumedMsg:
		m.setStatus("Log streaming resumed", 3*time.Second)
		
	case events.ClearedMsg:
		m.clearLogs()
		m.setStatus("Viewport cleared", 2*time.Second)
		
	case events.CommandErrorMsg:
		m.setStatus(m.theme.FormatError(msg.Event.Message), 5*time.Second)
		
	case events.HelpRequestedMsg:
		m.showHelp(msg.Event.HelpText)
		
	case events.QueryExecutedMsg:
		m.showQueryResult(msg.Event.DisplayText)
		
	case statusTickMsg:
		// Clear expired status messages
		if time.Now().After(m.statusExpiry) {
			m.status = ""
		}
		return m, m.tickStatus()
	}
	
	// Update text input
	m.textInput, cmd = m.textInput.Update(msg)
	cmds = append(cmds, cmd)
	
	return m, tea.Batch(cmds...)
}

// handleCommand processes the entered command
func (m *Model) handleCommand() (tea.Model, tea.Cmd) {
	input := strings.TrimSpace(m.textInput.Value())
	m.textInput.SetValue("")
	m.resetHistoryNavigation()
	
	if input == "" {
		return m, nil
	}
	
	// Add to command history (avoid duplicates)
	if len(m.commandHistory) == 0 || m.commandHistory[len(m.commandHistory)-1] != input {
		m.commandHistory = append(m.commandHistory, input)
		// Keep history size manageable
		if len(m.commandHistory) > 100 {
			m.commandHistory = m.commandHistory[1:]
		}
	}
	
	// Handle quit command
	if input == "quit" || input == "q" || input == "exit" {
		m.quitting = true
		m.cancel()
		return m, tea.Quit
	}
	
	// Parse command
	cmd, err := commands.Parse(input)
	if err != nil {
		m.setStatus(m.theme.FormatError(err.Error()), 5*time.Second)
		return m, nil
	}
	
	// Handle command
	resultEvents := m.handler.Handle(cmd)
	
	// Convert events to tea messages and send them
	var cmds []tea.Cmd
	for _, event := range resultEvents {
		if teaMsg := events.ToTeaMsg(event); teaMsg != nil {
			cmds = append(cmds, func(msg tea.Msg) tea.Cmd {
				return func() tea.Msg { return msg }
			}(teaMsg))
		}
	}
	
	return m, tea.Batch(cmds...)
}

// addLogLine adds a new log line to the viewport
func (m *Model) addLogLine(logLine events.LogLine) {
	// Format the log line with timestamp and gray styling
	formattedLine := m.theme.FormatLogLineWithTime(logLine.Text, logLine.Level, logLine.Time)
	
	// Add to the beginning of the slice for newest-first display
	m.logLines = append([]string{formattedLine}, m.logLines...)
	if len(m.logLines) > maxLogLines {
		m.logLines = m.logLines[:maxLogLines]
	}
	
	// Update viewport content
	m.updateViewportContent()
}

// updateViewportContent updates the viewport with current log lines
func (m *Model) updateViewportContent() {
	content := strings.Join(m.logLines, "\n")
	m.viewport.SetContent(content)
	
	// Auto-scroll to top for newest messages if enabled
	if m.autoScroll {
		m.viewport.GotoTop()
	}
}

// clearLogs clears all log lines
func (m *Model) clearLogs() {
	m.logLines = make([]string, 0, maxLogLines)
	m.viewport.SetContent("")
	m.autoScroll = true
}

// showHelp displays help text in the viewport
func (m *Model) showHelp(helpText string) {
	// Add help text as a special log entry
	formattedHelp := m.theme.FormatHelp(helpText)
	lines := strings.Split(formattedHelp, "\n")
	
	// Reverse lines to maintain newest-first order when prepending
	for i := len(lines) - 1; i >= 0; i-- {
		m.logLines = append([]string{lines[i]}, m.logLines...)
		if len(m.logLines) > maxLogLines {
			m.logLines = m.logLines[:maxLogLines]
		}
	}
	
	m.updateViewportContent()
}

// showQueryResult displays query result text in the viewport
func (m *Model) showQueryResult(resultText string) {
	// Add query result as a special log entry
	formattedResult := m.theme.FormatHelp(resultText)
	lines := strings.Split(formattedResult, "\n")
	
	// Reverse lines to maintain newest-first order when prepending
	for i := len(lines) - 1; i >= 0; i-- {
		m.logLines = append([]string{lines[i]}, m.logLines...)
		if len(m.logLines) > maxLogLines {
			m.logLines = m.logLines[:maxLogLines]
		}
	}
	
	m.updateViewportContent()
}

// setStatus sets a status message with expiry
func (m *Model) setStatus(message string, duration time.Duration) {
	m.status = message
	m.statusExpiry = time.Now().Add(duration)
}

// statusTickMsg is used for status message expiry
type statusTickMsg time.Time

// tickStatus returns a command that ticks for status updates
func (m *Model) tickStatus() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return statusTickMsg(t)
	})
}

// subscribeToWatcher sets up the subscription to watcher for log messages
func (m *Model) subscribeToWatcher() {
	m.watcherAdapter.Subscribe(func(level logrus.Level, message string) {
		// Convert logrus level to string
		levelStr := level.String()
		
		// Create log line event
		logLine := events.LogLine{
			Text:  message,
			Level: levelStr,
			Time:  time.Now(),
		}
		
		// Send the event through a goroutine to avoid blocking
		go func() {
			select {
			case <-m.ctx.Done():
				return
			default:
				// This creates a Bubble Tea command to send the message
				// We need to use the program's Send method, but since we don't have access to it here,
				// we'll store the log line and let the UI poll for it
				m.addLogLineFromWatcher(logLine)
			}
		}()
	})
}

// addLogLineFromWatcher adds a log line from the watcher subscription
func (m *Model) addLogLineFromWatcher(logLine events.LogLine) {
	// Format the log line with timestamp and gray styling
	formattedLine := m.theme.FormatLogLineWithTime(logLine.Text, logLine.Level, logLine.Time)
	
	// Add to the beginning of the slice for newest-first display
	m.logLines = append([]string{formattedLine}, m.logLines...)
	if len(m.logLines) > maxLogLines {
		m.logLines = m.logLines[:maxLogLines]
	}
	
	// Update viewport content
	m.updateViewportContent()
}

// View renders the model
func (m *Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	
	if m.quitting {
		return "\n  Goodbye!\n"
	}
	
	// Main viewport
	viewportView := m.theme.ViewportStyle.Render(m.viewport.View())
	
	// Status line
	statusView := ""
	if m.status != "" {
		statusView = m.theme.FormatStatus(m.status)
	}
	
	// Input line with metadata name and version
	promptText := fmt.Sprintf("%s:%s ", metadata.GetName(), metadata.GetVersion())
	prompt := m.theme.InputPromptStyle.Render(promptText)
	inputView := prompt + m.theme.InputStyle.Render(m.textInput.View())
	
	// Combine all parts
	var parts []string
	parts = append(parts, viewportView)
	
	if statusView != "" {
		parts = append(parts, statusView)
	}
	
	parts = append(parts, inputView)
	
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

// navigateHistory navigates through command history
func (m *Model) navigateHistory(direction int) tea.Model {
	if len(m.commandHistory) == 0 {
		return m
	}

	// Save current input if not already navigating
	if !m.navigatingHistory {
		m.currentInput = m.textInput.Value()
		m.navigatingHistory = true
		m.historyIndex = len(m.commandHistory) // Start at end + 1
	}

	// Navigate history
	m.historyIndex += direction

	// Bounds checking
	if m.historyIndex < 0 {
		m.historyIndex = 0
	} else if m.historyIndex >= len(m.commandHistory) {
		// Beyond history - restore current input
		m.historyIndex = len(m.commandHistory)
		m.textInput.SetValue(m.currentInput)
		return m
	}

	// Set the history item
	m.textInput.SetValue(m.commandHistory[m.historyIndex])
	return m
}

// resetHistoryNavigation resets history navigation state
func (m *Model) resetHistoryNavigation() {
	m.navigatingHistory = false
	m.historyIndex = -1
	m.currentInput = ""
}
