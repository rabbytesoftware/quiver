package health

import (
	"net/http"

	"github.com/rabbytesoftware/quiver/internal/logger"
	"github.com/rabbytesoftware/quiver/internal/metadata"
	"github.com/rabbytesoftware/quiver/internal/server/response"
)

// Handler handles health-related HTTP requests
type Handler struct {
	logger *logger.Logger
}

// NewHandler creates a new health handler instance
func NewHandler(logger *logger.Logger) *Handler {
	return &Handler{
		logger: logger.WithService("health-handler"),
	}
}

// HealthHandler handles health check requests
func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	responseData := map[string]interface{}{
		"status":  "ok",
		"version": metadata.Project.Version,
	}
	response.WriteJSON(w, http.StatusOK, responseData)
} 