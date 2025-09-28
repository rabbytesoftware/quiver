package ui

import (
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/rabbytesoftware/quiver/cmd/quiver/ui/domain/events"
	"github.com/rabbytesoftware/quiver/internal/core/watcher"
)

func TestNewModel(t *testing.T) {
	watcher := watcher.NewWatcherService()

	model := NewModel(watcher)

	if model == nil {
		t.Fatal("Expected model to be created")
	}

	if model.watcher != watcher {
		t.Error("Expected watcher to be set")
	}

	if model.handler == nil {
		t.Error("Expected handler to be created")
	}

	// Theme is a struct, not a pointer, so we can't check for nil
	// Just verify it exists by accessing a method
	_ = model.theme

	if model.queryService == nil {
		t.Error("Expected query service to be created")
	}

	if model.watcherAdapter == nil {
		t.Error("Expected watcher adapter to be created")
	}

	// Test ASCII service initialization
	if model.asciiService == nil {
		t.Error("Expected ASCII service to be created")
	}
}

func TestModel_ASCIIArtDisplay(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	// Test that ASCII art is added as a log line
	if len(model.logLines) == 0 {
		t.Error("Expected ASCII art to be added as initial log line")
	}

	// The first log line should contain the ASCII art
	firstLogLine := model.logLines[0]
	if firstLogLine == "" {
		t.Error("Expected first log line to contain ASCII art")
	}

	// Test that the ASCII service works correctly
	welcomeLogLine := model.asciiService.GetWelcomeLogLine()
	if welcomeLogLine.Text == "" {
		t.Error("Expected ASCII service to return ASCII art text")
	}
}

func TestModel_Init(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	cmd := model.Init()

	// Init should return a command (likely a batch command)
	if cmd == nil {
		t.Error("Expected Init to return a command")
	}
}

func TestModel_TickStatus(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	// Test that tickStatus returns a command
	cmd := model.tickStatus()
	if cmd == nil {
		t.Error("Expected tickStatus to return a command")
	}
}

func TestModel_SetStatus(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	// Test setting status
	model.setStatus("Test status", 2*time.Second)

	// We can't easily verify the internal state without exposing it,
	// but we can at least ensure the method doesn't panic
}

func TestModel_ClearLogs(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	// Test clearing logs
	model.clearLogs()

	// We can't easily verify the internal state without exposing it,
	// but we can at least ensure the method doesn't panic
}

func TestModel_ResetHistoryNavigation(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	// Test resetting history navigation
	model.resetHistoryNavigation()

	// We can't easily verify the internal state without exposing it,
	// but we can at least ensure the method doesn't panic
}

func TestModel_AddLogLine(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	// Test adding log lines
	logLine1 := events.LogLine{
		Text:  "Test log line 1",
		Level: "info",
		Time:  time.Now(),
	}
	logLine2 := events.LogLine{
		Text:  "Test log line 2",
		Level: "error",
		Time:  time.Now(),
	}

	model.addLogLine(logLine1)
	model.addLogLine(logLine2)

	// Method should not panic
}

func TestModel_UpdateViewportContent(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	// Add some log lines first
	logLine1 := events.LogLine{Text: "Line 1", Level: "info", Time: time.Now()}
	logLine2 := events.LogLine{Text: "Line 2", Level: "debug", Time: time.Now()}
	model.addLogLine(logLine1)
	model.addLogLine(logLine2)

	// Test updating viewport content
	model.updateViewportContent()

	// Method should not panic
}

func TestModel_ShowHelp(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	// Test showing help
	model.showHelp("Test help content")

	// Method should not panic
}

func TestModel_ShowQueryResult(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	// Test showing query result
	model.showQueryResult("test query", 200, "success response")

	// Method should not panic
}

func TestModel_ShowQueryError(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	// Test showing query error
	model.showQueryError("test query", 404, "not found")

	// Method should not panic
}

