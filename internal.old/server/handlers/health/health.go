package health

import (
	"github.com/gin-gonic/gin"

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

// HealthCheck handles health check requests
func (h *Handler) HealthCheck(c *gin.Context) {
	h.logger.Debug("Health check requested")
	response.HealthCheck(c, "Quiver", metadata.Project.Version)
}

// ReadinessProbe handles readiness probe requests
func (h *Handler) ReadinessProbe(c *gin.Context) {
	h.logger.Debug("Readiness probe requested")
	
	// TODO: Add actual readiness checks (database, external services, etc.)
	readinessData := gin.H{
		"status": "ready",
		"checks": gin.H{
			"database": "ok",
			"packages": "ok",
		},
	}
	
	response.Success(c, "Service is ready", readinessData)
}

// LivenessProbe handles liveness probe requests
func (h *Handler) LivenessProbe(c *gin.Context) {
	h.logger.Debug("Liveness probe requested")
	
	livenessData := gin.H{
		"status": "alive",
		"uptime": "unknown", // TODO: Add actual uptime calculation
	}
	
	response.Success(c, "Service is alive", livenessData)
} 