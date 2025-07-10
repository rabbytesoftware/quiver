package packages

import (
	"github.com/rabbytesoftware/quiver/internal/logger"
	"github.com/rabbytesoftware/quiver/internal/packages/server"
)

// ArrowsServer is an alias for the new server.ArrowsServer for backward compatibility
type ArrowsServer = server.ArrowsServer

// NewArrowsServer creates a new arrows server using the new architecture
func NewArrowsServer(repositories []string, installDir, dbPath string, logger *logger.Logger) *ArrowsServer {
	return server.NewArrowsServer(repositories, installDir, dbPath, logger)
}
