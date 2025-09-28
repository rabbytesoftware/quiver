package handlers

import (
	"testing"

	"github.com/rabbytesoftware/quiver/cmd/quiver/ui/domain/commands"
	uievents "github.com/rabbytesoftware/quiver/cmd/quiver/ui/domain/events"
	"github.com/rabbytesoftware/quiver/cmd/quiver/ui/queries"
	"github.com/rabbytesoftware/quiver/cmd/quiver/ui/queries/models"
	"github.com/rabbytesoftware/quiver/cmd/quiver/ui/services"
	"github.com/rabbytesoftware/quiver/internal/core/watcher"
)

func TestNewHandler(t *testing.T) {
	watcherAdapter := services.NewWatcherAdapter(nil)
	queryService := queries.NewService("http://example.com")

	handler := NewHandler(watcherAdapter, queryService)

	if handler == nil {
		t.Fatal("Expected handler to be created")
	}

	if handler.watcherAdapter != watcherAdapter {
		t.Error("Expected watcher adapter to be set")
	}

	if handler.queryService != queryService {
		t.Error("Expected query service to be set")
	}
}

func TestHandler_Handle_Help(t *testing.T) {
	watcherAdapter := services.NewWatcherAdapter(nil)
	queryService := queries.NewService("http://example.com")
	handler := NewHandler(watcherAdapter, queryService)

	cmd := commands.Command{
		Kind: commands.CmdHelp,
	}

	events := handler.Handle(cmd)

	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	// Check if it's a HelpRequested event
	switch event := events[0].(type) {
	case uievents.HelpRequested:
		if event.HelpText == "" {
			t.Error("Expected help text to be non-empty")
		}
	default:
		t.Errorf("Expected HelpRequested event, got %T", event)
	}
}

func TestHandler_Handle_Clear(t *testing.T) {
	watcherAdapter := services.NewWatcherAdapter(nil)
	handler := NewHandler(watcherAdapter, nil)

	cmd := commands.Command{
		Kind: commands.CmdClear,
	}

	events := handler.Handle(cmd)

	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	// Check if it's a Cleared event
	switch events[0].(type) {
	case uievents.Cleared:
		// Success
	default:
		t.Errorf("Expected Cleared event, got %T", events[0])
	}
}

func TestHandler_Handle_UnknownCommand(t *testing.T) {
	watcherAdapter := services.NewWatcherAdapter(nil)
	handler := NewHandler(watcherAdapter, nil)

	cmd := commands.Command{
		Kind: commands.CommandKind("invalid"), // Invalid command kind
	}

	events := handler.Handle(cmd)

	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	// Check if it's a CommandError event
	switch event := events[0].(type) {
	case uievents.CommandError:
		if event.Message == "" {
			t.Error("Expected non-empty error message")
		}
	default:
		t.Errorf("Expected CommandError event, got %T", event)
	}
}

func TestHandler_Handle_Query_NoService(t *testing.T) {
	watcherAdapter := services.NewWatcherAdapter(nil)
	handler := NewHandler(watcherAdapter, nil)

	cmd := commands.Command{
		Kind:          commands.CmdQuery,
		OriginalInput: "test query",
	}

	events := handler.Handle(cmd)

	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	// Check if it's a CommandError event
	switch event := events[0].(type) {
	case uievents.CommandError:
		if event.Message != "queries system not loaded or available" {
			t.Errorf("Unexpected error message: %q", event.Message)
		}
	default:
		t.Errorf("Expected CommandError event, got %T", event)
	}
}

func TestHandler_Handle_Filter_Clear(t *testing.T) {
	watcherService := watcher.NewWatcherService()
	watcherAdapter := services.NewWatcherAdapter(watcherService)
	handler := NewHandler(watcherAdapter, nil)

	cmd := commands.Command{
		Kind: commands.CmdFilter,
		Args: []string{},
	}

	events := handler.Handle(cmd)

	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	switch event := events[0].(type) {
	case uievents.FilterApplied:
		if event.Pattern != "" {
			t.Errorf("Expected empty pattern, got %q", event.Pattern)
		}
	default:
		t.Errorf("Expected FilterApplied event, got %T", event)
	}
}

func TestHandler_Handle_Filter_Set(t *testing.T) {
	watcherService := watcher.NewWatcherService()
	watcherAdapter := services.NewWatcherAdapter(watcherService)
	handler := NewHandler(watcherAdapter, nil)

	cmd := commands.Command{
		Kind: commands.CmdFilter,
		Args: []string{"error", "warning"},
	}

	events := handler.Handle(cmd)

	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	switch event := events[0].(type) {
	case uievents.FilterApplied:
		expectedPattern := "error warning"
		if event.Pattern != expectedPattern {
			t.Errorf("Expected pattern %q, got %q", expectedPattern, event.Pattern)
		}
	default:
		t.Errorf("Expected FilterApplied event, got %T", event)
	}
}

func TestHandler_Handle_Level_Valid(t *testing.T) {
	watcherService := watcher.NewWatcherService()
	watcherAdapter := services.NewWatcherAdapter(watcherService)
	handler := NewHandler(watcherAdapter, nil)

	cmd := commands.Command{
		Kind: commands.CmdLevel,
		Args: []string{"debug"},
	}

	events := handler.Handle(cmd)

	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	switch event := events[0].(type) {
	case uievents.LevelChanged:
		if event.Level != "debug" {
			t.Errorf("Expected level 'debug', got %q", event.Level)
		}
	default:
		t.Errorf("Expected LevelChanged event, got %T", event)
	}
}

