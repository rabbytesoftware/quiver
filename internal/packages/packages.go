package packages

import (
	"github.com/rabbytesoftware/quiver/internal/logger"
	"github.com/rabbytesoftware/quiver/internal/packages/core"
)

// ArrowsServer is an alias for the core.Manager for backward compatibility
type ArrowsServer = core.Manager

// NewArrowsServer creates a new arrows server using the core manager directly
func NewArrowsServer(repositories []string, installDir, dbPath string, logger *logger.Logger) *ArrowsServer {
	return core.NewManager(repositories, installDir, dbPath, logger)
}
