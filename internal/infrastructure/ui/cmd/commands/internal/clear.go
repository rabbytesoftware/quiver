package commands

type ClearCommand struct{}

func NewClearCommand() *ClearCommand {
	return &ClearCommand{}
}

func (cc *ClearCommand) Execute(args []string) (string, error) {
	return "CLEAR_SCREEN", nil
}

func (cc *ClearCommand) GetDescription() string {
	return "Clear the screen/logs"
}

func (cc *ClearCommand) GetUsage() string {
	return "clear"
}
