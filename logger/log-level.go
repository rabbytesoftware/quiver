package logger

type LogLevel int

const (
	Debug LogLevel = iota
	Info
	Ok
	Load
	Warn
	Error
	Fatal
)

func (l LogLevel) String() string {
	return [...]string{"DEBUG", "INFO", "OK", "LOAD", "WARN", "ERROR", "FATAL"}[l]
}

func (l LogLevel) Color() string {
	colors := [...]string{
		"\033[36m",    // Debug - Cyan
		"\033[37m",    // Info - White
		"\033[32m",    // Ok - Green
		"\033[34m",    // Load - Blue
		"\033[33m",    // Warn - Yellow
		"\033[31m",    // Error - Red
		"\033[35m",    // Fatal - Magenta
	}
	return colors[l]
}
