package main

import (
	"fmt"

	"github.com/rabbytesoftware/quiver/cmd/quiver/assets"
	"github.com/rabbytesoftware/quiver/cmd/quiver/ui"
	"github.com/rabbytesoftware/quiver/internal"
	"github.com/rabbytesoftware/quiver/internal/core/metadata"
)

func main() {
	iconManager := assets.NewIconManager()
	defer iconManager.Cleanup()

	internal := internal.NewInternal()
	watcher := internal.GetCore().GetWatcher()

	go internal.Run()

	go watcher.Info(fmt.Sprintf(
		"%s %s '%s' - Initializing with embedded icon support...",
		metadata.GetName(),
		metadata.GetVersion(),
		metadata.GetVersionCodename(),
	))

	err := ui.RunUI(watcher)
	if err != nil {
		watcher.Unforeseen(err.Error())
	}
}
