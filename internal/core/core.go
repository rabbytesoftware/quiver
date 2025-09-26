package core

// ? Core DI Container
// ? This is the initiator for core services (Dependency Injection)

import (
	"github.com/rabbytesoftware/quiver/internal/core/config"
	"github.com/rabbytesoftware/quiver/internal/core/metadata"
	"github.com/rabbytesoftware/quiver/internal/core/watcher"
)

type Core struct {
	metadata *metadata.Metadata
	config *config.Config
	watcher *watcher.Watcher
}

// ? Although metadata and config are independent services (global state),
// ? we initiate them here to avoid circular dependencies and to ensure
// ? data is loaded only once and on a controlled manner.

func Init() *Core {
	return &Core{
		metadata: metadata.Get(),
		config: config.Get(),
		watcher: watcher.NewWatcherService(),
	}
}

func (c *Core) GetMetadata() *metadata.Metadata {
	return c.metadata
}

func (c *Core) GetConfig() *config.Config {
	return c.config
}

func (c *Core) GetWatcher() *watcher.Watcher {
	return c.watcher
}
