package server

import (
	"net/http"

	"github.com/rabbytesoftware/quiver/internal/logger"
	"github.com/rabbytesoftware/quiver/internal/metadata"
	"github.com/rabbytesoftware/quiver/internal/packages"
	"github.com/rabbytesoftware/quiver/internal/server/response"
)

// Handler handles server-related HTTP requests
type Handler struct {
	pkgManager *packages.ArrowsServer
	logger     *logger.Logger
}

// NewHandler creates a new server handler instance
func NewHandler(pkgManager *packages.ArrowsServer, logger *logger.Logger) *Handler {
	return &Handler{
		pkgManager: pkgManager,
		logger:     logger.WithService("server-handler"),
	}
}

// ServerInfoHandler handles server info requests
func (h *Handler) ServerInfoHandler(w http.ResponseWriter, r *http.Request) {
	info := map[string]interface{}{
		"version": metadata.Project.Version,
		"description": metadata.Project.Description,
	}
	response.WriteJSON(w, http.StatusOK, info)
}

// ServerStatusHandler handles server status requests
func (h *Handler) ServerStatusHandler(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"status":           "running",
		"packages_loaded": len(h.pkgManager.Packages),
	}
	response.WriteJSON(w, http.StatusOK, status)
} 