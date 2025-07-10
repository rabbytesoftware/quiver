package netbridge

import "github.com/gin-gonic/gin"

// SetupRoutes configures the netbridge-related routes for the given router group
func (h *Handler) SetupRoutes(router *gin.RouterGroup) {
	netbridge := router.Group("/netbridge")
	{
		// Port management - JSON body style
		netbridge.POST("/open", h.OpenPort)
		netbridge.POST("/close", h.ClosePort)
		
		// Port management - URL parameter style (for CLI compatibility)
		netbridge.POST("/open/:port/:protocol", h.OpenPortByURL)
		netbridge.POST("/close/:port/:protocol", h.ClosePortByURL)
		
		// Alternative port management endpoints
		netbridge.POST("/port/:port/open", h.OpenPortByURL)
		netbridge.DELETE("/port/:port", h.ClosePortByURL)
		
		// Netbridge information and status
		netbridge.GET("/ports", h.ListOpenPorts)
		netbridge.GET("/status", h.GetStatus)
		
		// Public IP refresh - multiple endpoints for compatibility
		netbridge.POST("/refresh", h.RefreshPublicIP)
		netbridge.POST("/refresh-ip", h.RefreshPublicIP)
		
		// Auto port management
		netbridge.POST("/auto", h.OpenPortAuto)
		netbridge.POST("/auto/:port", h.OpenPortAutoByURL)
	}
} 