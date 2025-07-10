package arrows

import "github.com/gin-gonic/gin"

// SetupRoutes configures the arrow-related routes for the given router group
func (h *Handler) SetupRoutes(router *gin.RouterGroup) {
	arrows := router.Group("/arrows")
	{
		// Arrow search and management
		arrows.GET("/search", h.SearchArrows)
		arrows.POST("/:name/install", h.InstallArrow)
		arrows.POST("/:name/execute", h.ExecuteArrow)
		arrows.DELETE("/:name/uninstall", h.UninstallArrow)
		arrows.PUT("/:name/update", h.UpdateArrow)
		arrows.POST("/:name/validate", h.ValidateArrow)
		
		// Arrow information and status
		arrows.GET("/installed", h.GetInstalledArrows)
		arrows.GET("/:name/status", h.GetArrowStatus)
		arrows.GET("/status/:status", h.GetArrowsByStatus)
		arrows.GET("/:name", h.GetArrowInfo)
		arrows.GET("/statuses", h.ListArrowStatuses)
	}
} 