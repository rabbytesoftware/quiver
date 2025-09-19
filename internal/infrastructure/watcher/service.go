package watcher

import (
	"github.com/rabbytesoftware/quiver/internal/infrastructure/watcher/pool"
	"github.com/sirupsen/logrus"
)

// Watcher is an logging service, focus on an
// event sourcing approach to logging.

type Watcher struct {
	logger *logrus.Logger
	pool *pool.Pool
}

func NewWatcherService() *Watcher {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(nil)

	return &Watcher{
		logger: logger,
		pool: pool.NewPool(),
	}
}

func (w *Watcher) SetLevel(level logrus.Level) {
	logrus.SetLevel(level)
	w.logger.SetLevel(level)
}

func (w *Watcher) GetLevel() logrus.Level {
	return logrus.GetLevel()
}

func (w *Watcher) WithFields(fields logrus.Fields) *logrus.Entry {
	return w.logger.WithFields(fields)
}

func (w *Watcher) WithField(key string, value interface{}) *logrus.Entry {
	return w.logger.WithField(key, value)
}

func (w *Watcher) Subscribe(callback pool.Subscriber) {
	w.pool.Subscribe(callback)
}

func (w *Watcher) GetSubscriberCount() int {
	return w.pool.GetSubscriberCount()
}
