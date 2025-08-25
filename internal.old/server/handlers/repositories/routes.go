package repositories

import "github.com/gin-gonic/gin"

// SetupRoutes configures the repository-related routes for the given router group
func (h *Handler) SetupRoutes(router *gin.RouterGroup) {
	repositories := router.Group("/repositories")
	{
		// Repository management
		repositories.POST("/", h.AddRepository)
		repositories.DELETE("/", h.RemoveRepository)
		repositories.GET("/", h.GetRepositories)
		
		// Repository information and search
		repositories.GET("/:repository", h.GetRepositoryInfo)
		repositories.GET("/search", h.SearchRepositories)
		repositories.POST("/:repository/validate", h.ValidateRepository)
	}
} 