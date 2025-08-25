package health

import "github.com/gin-gonic/gin"

// SetupRoutes configures the health-related routes for the given router group
func (h *Handler) SetupRoutes(router *gin.RouterGroup) {
	// Health check endpoint
	router.GET("/health", h.HealthCheck)
	
	// Kubernetes style health checks
	router.GET("/ready", h.ReadinessProbe)
	router.GET("/live", h.LivenessProbe)
} 