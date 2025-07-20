package logger

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/pterm/pterm"
	"github.com/rabbytesoftware/quiver/internal/config"
)

// Logger represents the application logger
type Logger struct {
	fileLogger *slog.Logger
	config     config.LoggerConfig
	service    string
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

	// Create handler options for file logging (JSON)
	opts := &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Customize timestamp format for JSON
			if a.Key == slog.TimeKey {
				return slog.Attr{
					Key:   slog.TimeKey,
					Value: slog.StringValue(a.Value.Time().Format("2006-01-02 15:04:05")),
				}
			}
			return a
		},
	}

	// Create file logger with JSON format
	fileHandler := slog.NewJSONHandler(file, opts)
	fileLogger := slog.New(fileHandler)

	return &Logger{
		fileLogger: fileLogger,
		config:     cfg,
		service:    "",
	}
}

// WithService returns a logger with service context
func (l *Logger) WithService(service string) *Logger {
	return &Logger{
		fileLogger: l.fileLogger.With("service", service),
		config:     l.config,
		service:    service,
	}
}

// WithContext returns a logger with additional context
func (l *Logger) WithContext(ctx map[string]interface{}) *Logger {
	args := make([]interface{}, 0, len(ctx)*2)
	for k, v := range ctx {
		args = append(args, k, v)
	}
	return &Logger{
		fileLogger: l.fileLogger.With(args...),
		config:     l.config,
		service:    l.service,
	}
}

// logToFile logs to the file using JSON format
func (l *Logger) logToFile(level slog.Level, msg string) {
	switch level {
	case slog.LevelDebug:
		l.fileLogger.Debug(msg)
	case slog.LevelInfo:
		l.fileLogger.Info(msg)
	case slog.LevelWarn:
		l.fileLogger.Warn(msg)
	case slog.LevelError:
		l.fileLogger.Error(msg)
	}
}

// logToCLI logs to the CLI using pretty format with colors
func (l *Logger) logToCLI(level slog.Level, msg string) {
	if !l.config.Show {
		return
	}

	timestamp := time.Now().Format("15:04:05")
	service := l.service
	if service == "" {
		service = "main"
	}

	var levelStr string
	var coloredLevel string

	switch level {
	case slog.LevelDebug:
		levelStr = " DEBUG"
		coloredLevel = pterm.FgGray.Sprint(levelStr)
	case slog.LevelInfo:
		levelStr = " INFO "
		coloredLevel = pterm.FgCyan.Sprint(levelStr)
	case slog.LevelWarn:
		levelStr = " WARN "
		coloredLevel = pterm.FgYellow.Sprint(levelStr)
	case slog.LevelError:
		levelStr = "ERROR"
		coloredLevel = pterm.FgRed.Sprint(levelStr)
	default:
		levelStr = "UNFORESEEN"
		coloredLevel = pterm.FgDefault.Sprint(levelStr)
	}

	// Format: timestamp [ STATE ] - Service - Log Message
	formattedMsg := fmt.Sprintf("%s [%s] - %s - %s",
		pterm.FgLightBlue.Sprint(timestamp),
		coloredLevel,
		pterm.FgMagenta.Sprint(service),
		msg,
	)

	fmt.Println(formattedMsg)
}

// log handles logging to both file and CLI
func (l *Logger) log(level slog.Level, msg string, args ...interface{}) {
	formattedMsg := fmt.Sprintf(msg, args...)
	
	// Always log to file
	l.logToFile(level, formattedMsg)
	
	// Log to CLI if enabled
	l.logToCLI(level, formattedMsg)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.log(slog.LevelDebug, msg, args...)
}

// Info logs an info message
func (l *Logger) Info(msg string, args ...interface{}) {
	l.log(slog.LevelInfo, msg, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.log(slog.LevelWarn, msg, args...)
}

// Error logs an error message
func (l *Logger) Error(msg string, args ...interface{}) {
	l.log(slog.LevelError, msg, args...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.log(slog.LevelError, msg, args...)
	os.Exit(1)
}

// Load logs a loading message (alias for Info with specific formatting)
func (l *Logger) Load(msg string, args ...interface{}) {
	l.log(slog.LevelInfo, "🔄 "+msg, args...)
}

// Ok logs a success message (alias for Info with specific formatting)
func (l *Logger) Ok(msg string, args ...interface{}) {
	l.log(slog.LevelInfo, "✅ "+msg, args...)
} 