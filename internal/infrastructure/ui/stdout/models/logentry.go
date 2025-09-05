package models

import "time"

type LogEntry struct {
	Timestamp time.Time
	Level     LogLevel
	Message   string
}
