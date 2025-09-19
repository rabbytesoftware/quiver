package main

import (
	"log"
	"time"

	"github.com/rabbytesoftware/quiver/internal/infrastructure/ui"
	"github.com/rabbytesoftware/quiver/internal/infrastructure/watcher"
)

func main() {
	watcher := watcher.NewWatcherService()
	
	// Start the UI in a goroutine so we can send some initial logs
	go func() {
		// Give the UI a moment to initialize
		time.Sleep(500 * time.Millisecond)
		
		// Send some sample log messages to demonstrate the integration
		watcher.Info("Quiver TUI application started")
		watcher.Debug("Initializing watcher service...")
		watcher.Info("Watcher service ready")
		watcher.Warn("This is a sample warning message")
		watcher.Error("This is a sample error message")
		watcher.Info("Try commands like /help, /filter error, /level warn, /pause, /resume, /clear")
	}()
	
	err := ui.RunUI(watcher)
	if err != nil {
		log.Fatal(err)
	}
}
