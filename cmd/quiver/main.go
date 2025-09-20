package main

import (
	"time"

	"github.com/rabbytesoftware/quiver/cmd/quiver/ui"

	"github.com/rabbytesoftware/quiver/internal/core/watcher"
)

func main() {
	watcher := watcher.NewWatcherService()
	
	go func() {
		time.Sleep(500 * time.Millisecond)
		
		for {
			time.Sleep(5 * time.Second)
			watcher.Info("Quiver TUI application started")
			watcher.Debug("Initializing watcher service...")
			watcher.Info("Watcher service ready")
			watcher.Warn("This is a sample warning message")
			watcher.Error("This is a sample error message")
			watcher.Info("Try commands like /help, /filter error, /level warn, /pause, /resume, /clear")
		}
	}()
	
	err := ui.RunUI(watcher)
	if err != nil {
		watcher.Unforeseen(err.Error())
	}
}
