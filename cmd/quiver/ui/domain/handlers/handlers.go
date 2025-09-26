package handlers

import (
	"context"
	"strings"

	"github.com/rabbytesoftware/quiver/cmd/quiver/ui/domain/commands"
	"github.com/rabbytesoftware/quiver/cmd/quiver/ui/domain/events"
	"github.com/rabbytesoftware/quiver/cmd/quiver/ui/queries"
	"github.com/rabbytesoftware/quiver/cmd/quiver/ui/services"
)

// Handler handles commands and returns events
type Handler struct {
	watcherAdapter *services.WatcherAdapter
	queryService   *queries.QueryService
}

// NewHandler creates a new command handler
func NewHandler(watcherAdapter *services.WatcherAdapter, queryService *queries.QueryService) *Handler {
	return &Handler{
		watcherAdapter: watcherAdapter,
		queryService:   queryService,
	}
}

// Handle processes a command and returns resulting events
func (h *Handler) Handle(cmd commands.Command) []events.Event {
	switch cmd.Kind {
	case commands.CmdHelp:
		return h.handleHelp()
	case commands.CmdFilter:
		return h.handleFilter(cmd.Args)
	case commands.CmdLevel:
		return h.handleLevel(cmd.Args)
	case commands.CmdPause:
		return h.handlePause()
	case commands.CmdResume:
		return h.handleResume()
	case commands.CmdClear:
		return h.handleClear()
	case commands.CmdQuery:
		return h.handleQuery(cmd.OriginalInput)
	default:
		return []events.Event{
			events.CommandError{
				Message: "unknown command: " + cmd.Kind.String(),
			},
		}
	}
}

func (h *Handler) handleHelp() []events.Event {
	baseHelp := commands.GetHelpText()
	
	if h.queryService != nil {
		queryHelp := h.queryService.GetHelpText()
		helpText := baseHelp + "\n\n" + queryHelp
		return []events.Event{
			events.HelpRequested{
				HelpText: helpText,
			},
		}
	}
	
	return []events.Event{
		events.HelpRequested{
			HelpText: baseHelp,
		},
	}
}

func (h *Handler) handleFilter(args []string) []events.Event {
	if len(args) == 0 {
		if err := h.watcherAdapter.SetFilter(""); err != nil {
			return []events.Event{
				events.CommandError{
					Message: "failed to clear filter: " + err.Error(),
				},
			}
		}
		return []events.Event{
			events.FilterApplied{
				Pattern: "",
			},
		}
	}

	pattern := strings.Join(args, " ")
	if err := h.watcherAdapter.SetFilter(pattern); err != nil {
		return []events.Event{
			events.CommandError{
				Message: "failed to set filter: " + err.Error(),
			},
		}
	}

	return []events.Event{
		events.FilterApplied{
			Pattern: pattern,
		},
	}
}

func (h *Handler) handleLevel(args []string) []events.Event {
	if len(args) != 1 {
		return []events.Event{
			events.CommandError{
				Message: "level command requires exactly one argument",
			},
		}
	}

	level := strings.ToLower(args[0])
	if err := h.watcherAdapter.SetLevel(level); err != nil {
		return []events.Event{
			events.CommandError{
				Message: "failed to set level: " + err.Error(),
			},
		}
	}

	return []events.Event{
		events.LevelChanged{
			Level: level,
		},
	}
}

// handlePause handles the pause command
func (h *Handler) handlePause() []events.Event {
	if h.watcherAdapter.IsPaused() {
		return []events.Event{
			events.CommandError{
				Message: "log streaming is already paused",
			},
		}
	}

	h.watcherAdapter.Pause()
	return []events.Event{
		events.StreamPaused{},
	}
}

// handleResume handles the resume command
func (h *Handler) handleResume() []events.Event {
	if !h.watcherAdapter.IsPaused() {
		return []events.Event{
			events.CommandError{
				Message: "log streaming is not paused",
			},
		}
	}

	h.watcherAdapter.Resume()
	return []events.Event{
		events.StreamResumed{},
	}
}

func (h *Handler) handleClear() []events.Event {
	return []events.Event{
		events.Cleared{},
	}
}

func (h *Handler) handleQuery(input string) []events.Event {
	if h.queryService == nil {
		return []events.Event{
			events.CommandError{
				Message: "queries system not loaded or available",
			},
		}
	}

	ctx := context.Background()
	result, httpStatus, responseBody, err := h.queryService.HandleCommand(ctx, input)
	if err != nil {
		// If we have HTTP status and response body, create a QueryError event
		if httpStatus > 0 {
			return []events.Event{
				events.QueryError{
					UserInput:    input,
					HTTPStatus:   httpStatus,
					ResponseBody: responseBody,
				},
			}
		}
		// Otherwise, fall back to regular CommandError
		return []events.Event{
			events.CommandError{
				Message: "query execution error: " + err.Error(),
			},
		}
	}

	return []events.Event{
		events.QueryExecuted{
			DisplayText: result,
		},
	}
}
