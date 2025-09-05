package commands

import (
	"fmt"
	"runtime"
)

type VersionCommand struct{}

func NewVersionCommand() *VersionCommand {
	return &VersionCommand{}
}

func (vc *VersionCommand) Execute(args []string) (string, error) {
	version := fmt.Sprintf("Quiver Core\nGo Version: %s\nOS: %s\nArchitecture: %s",
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH)
	
	return version, nil
}

func (vc *VersionCommand) GetDescription() string {
	return "Show version information"
}

func (vc *VersionCommand) GetUsage() string {
	return "version"
}
