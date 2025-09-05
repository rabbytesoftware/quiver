package interpreter

import (
	"fmt"
	"strings"

	"github.com/rabbytesoftware/quiver/internal/infrastructure/ui/cmd/commands"
	"github.com/rabbytesoftware/quiver/internal/infrastructure/ui/cmd/models"
)

type CommandInterpreter struct {
	processor models.CmdRouter
}

func NewCommandInterpreter() *CommandInterpreter {
	return &CommandInterpreter{
		processor: commands.NewCommandProcessor(),
	}
}

func (ci *CommandInterpreter) ProcessCommand(input string) (string, error) {
	input = strings.TrimSpace(input)
	
	if input == "" {
		return "", nil
	}

	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "", fmt.Errorf("empty command")
	}

	command := parts[0]
	args := parts[1:]

	result, err := ci.processor.Execute(command, args)
	if err != nil {
		return "", fmt.Errorf("command execution failed: %w", err)
	}

	return result, nil
}

func (ci *CommandInterpreter) GetAvailableCommands() []string {
	return ci.processor.GetAvailableCommands()
}

func (ci *CommandInterpreter) GetHelp(command string) string {
	return ci.processor.GetHelp(command)
}
