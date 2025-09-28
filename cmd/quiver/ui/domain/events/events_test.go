package events

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestEventTypes(t *testing.T) {
	// Test all event type constants
	expectedTypes := map[EventType]string{
		EventLogLineReceived: "log_line_received",
		EventFilterApplied:   "filter_applied",
		EventLevelChanged:    "level_changed",
		EventStreamPaused:    "stream_paused",
		EventStreamResumed:   "stream_resumed",
		EventCleared:         "cleared",
		EventCommandError:    "command_error",
		EventHelpRequested:   "help_requested",
		EventQueryExecuted:   "query_executed",
		EventQueryError:      "query_error",
	}

	for eventType, expectedStr := range expectedTypes {
		if string(eventType) != expectedStr {
			t.Errorf("EventType %v should equal %q, got %q", eventType, expectedStr, string(eventType))
		}
	}
}

func TestLogLineReceived(t *testing.T) {
	logLine := LogLine{
		Text:  "test log message",
		Level: "info",
		Time:  time.Now(),
	}

	event := LogLineReceived{LogLine: logLine}

	// Test Type() method
	if event.Type() != EventLogLineReceived {
		t.Errorf("Expected type %v, got %v", EventLogLineReceived, event.Type())
	}

	// Test that it implements Event interface
	var _ Event = event

	// Test LogLine fields
	if event.LogLine.Text != "test log message" {
		t.Error("LogLine.Text not set correctly")
	}

	if event.LogLine.Level != "info" {
		t.Error("LogLine.Level not set correctly")
	}
}

func TestFilterApplied(t *testing.T) {
	event := FilterApplied{Pattern: "test.*pattern"}

	// Test Type() method
	if event.Type() != EventFilterApplied {
		t.Errorf("Expected type %v, got %v", EventFilterApplied, event.Type())
	}

	// Test that it implements Event interface
	var _ Event = event

	// Test Pattern field
	if event.Pattern != "test.*pattern" {
		t.Error("Pattern not set correctly")
	}
}

func TestLevelChanged(t *testing.T) {
	event := LevelChanged{Level: "debug"}

	// Test Type() method
	if event.Type() != EventLevelChanged {
		t.Errorf("Expected type %v, got %v", EventLevelChanged, event.Type())
	}

	// Test that it implements Event interface
	var _ Event = event

	// Test Level field
	if event.Level != "debug" {
		t.Error("Level not set correctly")
	}
}

func TestStreamPaused(t *testing.T) {
	event := StreamPaused{}

	// Test Type() method
	if event.Type() != EventStreamPaused {
		t.Errorf("Expected type %v, got %v", EventStreamPaused, event.Type())
	}

	// Test that it implements Event interface
	var _ Event = event
}

func TestStreamResumed(t *testing.T) {
	event := StreamResumed{}

	// Test Type() method
	if event.Type() != EventStreamResumed {
		t.Errorf("Expected type %v, got %v", EventStreamResumed, event.Type())
	}

	// Test that it implements Event interface
	var _ Event = event
}

func TestCleared(t *testing.T) {
	event := Cleared{}

	// Test Type() method
	if event.Type() != EventCleared {
		t.Errorf("Expected type %v, got %v", EventCleared, event.Type())
	}

	// Test that it implements Event interface
	var _ Event = event
}

func TestCommandError(t *testing.T) {
	event := CommandError{Message: "test error message"}

	// Test Type() method
	if event.Type() != EventCommandError {
		t.Errorf("Expected type %v, got %v", EventCommandError, event.Type())
	}

	// Test that it implements Event interface
	var _ Event = event

	// Test Message field
	if event.Message != "test error message" {
		t.Error("Message not set correctly")
	}
}

func TestHelpRequested(t *testing.T) {
	event := HelpRequested{HelpText: "test help text"}

	// Test Type() method
	if event.Type() != EventHelpRequested {
		t.Errorf("Expected type %v, got %v", EventHelpRequested, event.Type())
	}

	// Test that it implements Event interface
	var _ Event = event

	// Test HelpText field
	if event.HelpText != "test help text" {
		t.Error("HelpText not set correctly")
	}
}

func TestQueryExecuted(t *testing.T) {
	event := QueryExecuted{
		UserInput:    "test query",
		HTTPStatus:   200,
		ResponseBody: "success response",
	}

	// Test Type() method
	if event.Type() != EventQueryExecuted {
		t.Errorf("Expected type %v, got %v", EventQueryExecuted, event.Type())
	}

	// Test that it implements Event interface
	var _ Event = event

	// Test fields
	if event.UserInput != "test query" {
		t.Error("UserInput not set correctly")
	}

	if event.HTTPStatus != 200 {
		t.Error("HTTPStatus not set correctly")
	}

	if event.ResponseBody != "success response" {
		t.Error("ResponseBody not set correctly")
	}
}

