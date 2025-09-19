package logstream

import (
	"github.com/rabbytesoftware/quiver/internal/infrastructure/watcher/pool"
	"github.com/sirupsen/logrus"
)

type LogService interface {
	SetLevel(level logrus.Level)
	GetLevel() logrus.Level
	WithFields(fields logrus.Fields) *logrus.Entry
	WithField(key string, value interface{}) *logrus.Entry
	Subscribe(callback pool.Subscriber)
	GetSubscriberCount() int

	Debug(message string)
	Info(message string)
	Warn(message string)
	Error(message string)
	Unforeseen(message string)
}
