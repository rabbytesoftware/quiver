package server

import (
	"github.com/rabbytesoftware/quiver/internal/logger"
	"github.com/rabbytesoftware/quiver/internal/packages"
	"github.com/rabbytesoftware/quiver/internal/server/handlers/arrows"
	"github.com/rabbytesoftware/quiver/internal/server/handlers/health"
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
}

// NewHandlers creates a new handlers instance with all subsystem handlers
func NewHandlers(pkgManager *packages.ArrowsServer, logger *logger.Logger) *Handlers {
	return &Handlers{
		Health:       health.NewHandler(logger),
		Packages:     packageHandlers.NewHandler(pkgManager, logger),
		Arrows:       arrows.NewHandler(pkgManager, logger),
		Repositories: repositories.NewHandler(pkgManager, logger),
		Server:       serverHandlers.NewHandler(pkgManager, logger),
	}
} 