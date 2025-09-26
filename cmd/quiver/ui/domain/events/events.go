package events

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Event represents a domain-level event that can occur in the application
type Event interface {
	Type() EventType
}

// EventType represents the type of event
type EventType string

const (
	EventLogLineReceived EventType = "log_line_received"
	EventFilterApplied   EventType = "filter_applied"
	EventLevelChanged    EventType = "level_changed"
	EventStreamPaused    EventType = "stream_paused"
	EventStreamResumed   EventType = "stream_resumed"
	EventCleared         EventType = "cleared"
	EventCommandError    EventType = "command_error"
	EventHelpRequested   EventType = "help_requested"
	EventQueryExecuted   EventType = "query_executed"
)

// LogLine represents a single log line with metadata
type LogLine struct {
	Text  string
	Level string
	Time  time.Time
}

// LogLineReceived event when a new log line is received
type LogLineReceived struct {
	LogLine LogLine
}

func (e LogLineReceived) Type() EventType {
	return EventLogLineReceived
}

// FilterApplied event when a filter is applied
type FilterApplied struct {
	Pattern string
}

func (e FilterApplied) Type() EventType {
	return EventFilterApplied
}

// LevelChanged event when log level filter is changed
type LevelChanged struct {
	Level string
}

func (e LevelChanged) Type() EventType {
	return EventLevelChanged
}

// StreamPaused event when log streaming is paused
type StreamPaused struct{}

func (e StreamPaused) Type() EventType {
	return EventStreamPaused
}

// StreamResumed event when log streaming is resumed
type StreamResumed struct{}

func (e StreamResumed) Type() EventType {
	return EventStreamResumed
}

// Cleared event when the viewport is cleared
type Cleared struct{}

func (e Cleared) Type() EventType {
	return EventCleared
}

// CommandError event when a command error occurs
type CommandError struct {
	Message string
}

func (e CommandError) Type() EventType {
	return EventCommandError
}

// HelpRequested event when help is requested
type HelpRequested struct {
	HelpText string
}

func (e HelpRequested) Type() EventType {
	return EventHelpRequested
}

// QueryExecuted event when a query command is executed
type QueryExecuted struct {
	DisplayText string
}

func (e QueryExecuted) Type() EventType {
	return EventQueryExecuted
}

// Bubble Tea message types
// These wrap domain events for the Bubble Tea message system

// LogLineReceivedMsg wraps LogLineReceived for Bubble Tea
type LogLineReceivedMsg struct {
	Event LogLineReceived
}

// FilterAppliedMsg wraps FilterApplied for Bubble Tea
type FilterAppliedMsg struct {
	Event FilterApplied
}

// LevelChangedMsg wraps LevelChanged for Bubble Tea
type LevelChangedMsg struct {
	Event LevelChanged
}

// StreamPausedMsg wraps StreamPaused for Bubble Tea
type StreamPausedMsg struct {
	Event StreamPaused
}

// StreamResumedMsg wraps StreamResumed for Bubble Tea
type StreamResumedMsg struct {
	Event StreamResumed
}

// ClearedMsg wraps Cleared for Bubble Tea
type ClearedMsg struct {
	Event Cleared
}

// CommandErrorMsg wraps CommandError for Bubble Tea
type CommandErrorMsg struct {
	Event CommandError
}

// HelpRequestedMsg wraps HelpRequested for Bubble Tea
type HelpRequestedMsg struct {
	Event HelpRequested
}

// QueryExecutedMsg wraps QueryExecuted for Bubble Tea
type QueryExecutedMsg struct {
	Event QueryExecuted
}

// ToTeaMsg converts a domain event to a Bubble Tea message
func ToTeaMsg(event Event) tea.Msg {
	switch e := event.(type) {
	case LogLineReceived:
		return LogLineReceivedMsg{Event: e}
	case FilterApplied:
		return FilterAppliedMsg{Event: e}
	case LevelChanged:
		return LevelChangedMsg{Event: e}
	case StreamPaused:
		return StreamPausedMsg{Event: e}
	case StreamResumed:
		return StreamResumedMsg{Event: e}
	case Cleared:
		return ClearedMsg{Event: e}
	case CommandError:
		return CommandErrorMsg{Event: e}
	case HelpRequested:
		return HelpRequestedMsg{Event: e}
	case QueryExecuted:
		return QueryExecutedMsg{Event: e}
	default:
		return nil
	}
}
