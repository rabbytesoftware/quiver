package packages

import "github.com/gin-gonic/gin"

// SetupRoutes configures the package-related routes for the given router group
func (h *Handler) SetupRoutes(router *gin.RouterGroup) {
	packages := router.Group("/packages")
	{
		// Package management
		packages.GET("/", h.ListPackages)
		packages.GET("/:id", h.GetPackage)
		packages.POST("/:id/start", h.StartPackage)
		packages.POST("/:id/stop", h.StopPackage)
		packages.GET("/:id/status", h.GetPackageStatus)
		packages.GET("/:id/processes", h.GetPackageProcesses)
		
		// Package listings and status
		packages.GET("/installed", h.GetInstalledPackages)
		packages.GET("/status/:status", h.GetPackagesByStatus)
		packages.GET("/statuses", h.GetAllPackageStatuses)
		
		// Execution status tracking
		packages.GET("/:id/execution/status", h.GetPackageExecutionStatus)
		packages.GET("/executions/status", h.GetAllPackageExecutionStatuses)
	}
} 