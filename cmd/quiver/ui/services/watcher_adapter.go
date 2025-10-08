package services

import (
	"strings"

	"github.com/rabbytesoftware/quiver/internal/core/watcher"
	"github.com/sirupsen/logrus"
)

// WatcherAdapter adapts the watcher service to provide the interface needed by UI handlers
type WatcherAdapter struct {
	watcher  *watcher.Watcher
	filter   string
	paused   bool
	minLevel logrus.Level
}

// NewWatcherAdapter creates a new watcher adapter
func NewWatcherAdapter(w *watcher.Watcher) *WatcherAdapter {
	return &WatcherAdapter{
		watcher:  w,
		filter:   "",
		paused:   false,
		minLevel: logrus.DebugLevel,
	}
}

// SetFilter sets a filter pattern for log messages
func (wa *WatcherAdapter) SetFilter(pattern string) error {
	wa.filter = pattern
	return nil
}

// GetFilter returns the current filter pattern
func (wa *WatcherAdapter) GetFilter() string {
	return wa.filter
}

// SetLevel sets the log level from a string
func (wa *WatcherAdapter) SetLevel(level string) error {
	parsedLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	wa.minLevel = parsedLevel
	watcher.SetLevel(parsedLevel)
	return nil
}

// GetLevel returns the current log level as string
func (wa *WatcherAdapter) GetLevel() string {
	return wa.minLevel.String()
}

// IsPaused returns whether log streaming is paused
func (wa *WatcherAdapter) IsPaused() bool {
	return wa.paused
}

// Pause pauses log streaming
func (wa *WatcherAdapter) Pause() {
	wa.paused = true
}

// Resume resumes log streaming
func (wa *WatcherAdapter) Resume() {
	wa.paused = false
}

// Subscribe subscribes to log messages with filtering and pausing logic
func (wa *WatcherAdapter) Subscribe(callback func(logrus.Level, string)) {
	watcher.Subscribe(func(level logrus.Level, message string) {
		// Skip if paused
		if wa.paused {
			return
		}

		// Skip if level is below minimum
		if level > wa.minLevel {
			return
		}

		// Skip if filter is set and message doesn't match
		if wa.filter != "" && !strings.Contains(strings.ToLower(message), strings.ToLower(wa.filter)) {
			return
		}

		// Call the callback
		callback(level, message)
	})
}

// GetWatcher returns the underlying watcher service
func (wa *WatcherAdapter) GetWatcher() *watcher.Watcher {
	return wa.watcher
}