func TestHandler_Handle_Level_NoArgs(t *testing.T) {
	watcherService := watcher.NewWatcherService()
	watcherAdapter := services.NewWatcherAdapter(watcherService)
	handler := NewHandler(watcherAdapter, nil)

	cmd := commands.Command{
		Kind: commands.CmdLevel,
		Args: []string{},
	}

	events := handler.Handle(cmd)

	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	switch event := events[0].(type) {
	case uievents.CommandError:
		if event.Message != "level command requires exactly one argument" {
			t.Errorf("Unexpected error message: %q", event.Message)
		}
	default:
		t.Errorf("Expected CommandError event, got %T", event)
	}
}

func TestHandler_Handle_Level_TooManyArgs(t *testing.T) {
	watcherService := watcher.NewWatcherService()
	watcherAdapter := services.NewWatcherAdapter(watcherService)
	handler := NewHandler(watcherAdapter, nil)

	cmd := commands.Command{
		Kind: commands.CmdLevel,
		Args: []string{"debug", "info"},
	}

	events := handler.Handle(cmd)

	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	switch event := events[0].(type) {
	case uievents.CommandError:
		if event.Message != "level command requires exactly one argument" {
			t.Errorf("Unexpected error message: %q", event.Message)
		}
	default:
		t.Errorf("Expected CommandError event, got %T", event)
	}
}

func TestHandler_Handle_Pause_NotPaused(t *testing.T) {
	watcherService := watcher.NewWatcherService()
	watcherAdapter := services.NewWatcherAdapter(watcherService)
	handler := NewHandler(watcherAdapter, nil)

	cmd := commands.Command{
		Kind: commands.CmdPause,
	}

	events := handler.Handle(cmd)

	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	switch events[0].(type) {
	case uievents.StreamPaused:
		// Success
	default:
		t.Errorf("Expected StreamPaused event, got %T", events[0])
	}
}

func TestHandler_Handle_Pause_AlreadyPaused(t *testing.T) {
	watcherService := watcher.NewWatcherService()
	watcherAdapter := services.NewWatcherAdapter(watcherService)
	watcherAdapter.Pause() // Pause first
	handler := NewHandler(watcherAdapter, nil)

	cmd := commands.Command{
		Kind: commands.CmdPause,
	}

	events := handler.Handle(cmd)

	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	switch event := events[0].(type) {
	case uievents.CommandError:
		if event.Message != "log streaming is already paused" {
			t.Errorf("Unexpected error message: %q", event.Message)
		}
	default:
		t.Errorf("Expected CommandError event, got %T", event)
	}
}

func TestHandler_Handle_Resume_Paused(t *testing.T) {
	watcherService := watcher.NewWatcherService()
	watcherAdapter := services.NewWatcherAdapter(watcherService)
	watcherAdapter.Pause() // Pause first
	handler := NewHandler(watcherAdapter, nil)

	cmd := commands.Command{
		Kind: commands.CmdResume,
	}

	events := handler.Handle(cmd)

	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	switch events[0].(type) {
	case uievents.StreamResumed:
		// Success
	default:
		t.Errorf("Expected StreamResumed event, got %T", events[0])
	}
}

func TestHandler_Handle_Resume_NotPaused(t *testing.T) {
	watcherService := watcher.NewWatcherService()
	watcherAdapter := services.NewWatcherAdapter(watcherService)
	handler := NewHandler(watcherAdapter, nil)

	cmd := commands.Command{
		Kind: commands.CmdResume,
	}

	events := handler.Handle(cmd)

	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	switch event := events[0].(type) {
	case uievents.CommandError:
		if event.Message != "log streaming is not paused" {
			t.Errorf("Unexpected error message: %q", event.Message)
		}
	default:
		t.Errorf("Expected CommandError event, got %T", event)
	}
}

func TestHandler_Handle_Query_WithRealService(t *testing.T) {
	watcherAdapter := services.NewWatcherAdapter(nil)
	queryService := queries.NewService("http://invalid-url-for-testing.example")

	// Add a simple query for testing
	queryService.Queries = []models.Query{
		{
			Syntax:      "test",
			Description: "Test query",
			REST: &models.REST{
				Method: "GET",
				URL:    "/test",
			},
		},
	}

	handler := NewHandler(watcherAdapter, queryService)

	cmd := commands.Command{
		Kind:          commands.CmdQuery,
		OriginalInput: "test",
	}

	events := handler.Handle(cmd)

	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	// Since we're using an invalid URL, this should result in an error
	// but we're testing the handler logic, not the actual HTTP call
	switch event := events[0].(type) {
	case uievents.QueryExecuted:
		if event.UserInput != "test" {
			t.Errorf("Expected user input 'test', got %q", event.UserInput)
		}
	case uievents.QueryError:
		if event.UserInput != "test" {
			t.Errorf("Expected user input 'test', got %q", event.UserInput)
		}
	case uievents.CommandError:
		// This is also acceptable as the HTTP call will fail
		if event.Message == "" {
			t.Error("Expected non-empty error message")
		}
	default:
		t.Errorf("Expected QueryExecuted, QueryError, or CommandError event, got %T", event)
	}
}
