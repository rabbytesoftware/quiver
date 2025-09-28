package commands

import (
	"fmt"
	"strings"
)

// CommandKind represents the type of command
type CommandKind string

const (
	CmdHelp   CommandKind = "help"
	CmdFilter CommandKind = "filter"
	CmdLevel  CommandKind = "level"
	CmdPause  CommandKind = "pause"
	CmdResume CommandKind = "resume"
	CmdClear  CommandKind = "clear"
	CmdQuery  CommandKind = "query" // For query system commands
)

// Command represents a parsed command with its arguments
type Command struct {
	Kind CommandKind
	Args []string
	// OriginalInput stores the original input for query commands
	OriginalInput string
}

// CommandError represents an error that occurred during command parsing or execution
type CommandError struct {
	Message string
	Input   string
}

func (e CommandError) Error() string {
	return fmt.Sprintf("command error: %s (input: %q)", e.Message, e.Input)
}

// NewCommandError creates a new command error
func NewCommandError(message, input string) error {
	return CommandError{
		Message: message,
		Input:   input,
	}
}

// ValidLogLevels contains the accepted log levels
var ValidLogLevels = map[string]bool{
	"debug": true,
	"info":  true,
	"warn":  true,
	"error": true,
}

// String returns the string representation of a CommandKind
func (ck CommandKind) String() string {
	return string(ck)
}

// IsValid checks if a CommandKind is valid
func (ck CommandKind) IsValid() bool {
	switch ck {
	case CmdHelp, CmdFilter, CmdLevel, CmdPause, CmdResume, CmdClear, CmdQuery:
		return true
	default:
		return false
	}
}

// Validate validates a command and its arguments
func (c Command) Validate() error {
	if !c.Kind.IsValid() {
		return NewCommandError("unknown command", c.Kind.String())
	}

	switch c.Kind {
	case CmdFilter:
		if len(c.Args) == 0 {
			return NewCommandError("filter command requires a regex pattern", strings.Join(c.Args, " "))
		}
	case CmdLevel:
		if len(c.Args) != 1 {
			return NewCommandError("level command requires exactly one argument", strings.Join(c.Args, " "))
		}
		level := strings.ToLower(c.Args[0])
		if !ValidLogLevels[level] {
			return NewCommandError("invalid log level, must be one of: debug, info, warn, error", level)
		}
	case CmdHelp, CmdPause, CmdResume, CmdClear:
		// These commands don't accept arguments
		if len(c.Args) > 0 {
			return NewCommandError(fmt.Sprintf("%s command does not accept arguments", c.Kind.String()), strings.Join(c.Args, " "))
		}
	case CmdQuery:
		// Query commands are validated by the queries system
	}

	return nil
}
