package logger

import (
	"fmt"
	"time"
)

type LogEntry struct {
	Timestamp time.Time
	Level     LogLevel
	Service   string
	Message   string
}

func (e LogEntry) String() string {
	return fmt.Sprintf("[%s] %s: %s",
		e.Timestamp.Format("2006-01-02 15:04:05"),
		e.Level.String(),
		e.Message,
	)
}

func (e LogEntry) FormattedString() string {
	return fmt.Sprintf("\033[1;37m%s \033[%s[%s]\033[0m %s",
		e.Timestamp.Format("2006-01-02 15:04:05"),
		e.Level.Color(),
		e.Level.String(),
		e.Message,
	)
}