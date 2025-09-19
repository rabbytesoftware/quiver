package handlers

import (
	"strings"

	"github.com/rabbytesoftware/quiver/internal/infrastructure/ui/domain/commands"
	"github.com/rabbytesoftware/quiver/internal/infrastructure/ui/domain/events"
	"github.com/rabbytesoftware/quiver/internal/infrastructure/ui/services/logstream"
)

// Handler handles commands and returns events
type Handler struct {
	logService logstream.LogService
}

// NewHandler creates a new command handler
func NewHandler(logService logstream.LogService) *Handler {
	return &Handler{
		logService: logService,
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
	default:
		return []events.Event{
			events.CommandError{
				Message: "unknown command: " + cmd.Kind.String(),
			},
		}
	}
}

// handleHelp handles the help command
func (h *Handler) handleHelp() []events.Event {
	return []events.Event{
		events.HelpRequested{
			HelpText: commands.GetHelpText(),
		},
	}
}

// handleFilter handles the filter command
func (h *Handler) handleFilter(args []string) []events.Event {
	if len(args) == 0 {
		// Clear filter
		if err := h.logService.SetFilter(""); err != nil {
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

	// Apply filter
	pattern := strings.Join(args, " ")
	if err := h.logService.SetFilter(pattern); err != nil {
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

// handleLevel handles the level command
func (h *Handler) handleLevel(args []string) []events.Event {
	if len(args) != 1 {
		return []events.Event{
			events.CommandError{
				Message: "level command requires exactly one argument",
			},
		}
	}

	level := strings.ToLower(args[0])
	if err := h.logService.SetLevel(level); err != nil {
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
	if h.logService.IsPaused() {
		return []events.Event{
			events.CommandError{
				Message: "log streaming is already paused",
			},
		}
	}

	h.logService.Pause()
	return []events.Event{
		events.StreamPaused{},
	}
}

// handleResume handles the resume command
func (h *Handler) handleResume() []events.Event {
	if !h.logService.IsPaused() {
		return []events.Event{
			events.CommandError{
				Message: "log streaming is not paused",
			},
		}
	}

	h.logService.Resume()
	return []events.Event{
		events.StreamResumed{},
	}
}

// handleClear handles the clear command
func (h *Handler) handleClear() []events.Event {
	return []events.Event{
		events.Cleared{},
	}
}
