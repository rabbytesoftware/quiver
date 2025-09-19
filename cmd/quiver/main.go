package main

import (
	"time"

	"github.com/rabbytesoftware/quiver/internal/infrastructure/watcher"
)

func main() {
	watcher := watcher.NewWatcherService()

	// err := ui.RunUI(watcher)
	// if err != nil {
	// 	watcher.Unforeseen(err.Error())
	// }

	time.Sleep(5 * time.Second)

	watcher.Debug("Freeman!")
	watcher.Info("Freeman!")
	watcher.Warn("Freeman!")
	watcher.Error("Freeman!")

	time.Sleep(60 * time.Second)
}