func TestQueryError(t *testing.T) {
	event := QueryError{
		UserInput:    "test query",
		HTTPStatus:   404,
		ResponseBody: "error response",
	}

	// Test Type() method
	if event.Type() != EventQueryError {
		t.Errorf("Expected type %v, got %v", EventQueryError, event.Type())
	}

	// Test that it implements Event interface
	var _ Event = event

	// Test fields
	if event.UserInput != "test query" {
		t.Error("UserInput not set correctly")
	}

	if event.HTTPStatus != 404 {
		t.Error("HTTPStatus not set correctly")
	}

	if event.ResponseBody != "error response" {
		t.Error("ResponseBody not set correctly")
	}
}

func TestToTeaMsg(t *testing.T) {
	testCases := []struct {
		event       Event
		expectedMsg tea.Msg
		description string
	}{
		{
			LogLineReceived{LogLine: LogLine{Text: "test", Level: "info", Time: time.Now()}},
			LogLineReceivedMsg{Event: LogLineReceived{LogLine: LogLine{Text: "test", Level: "info", Time: time.Now()}}},
			"LogLineReceived",
		},
		{
			FilterApplied{Pattern: "test"},
			FilterAppliedMsg{Event: FilterApplied{Pattern: "test"}},
			"FilterApplied",
		},
		{
			LevelChanged{Level: "debug"},
			LevelChangedMsg{Event: LevelChanged{Level: "debug"}},
			"LevelChanged",
		},
		{
			StreamPaused{},
			StreamPausedMsg{Event: StreamPaused{}},
			"StreamPaused",
		},
		{
			StreamResumed{},
			StreamResumedMsg{Event: StreamResumed{}},
			"StreamResumed",
		},
		{
			Cleared{},
			ClearedMsg{Event: Cleared{}},
			"Cleared",
		},
		{
			CommandError{Message: "error"},
			CommandErrorMsg{Event: CommandError{Message: "error"}},
			"CommandError",
		},
		{
			HelpRequested{HelpText: "help"},
			HelpRequestedMsg{Event: HelpRequested{HelpText: "help"}},
			"HelpRequested",
		},
		{
			QueryExecuted{UserInput: "query", HTTPStatus: 200, ResponseBody: "ok"},
			QueryExecutedMsg{Event: QueryExecuted{UserInput: "query", HTTPStatus: 200, ResponseBody: "ok"}},
			"QueryExecuted",
		},
		{
			QueryError{UserInput: "query", HTTPStatus: 404, ResponseBody: "error"},
			QueryErrorMsg{Event: QueryError{UserInput: "query", HTTPStatus: 404, ResponseBody: "error"}},
			"QueryError",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			msg := ToTeaMsg(tc.event)

			if msg == nil {
				t.Error("ToTeaMsg returned nil")
				return
			}

			// Test that the message type matches expected
			switch expectedMsg := tc.expectedMsg.(type) {
			case LogLineReceivedMsg:
				if actualMsg, ok := msg.(LogLineReceivedMsg); ok {
					if actualMsg.Event.LogLine.Text != expectedMsg.Event.LogLine.Text {
						t.Error("LogLineReceivedMsg text mismatch")
					}
				} else {
					t.Errorf("Expected LogLineReceivedMsg, got %T", msg)
				}
			case FilterAppliedMsg:
				if actualMsg, ok := msg.(FilterAppliedMsg); ok {
					if actualMsg.Event.Pattern != expectedMsg.Event.Pattern {
						t.Error("FilterAppliedMsg pattern mismatch")
					}
				} else {
					t.Errorf("Expected FilterAppliedMsg, got %T", msg)
				}
			case LevelChangedMsg:
				if actualMsg, ok := msg.(LevelChangedMsg); ok {
					if actualMsg.Event.Level != expectedMsg.Event.Level {
						t.Error("LevelChangedMsg level mismatch")
					}
				} else {
					t.Errorf("Expected LevelChangedMsg, got %T", msg)
				}
			case StreamPausedMsg:
				if _, ok := msg.(StreamPausedMsg); !ok {
					t.Errorf("Expected StreamPausedMsg, got %T", msg)
				}
			case StreamResumedMsg:
				if _, ok := msg.(StreamResumedMsg); !ok {
					t.Errorf("Expected StreamResumedMsg, got %T", msg)
				}
			case ClearedMsg:
				if _, ok := msg.(ClearedMsg); !ok {
					t.Errorf("Expected ClearedMsg, got %T", msg)
				}
			case CommandErrorMsg:
				if actualMsg, ok := msg.(CommandErrorMsg); ok {
					if actualMsg.Event.Message != expectedMsg.Event.Message {
						t.Error("CommandErrorMsg message mismatch")
					}
				} else {
					t.Errorf("Expected CommandErrorMsg, got %T", msg)
				}
			case HelpRequestedMsg:
				if actualMsg, ok := msg.(HelpRequestedMsg); ok {
					if actualMsg.Event.HelpText != expectedMsg.Event.HelpText {
						t.Error("HelpRequestedMsg help text mismatch")
					}
				} else {
					t.Errorf("Expected HelpRequestedMsg, got %T", msg)
				}
			case QueryExecutedMsg:
				if actualMsg, ok := msg.(QueryExecutedMsg); ok {
					if actualMsg.Event.UserInput != expectedMsg.Event.UserInput {
						t.Error("QueryExecutedMsg user input mismatch")
					}
				} else {
					t.Errorf("Expected QueryExecutedMsg, got %T", msg)
				}
			case QueryErrorMsg:
				if actualMsg, ok := msg.(QueryErrorMsg); ok {
					if actualMsg.Event.UserInput != expectedMsg.Event.UserInput {
						t.Error("QueryErrorMsg user input mismatch")
					}
				} else {
					t.Errorf("Expected QueryErrorMsg, got %T", msg)
				}
			}
		})
	}
}

