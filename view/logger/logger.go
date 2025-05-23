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

func (l *Logger) log(level LogLevel, service string, msg string,  attributes ...any) {
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Service:   service,
		Message:   fmt.Sprintf(msg, attributes...),
	}
	
	err := SaveLogToFile("./logs", entry)

	if err != nil {
		fmt.Println("[Could not save log to file] ", err)
		return
	}
	
	fmt.Println(entry.FormattedString())
}

func (l *Logger) Debug(service string, msg string, attributes ...any) {
	l.log(Debug, service, msg, attributes...)
}

func (l *Logger) Info(service string, msg string, attributes ...any) {
	l.log(Info, service, msg, attributes...)
}

func (l *Logger) Ok(service string, msg string, attributes ...any) {
	l.log(Ok, service, msg, attributes...)
}

func (l *Logger) Load(service string, msg string, attributes ...any) {
	l.log(Load, service, msg, attributes...)
}

func (l *Logger) Warn(service string, msg string, attributes ...any) {
	l.log(Warn, service, msg, attributes...)
}

func (l *Logger) Error(service string, msg string, attributes ...any) {
	l.log(Error, service, msg, attributes...)
}

func (l *Logger) Fatal(service string, msg string, attributes ...any) {
	l.log(Fatal, service, msg, attributes...)
}
