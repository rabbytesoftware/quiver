package watcher

import (
	"fmt"

	"github.com/rabbytesoftware/quiver/internal/core/errors"
	"github.com/rabbytesoftware/quiver/internal/core/watcher/pool"
	"github.com/sirupsen/logrus"
)

func (w *Watcher) Debug(message string) {
	w.logger.Debug(message)
	w.pool.AddMessage(pool.Message{
		Level:   logrus.DebugLevel,
		Message: message,
	})
}

func (w *Watcher) Info(message string) {
	w.logger.Info(message)
	w.pool.AddMessage(pool.Message{
		Level:   logrus.InfoLevel,
		Message: message,
	})
}

func (w *Watcher) Warn(message string) {
	w.logger.Warning(message)
	w.pool.AddMessage(pool.Message{
		Level:   logrus.WarnLevel,
		Message: message,
	})
}

func (w *Watcher) Errorf(message string, args ...interface{}) {
	w.logger.Info(fmt.Sprintf(message, args...))
	w.pool.AddMessage(pool.Message{
		Level:   logrus.InfoLevel,
		Message: fmt.Sprintf(message, args...),
	})
}

// ? Error forces you to use the errors.Error type
func (w *Watcher) Error(message errors.Error) {
	w.logger.Error(message.Error())
	w.pool.AddMessage(pool.Message{
		Level:   logrus.ErrorLevel,
		Message: message.Error(),
	})
}

// ? Unforeseen forces you to use the errors.Error type
func (w *Watcher) Unforeseen(message errors.Error) {
	w.logger.Fatal(message)
	w.pool.AddMessage(pool.Message{
		Level:   logrus.FatalLevel,
		Message: message.Error(),
	})
}