func TestModel_AddLogLineFromWatcher(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	// Test adding log line from watcher
	logLine := events.LogLine{
		Text:  "Test message from watcher",
		Level: "info",
		Time:  time.Now(),
	}
	model.addLogLineFromWatcher(logLine)

	// Method should not panic
}

func TestModel_SubscribeToWatcher(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	// Test subscribing to watcher (method doesn't return anything)
	model.subscribeToWatcher()

	// Method should not panic
}

func TestModel_NavigateHistory(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	// Add some command history first
	model.commandHistory = []string{"command1", "command2", "command3"}

	// Test navigating history up
	result := model.navigateHistory(-1)
	if result == nil {
		t.Error("Expected navigateHistory to return a model")
	}

	// Test navigating history down
	result = model.navigateHistory(1)
	if result == nil {
		t.Error("Expected navigateHistory to return a model")
	}
}

func TestModel_Update_WindowSize(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	// Test window size message
	msg := tea.WindowSizeMsg{
		Width:  100,
		Height: 50,
	}

	updatedModel, cmd := model.Update(msg)

	if updatedModel == nil {
		t.Error("Expected Update to return a model")
	}

	// cmd can be nil or a command, both are valid
	_ = cmd
}

func TestModel_Update_KeyMsg(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	// Test escape key
	escMsg := tea.KeyMsg{
		Type: tea.KeyEsc,
	}

	updatedModel, cmd := model.Update(escMsg)

	if updatedModel == nil {
		t.Error("Expected Update to return a model")
	}

	_ = cmd
}

func TestModel_HandleCommand(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	// Set some input text
	model.textInput.SetValue("help")

	// Test handling command
	updatedModel, cmd := model.handleCommand()

	if updatedModel == nil {
		t.Error("Expected handleCommand to return a model")
	}

	_ = cmd
}

func TestModel_View(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	// Set model as ready
	model.ready = true
	model.width = 100
	model.height = 50

	// Test view rendering
	view := model.View()

	if view == "" {
		t.Error("Expected View to return non-empty string")
	}
}

func TestModel_Update_TickMsg(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	// Test tick message
	msg := model.tickStatus()
	updatedModel, cmd := model.Update(msg)

	if updatedModel == nil {
		t.Fatal("Expected Update to return a model")
	}

	_ = cmd
}

func TestModel_Update_TextInputMsg(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	// Test text input message
	msg := textinput.New()
	updatedModel, cmd := model.Update(msg)

	if updatedModel == nil {
		t.Fatal("Expected Update to return a model")
	}

	_ = cmd
}

func TestModel_Update_ViewportMsg(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	// Test viewport message
	msg := viewport.New(80, 24)
	updatedModel, cmd := model.Update(msg)

	if updatedModel == nil {
		t.Fatal("Expected Update to return a model")
	}

	_ = cmd
}

func TestModel_Update_KeyMsgEnter(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	// Set some input text
	model.textInput.SetValue("help")

	// Test Enter key
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, cmd := model.Update(msg)

	if updatedModel == nil {
		t.Fatal("Expected Update to return a model")
	}

	_ = cmd
}

func TestModel_Update_KeyMsgUp(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	// Add some history
	model.commandHistory = []string{"help", "clear"}

	// Test Up key
	msg := tea.KeyMsg{Type: tea.KeyUp}
	updatedModel, cmd := model.Update(msg)

	if updatedModel == nil {
		t.Fatal("Expected Update to return a model")
	}

	_ = cmd
}

func TestModel_Update_KeyMsgDown(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	// Add some history and navigate up first
	model.commandHistory = []string{"help", "clear"}
	model.historyIndex = 0

	// Test Down key
	msg := tea.KeyMsg{Type: tea.KeyDown}
	updatedModel, cmd := model.Update(msg)

	if updatedModel == nil {
		t.Fatal("Expected Update to return a model")
	}

	_ = cmd
}

