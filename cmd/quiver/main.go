package main

import (
	"time"

	"github.com/rabbytesoftware/quiver/internal/infrastructure/ui/tui"
	"github.com/rabbytesoftware/quiver/internal/infrastructure/watcher"
)

func main() {
	watcher := watcher.NewWatcherService()

	err := tui.RunTUI(watcher)
	if err != nil {
		watcher.Unforeseen(err.Error())
	}

	time.Sleep(5 * time.Second)

	watcher.Debug("Freeman!")
	watcher.Info("Freeman!")
	watcher.Warning("Freeman!")
	watcher.Error("Freeman!")
	watcher.Unforeseen("Freeman!")
}
