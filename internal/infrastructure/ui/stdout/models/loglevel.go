package models

type LogLevel string

const (
	LogLevelDebug     LogLevel = "DEBUG"
	LogLevelInfo      LogLevel = "INFO"
	LogLevelWarning   LogLevel = "WARNING"
	LogLevelError     LogLevel = "ERROR"
	LogLevelUnforeseen LogLevel = "UNFORESEEN"
)

// LogLevelPriority defines the priority order for log levels
var LogLevelPriority = map[LogLevel]int{
	LogLevelDebug:     0,
	LogLevelInfo:      1,
	LogLevelWarning:   2,
	LogLevelError:     3,
	LogLevelUnforeseen: 4,
}