// UnknownEvent is a test event type that doesn't match any known types
type UnknownEvent struct{}

func (e UnknownEvent) Type() EventType { return EventType("unknown") }

func TestToTeaMsgWithUnknownEvent(t *testing.T) {
	// Create a mock event that doesn't match any known types
	unknownEvent := UnknownEvent{}
	msg := ToTeaMsg(unknownEvent)

	if msg != nil {
		t.Error("ToTeaMsg should return nil for unknown event types")
	}
}

func TestLogLine(t *testing.T) {
	now := time.Now()
	logLine := LogLine{
		Text:  "test message",
		Level: "error",
		Time:  now,
	}

	if logLine.Text != "test message" {
		t.Error("LogLine.Text not set correctly")
	}

	if logLine.Level != "error" {
		t.Error("LogLine.Level not set correctly")
	}

	if logLine.Time != now {
		t.Error("LogLine.Time not set correctly")
	}
}

func TestBubbleTeaMessageTypes(t *testing.T) {
	// Test that all Bubble Tea message types wrap their events correctly
	logEvent := LogLineReceived{LogLine: LogLine{Text: "test", Level: "info", Time: time.Now()}}
	logMsg := LogLineReceivedMsg{Event: logEvent}

	if logMsg.Event.LogLine.Text != "test" {
		t.Error("LogLineReceivedMsg does not wrap event correctly")
	}

	filterEvent := FilterApplied{Pattern: "test"}
	filterMsg := FilterAppliedMsg{Event: filterEvent}

	if filterMsg.Event.Pattern != "test" {
		t.Error("FilterAppliedMsg does not wrap event correctly")
	}

	// Test other message types
	levelMsg := LevelChangedMsg{Event: LevelChanged{Level: "debug"}}
	if levelMsg.Event.Level != "debug" {
		t.Error("LevelChangedMsg does not wrap event correctly")
	}

	pausedMsg := StreamPausedMsg{Event: StreamPaused{}}
	_ = pausedMsg // Just test it compiles

	resumedMsg := StreamResumedMsg{Event: StreamResumed{}}
	_ = resumedMsg // Just test it compiles

	clearedMsg := ClearedMsg{Event: Cleared{}}
	_ = clearedMsg // Just test it compiles

	errorMsg := CommandErrorMsg{Event: CommandError{Message: "error"}}
	if errorMsg.Event.Message != "error" {
		t.Error("CommandErrorMsg does not wrap event correctly")
	}

	helpMsg := HelpRequestedMsg{Event: HelpRequested{HelpText: "help"}}
	if helpMsg.Event.HelpText != "help" {
		t.Error("HelpRequestedMsg does not wrap event correctly")
	}

	queryMsg := QueryExecutedMsg{Event: QueryExecuted{UserInput: "query", HTTPStatus: 200, ResponseBody: "ok"}}
	if queryMsg.Event.UserInput != "query" {
		t.Error("QueryExecutedMsg does not wrap event correctly")
	}

	queryErrMsg := QueryErrorMsg{Event: QueryError{UserInput: "query", HTTPStatus: 404, ResponseBody: "error"}}
	if queryErrMsg.Event.UserInput != "query" {
		t.Error("QueryErrorMsg does not wrap event correctly")
	}
}
