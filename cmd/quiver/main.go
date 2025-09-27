package main

import (
	"fmt"
	"time"

	"github.com/rabbytesoftware/quiver/cmd/quiver/ui"
	"github.com/rabbytesoftware/quiver/internal"

	"github.com/rabbytesoftware/quiver/internal/core/metadata"
)

func main() {
	internal := internal.NewInternal()
	watcher := internal.GetCore().GetWatcher()

	go internal.Run()

	go func() {
		time.Sleep(5 * time.Second)

		watcher.Info(fmt.Sprintf(
			"%s %s - Initializing...",
			metadata.GetName(),
			metadata.GetVersion(),
		))
	}()

	err := ui.RunUI(watcher)
	if err != nil {
		watcher.Unforeseen(err.Error())
	}
}
