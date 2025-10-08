package main

import (
	"fmt"

	"github.com/rabbytesoftware/quiver/cmd/quiver/assets"
	"github.com/rabbytesoftware/quiver/cmd/quiver/ui"
	"github.com/rabbytesoftware/quiver/internal"
	"github.com/rabbytesoftware/quiver/internal/core/errors"
	"github.com/rabbytesoftware/quiver/internal/core/metadata"
	"github.com/rabbytesoftware/quiver/internal/core/watcher"
)

func main() {
	iconManager := assets.NewIconManager()
	defer iconManager.Cleanup()

	internal := internal.NewInternal()

	go internal.Run()

	go watcher.Info(fmt.Sprintf(
		"%s %s '%s' - Initializing with embedded icon support...",
		metadata.GetName(),
		metadata.GetVersion(),
		metadata.GetVersionCodename(),
	))

	err := ui.RunUI(watcher.GetWatcher())
	if err != nil {
		watcher.Unforeseen(errors.Throw(errors.FailedDependency, err.Error(), nil))
	}
}
