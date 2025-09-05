package models

type CmdRouter interface {
	Execute(command string, args []string) (string, error)
	GetAvailableCommands() []string
	GetHelp(command string) string
}
