package watcher

func (w *Watcher) Debug(message string) {
	w.logger.Debug(message)
	w.pool.AddMessage(message)
}

func (w *Watcher) Info(message string) {
	w.logger.Info(message)
	w.pool.AddMessage(message)
}

func (w *Watcher) Warning(message string) {
	w.logger.Warning(message)
	w.pool.AddMessage(message)
}

func (w *Watcher) Error(message string) {
	w.logger.Error(message)
	w.pool.AddMessage(message)
}

func (w *Watcher) Unforeseen(message string) {
	w.logger.Fatal(message)
	w.pool.AddMessage(message)
}
