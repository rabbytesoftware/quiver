package server

import "github.com/gin-gonic/gin"

// SetupRoutes configures the server-related routes for the given router group
func (h *Handler) SetupRoutes(router *gin.RouterGroup) {
	server := router.Group("/server")
	{
		// Server information and monitoring
		server.GET("/info", h.GetServerInfo)
		server.GET("/status", h.GetServerStatus)
		server.GET("/health", h.GetServerHealth)
		server.GET("/metrics", h.GetServerMetrics)
		server.GET("/version", h.GetServerVersion)
	}
} 