package commands

import (
	"regexp"
	"strings"
)

// Parse parses a command input string and returns a Command or an error
func Parse(input string) (Command, error) {
	// Trim whitespace
	input = strings.TrimSpace(input)
	
	// Ignore blank lines
	if input == "" {
		return Command{}, NewCommandError("empty command", input)
	}
	
	// Split command and arguments
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return Command{}, NewCommandError("invalid command format", input)
	}
	
	commandStr := strings.ToLower(parts[0])
	args := parts[1:]
	
	// Map string to CommandKind
	var kind CommandKind
	switch commandStr {
	case "help", "h":
		kind = CmdHelp
	case "filter", "f":
		kind = CmdFilter
	case "level", "l":
		kind = CmdLevel
	case "pause", "p":
		kind = CmdPause
	case "resume", "r":
		kind = CmdResume
	case "clear", "c":
		kind = CmdClear
	default:
		// If it's not a base TUI command, treat it as a query command
		kind = CmdQuery
	}
	
	cmd := Command{
		Kind:          kind,
		Args:          args,
		OriginalInput: input,
	}
	
	// Validate the command
	if err := cmd.Validate(); err != nil {
		return Command{}, err
	}
	
	// Special validation for filter regex
	if kind == CmdFilter && len(args) > 0 {
		pattern := strings.Join(args, " ")
		if _, err := regexp.Compile(pattern); err != nil {
			return Command{}, NewCommandError("invalid regex pattern: "+err.Error(), pattern)
		}
	}
	
	return cmd, nil
}

// GetHelpText returns help text for all available commands
func GetHelpText() string {
	return strings.TrimSpace(`
Available Commands:
  help, h              Show this help message
  filter, f <regex>    Apply regex filter to log lines
  level, l <level>     Set log level filter (debug|info|warn|error)
  pause, p             Pause log streaming
  resume, r            Resume log streaming
  clear, c             Clear the log viewport

Navigation:
  Up/Down/PgUp/PgDn      Scroll viewport
  Esc                    Clear input
  Ctrl+C                 Quit application
`)
}
