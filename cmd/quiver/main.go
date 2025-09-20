package main

import (
	"fmt"
	"time"

	"github.com/rabbytesoftware/quiver/cmd/quiver/ui"

	"github.com/rabbytesoftware/quiver/internal/core/metadata"
	"github.com/rabbytesoftware/quiver/internal/core/watcher"
)

func main() {
	watcher := watcher.NewWatcherService()
	
	err := ui.RunUI(watcher)
	if err != nil {
		watcher.Unforeseen(err.Error())
	}

	go func() {
		time.Sleep(5 * time.Second)
	
		watcher.Info(fmt.Sprintf(
		"%s %s - Initializing...",
			metadata.GetName(), 
			metadata.GetVersion(),
		))
	}()
}