func TestModel_Update_KeyMsgCtrlC(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	// Test Ctrl+C key
	msg := tea.KeyMsg{Type: tea.KeyCtrlC}
	updatedModel, cmd := model.Update(msg)

	if updatedModel == nil {
		t.Fatal("Expected Update to return a model")
	}

	// Check that quitting flag is set
	modelPtr := updatedModel.(*Model)
	if !modelPtr.quitting {
		t.Error("Expected quitting to be true after Ctrl+C")
	}

	_ = cmd
}

func TestModel_Update_KeyMsgEsc(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	// Set some input text first
	model.textInput.SetValue("test input")

	// Test Esc key
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, cmd := model.Update(msg)

	if updatedModel == nil {
		t.Fatal("Expected Update to return a model")
	}

	// Check that input is cleared
	modelPtr := updatedModel.(*Model)
	if modelPtr.textInput.Value() != "" {
		t.Error("Expected text input to be cleared after Esc")
	}

	_ = cmd
}

func TestModel_Update_KeyMsgUpDown_ViewportMode(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	// Unfocus text input to test viewport scrolling
	model.textInput.Blur()

	// Test Up key in viewport mode
	msgUp := tea.KeyMsg{Type: tea.KeyUp}
	updatedModel, cmd := model.Update(msgUp)

	if updatedModel == nil {
		t.Fatal("Expected Update to return a model")
	}

	_ = cmd

	// Test Down key in viewport mode
	msgDown := tea.KeyMsg{Type: tea.KeyDown}
	updatedModel2, cmd2 := updatedModel.Update(msgDown)

	if updatedModel2 == nil {
		t.Fatal("Expected Update to return a model")
	}

	_ = cmd2
}

func TestModel_Update_KeyMsgOther(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	// Test other key (should update text input)
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
	updatedModel, cmd := model.Update(msg)

	if updatedModel == nil {
		t.Fatal("Expected Update to return a model")
	}

	_ = cmd
}

func TestModel_Update_StatusTickMsg(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	// Set a status first
	model.setStatus("Test status", 2*time.Second)

	// Test status tick message
	msg := statusTickMsg(time.Now())
	updatedModel, cmd := model.Update(msg)

	if updatedModel == nil {
		t.Fatal("Expected Update to return a model")
	}

	_ = cmd
}

func TestModel_Update_UnknownMsg(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	// Test unknown message type
	msg := "unknown message"
	updatedModel, cmd := model.Update(msg)

	if updatedModel == nil {
		t.Fatal("Expected Update to return a model")
	}

	_ = cmd
}

func TestModel_SubscribeToWatcher_WithSubscriber(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	// Test subscribe to watcher
	model.subscribeToWatcher()

	// Verify subscriber was added (indirectly by checking no panic)
	if model.watcher == nil {
		t.Error("Expected watcher to be set")
	}
}

func TestModel_SubscribeToWatcher_Comprehensive(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	// Set up the model properly
	model.ready = true

	// Test subscribe to watcher - this should set up the callback
	model.subscribeToWatcher()

	// Verify the watcher adapter exists
	if model.watcherAdapter == nil {
		t.Error("Expected watcherAdapter to be set")
	}

	// Test that the subscription doesn't panic with different model states
	model.subscribeToWatcher() // Should be safe to call multiple times
}

func TestModel_AddLogLineFromWatcher_Coverage(t *testing.T) {
	watcher := watcher.NewWatcherService()
	model := NewModel(watcher)

	// Test addLogLineFromWatcher with different log lines
	logLine1 := events.LogLine{
		Text:  "Test message 1",
		Level: "info",
		Time:  time.Now(),
	}

	logLine2 := events.LogLine{
		Text:  "Test message 2",
		Level: "error",
		Time:  time.Now(),
	}

	// These should not panic
	model.addLogLineFromWatcher(logLine1)
	model.addLogLineFromWatcher(logLine2)

	// Verify logs were added (basic check)
	if len(model.logLines) == 0 {
		t.Error("Expected logs to be added")
	}
}
