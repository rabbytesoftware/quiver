package watcher

import (
	"os"
	"path/filepath"

	"github.com/rabbytesoftware/quiver/internal/core/config"
	"github.com/rabbytesoftware/quiver/internal/core/watcher/pool"
	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// ? Watcher is an logging service, focus on an
// ? event sourcing approach to logging.

type Watcher struct {
	logger *logrus.Logger
	pool   *pool.Pool
}

func NewWatcherService() *Watcher {
	watcherConfig := config.GetWatcher()

	return &Watcher{
		logger: initLogger(watcherConfig),
		pool:   pool.NewPool(),
	}
}

func (w *Watcher) SetLevel(level logrus.Level) {
	w.logger.SetLevel(level)
}

func (w *Watcher) GetLevel() logrus.Level {
	return w.logger.GetLevel()
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

func (w *Watcher) GetConfig() config.Watcher {
	return config.GetWatcher()
}

func (w *Watcher) IsEnabled() bool {
	return config.GetWatcher().Enabled
}

func initLogger(watcherConfig config.Watcher) *logrus.Logger {
	logger := logrus.New()

	level, err := logrus.ParseLevel(watcherConfig.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	if err := os.MkdirAll(watcherConfig.Folder, 0755); err != nil || !watcherConfig.Enabled {
		logger.SetOutput(os.Stderr)

		return logger
	}

	logFile := filepath.Join(watcherConfig.Folder, "quiver.log")

	logger.SetOutput(&lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    watcherConfig.MaxSize,
		MaxAge:     watcherConfig.MaxAge,
		MaxBackups: 3,
		Compress:   watcherConfig.Compress,
		LocalTime:  true,
	})

	return logger
}
