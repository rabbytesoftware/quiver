package logger

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/rabbytesoftware/quiver/internal/config"
)

// Logger represents the application logger
type Logger struct {
	*slog.Logger
	config config.LoggerConfig
}

// Level represents log levels
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// String returns the string representation of the log level
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// New creates a new logger instance
func New(cfg config.LoggerConfig) *Logger {
	// Ensure log directory exists
	if err := os.MkdirAll(cfg.LogDir, 0755); err != nil {
		fmt.Printf("Failed to create log directory: %v\n", err)
	}

	// Create log file
	logFile := filepath.Join(cfg.LogDir, fmt.Sprintf("quiver-%s.log", time.Now().Format("2006-01-02")))
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		file = os.Stdout
	}

	// Parse log level
	var level slog.Level
	switch cfg.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// Create handler options
	opts := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Customize timestamp format
			if a.Key == slog.TimeKey {
				return slog.Attr{
					Key:   slog.TimeKey,
					Value: slog.StringValue(a.Value.Time().Format("2006-01-02 15:04:05")),
				}
			}
			return a
		},
	}

	// Create structured logger
	handler := slog.NewJSONHandler(file, opts)
	logger := slog.New(handler)

	return &Logger{
		Logger: logger,
		config: cfg,
	}
}

// WithService returns a logger with service context
func (l *Logger) WithService(service string) *Logger {
	return &Logger{
		Logger: l.Logger.With("service", service),
		config: l.config,
	}
}

// WithContext returns a logger with additional context
func (l *Logger) WithContext(ctx map[string]interface{}) *Logger {
	args := make([]interface{}, 0, len(ctx)*2)
	for k, v := range ctx {
		args = append(args, k, v)
	}
	return &Logger{
		Logger: l.Logger.With(args...),
		config: l.config,
	}
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.Logger.Debug(fmt.Sprintf(msg, args...))
}

// Info logs an info message
func (l *Logger) Info(msg string, args ...interface{}) {
	l.Logger.Info(fmt.Sprintf(msg, args...))
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.Logger.Warn(fmt.Sprintf(msg, args...))
}

// Error logs an error message
func (l *Logger) Error(msg string, args ...interface{}) {
	l.Logger.Error(fmt.Sprintf(msg, args...))
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.Logger.Error(fmt.Sprintf(msg, args...))
	os.Exit(1)
}

// Load logs a loading message (alias for Info with specific formatting)
func (l *Logger) Load(msg string, args ...interface{}) {
	l.Logger.Info(fmt.Sprintf("🔄 "+msg, args...))
}

// Ok logs a success message (alias for Info with specific formatting)
func (l *Logger) Ok(msg string, args ...interface{}) {
	l.Logger.Info(fmt.Sprintf("✅ "+msg, args...))
} 