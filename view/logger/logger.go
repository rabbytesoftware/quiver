package logger

import (
	"fmt"
	"time"
)

var It *Logger

type Logger struct {}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) log(level LogLevel, msg string, attributes ...any) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   fmt.Sprintf(msg, attributes...),
	}

	fmt.Println(entry.FormattedString())
}

func (l *Logger) Debug(msg string, attributes ...any) {
	l.log(Debug, msg)
}

func (l *Logger) Info(msg string, attributes ...any) {
	l.log(Info, msg)
}

func (l *Logger) Ok(msg string, attributes ...any) {
	l.log(Ok, msg)
}

func (l *Logger) Load(msg string, attributes ...any) {
	l.log(Load, msg)
}

func (l *Logger) Warn(msg string, attributes ...any) {
	l.log(Warn, msg)
}

func (l *Logger) Error(msg string, attributes ...any) {
	l.log(Error, msg)
}

func (l *Logger) Fatal(msg string, attributes ...any) {
	l.log(Fatal, msg)
}
