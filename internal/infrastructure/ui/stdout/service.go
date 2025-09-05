package stdout

import (
	"sync"
	"time"

	"github.com/rabbytesoftware/quiver/internal/infrastructure/ui/stdout/models"
)

type LogService struct {
	logger      Logger
	subscribers []func(models.LogEntry)
	mu          sync.RWMutex
}

func NewLogService() *LogService {
	return &LogService{
		logger:      NewUILogger(),
		subscribers: make([]func(models.LogEntry), 0),
	}
}

func (ls *LogService) GetLogger() Logger {
	return ls.logger
}

func (ls *LogService) Debug(message string) {
	ls.logger.Debug(message)
	ls.notifySubscribers(models.LogEntry{
		Timestamp: time.Now(),
		Level:     models.LogLevelDebug,
		Message:   message,
	})
}

func (ls *LogService) Info(message string) {
	ls.logger.Info(message)
	ls.notifySubscribers(models.LogEntry{
		Timestamp: time.Now(),
		Level:     models.LogLevelInfo,
		Message:   message,
	})
}

func (ls *LogService) Warning(message string) {
	ls.logger.Warning(message)
	ls.notifySubscribers(models.LogEntry{
		Timestamp: time.Now(),
		Level:     models.LogLevelWarning,
		Message:   message,
	})
}

func (ls *LogService) Error(message string) {
	ls.logger.Error(message)
	ls.notifySubscribers(models.LogEntry{
		Timestamp: time.Now(),
		Level:     models.LogLevelError,
		Message:   message,
	})
}

func (ls *LogService) Unforeseen(message string) {
	ls.logger.Unforeseen(message)
	ls.notifySubscribers(models.LogEntry{
		Timestamp: time.Now(),
		Level:     models.LogLevelUnforeseen,
		Message:   message,
	})
}

func (ls *LogService) SetLevel(level models.LogLevel) {
	ls.logger.SetLevel(level)
}

func (ls *LogService) GetLevel() models.LogLevel {
	return ls.logger.GetLevel()
}

func (ls *LogService) GetEntries() []models.LogEntry {
	return ls.logger.GetEntries()
}

func (ls *LogService) Clear() {
	ls.logger.Clear()
}

func (ls *LogService) Subscribe(callback func(models.LogEntry)) {
	ls.mu.Lock()
	defer ls.mu.Unlock()
	ls.subscribers = append(ls.subscribers, callback)
}

func (ls *LogService) notifySubscribers(entry models.LogEntry) {
	ls.mu.RLock()
	defer ls.mu.RUnlock()
	
	for _, subscriber := range ls.subscribers {
		go subscriber(entry)
	}
}
