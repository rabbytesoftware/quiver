package watcher

import (
	"github.com/rabbytesoftware/quiver/internal/infrastructure/watcher/pool"
	"github.com/sirupsen/logrus"
)

func (w *Watcher) Debug(message string) {
	w.logger.Debug(message)
	w.pool.AddMessage(pool.Message{
		Level: logrus.DebugLevel, 
		Message: message,
	})
}

func (w *Watcher) Info(message string) {
	w.logger.Info(message)
	w.pool.AddMessage(pool.Message{
		Level: logrus.InfoLevel, 
		Message: message,
	})
}

func (w *Watcher) Warn(message string) {
	w.logger.Warning(message)
	w.pool.AddMessage(pool.Message{
		Level: logrus.WarnLevel, 
		Message: message,
	})
}

func (w *Watcher) Error(message string) {
	w.logger.Error(message)
	w.pool.AddMessage(pool.Message{
		Level: logrus.ErrorLevel, 
		Message: message,
	})
}

func (w *Watcher) Unforeseen(message string) {
	w.logger.Fatal(message)
	w.pool.AddMessage(pool.Message{
		Level: logrus.FatalLevel, 
		Message: message,
	})
}
