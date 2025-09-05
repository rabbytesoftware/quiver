package stdout

import (
	"sync"
	"time"

	"github.com/rabbytesoftware/quiver/internal/infrastructure/ui/stdout/models"
)

// Logger defines the interface for logging functionality
type Logger interface {
	Debug(message string)
	Info(message string)
	Warning(message string)
	Error(message string)
	Unforeseen(message string)
	SetLevel(level models.LogLevel)
	GetLevel() models.LogLevel
	AddEntry(entry models.LogEntry)
	GetEntries() []models.LogEntry
	Clear()
}

// UILogger implements the Logger interface for UI logging
type UILogger struct {
	entries    []models.LogEntry
	level      models.LogLevel
	mu         sync.RWMutex
	subscribers []func(models.LogEntry)
}

// NewUILogger creates a new UI logger instance
func NewUILogger() *UILogger {
	return &UILogger{
		entries:    make([]models.LogEntry, 0),
		level:      models.LogLevelInfo, // Default to INFO level
		subscribers: make([]func(models.LogEntry), 0),
	}
}

// Debug logs a debug message
func (l *UILogger) Debug(message string) {
	l.log(models.LogLevelDebug, message)
}

// Info logs an info message
func (l *UILogger) Info(message string) {
	l.log(models.LogLevelInfo, message)
}

// Warning logs a warning message
func (l *UILogger) Warning(message string) {
	l.log(models.LogLevelWarning, message)
}

// Error logs an error message
func (l *UILogger) Error(message string) {
	l.log(models.LogLevelError, message)
}

// Unforeseen logs an unforeseen message
func (l *UILogger) Unforeseen(message string) {
	l.log(models.LogLevelUnforeseen, message)
}

// SetLevel sets the minimum log level to display
func (l *UILogger) SetLevel(level models.LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// GetLevel returns the current log level
func (l *UILogger) GetLevel() models.LogLevel {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.level
}

// AddEntry adds a log entry directly
func (l *UILogger) AddEntry(entry models.LogEntry) {
	l.mu.Lock()
	defer l.mu.Unlock()
	
	// Check if the entry level meets the minimum threshold
	if models.LogLevelPriority[entry.Level] >= models.LogLevelPriority[l.level] {
		l.entries = append(l.entries, entry)
		
		// Notify subscribers
		for _, subscriber := range l.subscribers {
			go subscriber(entry)
		}
	}
}

// GetEntries returns all log entries
func (l *UILogger) GetEntries() []models.LogEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	
	// Return a copy to prevent external modification
	entries := make([]models.LogEntry, len(l.entries))
	copy(entries, l.entries)
	return entries
}

// Clear clears all log entries
func (l *UILogger) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = make([]models.LogEntry, 0)
}

// Subscribe adds a callback function that will be called when new entries are added
func (l *UILogger) Subscribe(callback func(models.LogEntry)) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.subscribers = append(l.subscribers, callback)
}

// log is the internal method that handles the actual logging
func (l *UILogger) log(level models.LogLevel, message string) {
	entry := models.LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
	}
	l.AddEntry(entry)
}

// IsLevelEnabled checks if a given log level is enabled
func (l *UILogger) IsLevelEnabled(level models.LogLevel) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return models.LogLevelPriority[level] >= models.LogLevelPriority[l.level]
}
