package watcher

import (
	"fmt"

	"github.com/rabbytesoftware/quiver/internal/core/errors"
	"github.com/rabbytesoftware/quiver/internal/core/watcher/pool"
	"github.com/sirupsen/logrus"
)

func Debug(message string) {
	w.logger.Debug(message)
	w.pool.AddMessage(pool.Message{
		Level:   logrus.DebugLevel,
		Message: message,
	})
}

func Info(message string) {
	w.logger.Info(message)
	w.pool.AddMessage(pool.Message{
		Level:   logrus.InfoLevel,
		Message: message,
	})
}

func Warn(message string) {
	w.logger.Warning(message)
	w.pool.AddMessage(pool.Message{
		Level:   logrus.WarnLevel,
		Message: message,
	})
}

func Errorf(message string, args ...interface{}) {
	w.logger.Info(fmt.Sprintf(message, args...))
	w.pool.AddMessage(pool.Message{
		Level:   logrus.InfoLevel,
		Message: fmt.Sprintf(message, args...),
	})
}

// ? Error forces you to use the errors.Error type
func Error(message errors.Error) {
	w.logger.Error(message.Error())
	w.pool.AddMessage(pool.Message{
		Level:   logrus.ErrorLevel,
		Message: message.Error(),
	})
}

// ? Unforeseen forces you to use the errors.Error type
func Unforeseen(message errors.Error) {
	w.logger.Fatal(message.Error())
	w.pool.AddMessage(pool.Message{
		Level:   logrus.FatalLevel,
		Message: message.Error(),
	})
}
