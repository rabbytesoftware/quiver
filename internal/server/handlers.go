package server

import (
	"github.com/rabbytesoftware/quiver/internal/logger"
	"github.com/rabbytesoftware/quiver/internal/netbridge"
	"github.com/rabbytesoftware/quiver/internal/packages"
	"github.com/rabbytesoftware/quiver/internal/server/handlers/arrows"
	"github.com/rabbytesoftware/quiver/internal/server/handlers/health"
	netbridgeHandlers "github.com/rabbytesoftware/quiver/internal/server/handlers/netbridge"
	packageHandlers "github.com/rabbytesoftware/quiver/internal/server/handlers/packages"
	"github.com/rabbytesoftware/quiver/internal/server/handlers/repositories"
	serverHandlers "github.com/rabbytesoftware/quiver/internal/server/handlers/server"
)

// Handlers contains all HTTP handlers organized by subsystem
type Handlers struct {
	Health       *health.Handler
	Packages     *packageHandlers.Handler
	Arrows       *arrows.Handler
	Repositories *repositories.Handler
	Server       *serverHandlers.Handler
	Netbridge    *netbridgeHandlers.Handler
}

// NewHandlers creates a new handlers instance with all subsystem handlers
func NewHandlers(pkgManager *packages.ArrowsServer, netbridgeInstance *netbridge.Netbridge, logger *logger.Logger) *Handlers {
	return &Handlers{
		Health:       health.NewHandler(logger),
		Packages:     packageHandlers.NewHandler(pkgManager, logger),
		Arrows:       arrows.NewHandler(pkgManager, logger),
		Repositories: repositories.NewHandler(pkgManager, logger),
		Server:       serverHandlers.NewHandler(pkgManager, logger),
		Netbridge:    netbridgeHandlers.NewHandler(netbridgeInstance, logger),
	}
} 