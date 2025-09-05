package commands

import (
	"fmt"

	internalCommands "github.com/rabbytesoftware/quiver/internal/infrastructure/ui/cmd/commands/internal"
	"github.com/rabbytesoftware/quiver/internal/infrastructure/ui/cmd/models"
)

type CommandProcessor struct {
	commands map[string]models.Command
}

func NewCommandProcessor() *CommandProcessor {
	cp := &CommandProcessor{
		commands: make(map[string]models.Command),
	}
	
	cp.RegisterCommand("help", internalCommands.NewHelpCommand(cp))
	cp.RegisterCommand("clear", internalCommands.NewClearCommand())
	cp.RegisterCommand("version", internalCommands.NewVersionCommand())
	
	return cp
}

func (cp *CommandProcessor) RegisterCommand(name string, command models.Command) {
	cp.commands[name] = command
}

func (cp *CommandProcessor) Execute(command string, args []string) (string, error) {
	cmd, exists := cp.commands[command]
	if !exists {
		return "", fmt.Errorf("unknown command: %s", command)
	}
	
	return cmd.Execute(args)
}

func (cp *CommandProcessor) GetAvailableCommands() []string {
	commands := make([]string, 0, len(cp.commands))
	for name := range cp.commands {
		commands = append(commands, name)
	}
	return commands
}

func (cp *CommandProcessor) GetHelp(command string) string {
	cmd, exists := cp.commands[command]
	if !exists {
		return fmt.Sprintf("Unknown command: %s", command)
	}
	
	return fmt.Sprintf("%s\nUsage: %s", cmd.GetDescription(), cmd.GetUsage())
}
