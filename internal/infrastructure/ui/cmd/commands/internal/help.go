package commands

import (
	"fmt"
	"strings"

	"github.com/rabbytesoftware/quiver/internal/infrastructure/ui/cmd/models"
)

type HelpCommand struct {
	processor models.CmdRouter
}

func NewHelpCommand(processor models.CmdRouter) *HelpCommand {
	return &HelpCommand{processor: processor}
}

func (hc *HelpCommand) Execute(args []string) (string, error) {
	if len(args) == 0 {
		var help strings.Builder
		help.WriteString("Available commands:\n")
		commands := hc.processor.GetAvailableCommands()

		for _, name := range commands {
			cmd := hc.processor.GetHelp(name)
			help.WriteString(fmt.Sprintf("  %-10s - %s\n", name, cmd))
		}
		
		help.WriteString("\nUse 'help <command>' for detailed information about a specific command.")
		return help.String(), nil
	}
	
	command := args[0]
	return hc.processor.GetHelp(command), nil
}

func (hc *HelpCommand) GetDescription() string {
	return "Show help information"
}

func (hc *HelpCommand) GetUsage() string {
	return "help [command]"
}
