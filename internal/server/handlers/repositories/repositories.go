package repositories

import (
	"github.com/gin-gonic/gin"

	"github.com/rabbytesoftware/quiver/internal/logger"
	"github.com/rabbytesoftware/quiver/internal/packages"
	"github.com/rabbytesoftware/quiver/internal/server/response"
)

// Handler handles repository-related HTTP requests
type Handler struct {
	pkgManager *packages.ArrowsServer
	logger     *logger.Logger
}

// NewHandler creates a new repositories handler instance
func NewHandler(pkgManager *packages.ArrowsServer, logger *logger.Logger) *Handler {
	return &Handler{
		pkgManager: pkgManager,
		logger:     logger.WithService("repositories-handler"),
	}
}

// AddRepositoryRequest represents the request payload for adding a repository
type AddRepositoryRequest struct {
	Repository string `json:"repository" binding:"required"`
	Type       string `json:"type,omitempty"` // "local" or "remote"
	Name       string `json:"name,omitempty"` // Optional alias name
}

// RemoveRepositoryRequest represents the request payload for removing a repository
type RemoveRepositoryRequest struct {
	Repository string `json:"repository" binding:"required"`
}

// AddRepository handles adding a new repository
func (h *Handler) AddRepository(c *gin.Context) {
	var req AddRepositoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	h.logger.Info("Adding repository: %s", req.Repository)

	h.pkgManager.AddRepository(req.Repository)

	responseData := gin.H{
		"repository": req.Repository,
		"type":       req.Type,
		"name":       req.Name,
	}

	response.Created(c, "Repository added successfully", responseData)
}

// RemoveRepository handles removing a repository
func (h *Handler) RemoveRepository(c *gin.Context) {
	var req RemoveRepositoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body", err.Error())
		return
	}

	h.logger.Info("Removing repository: %s", req.Repository)

	h.pkgManager.RemoveRepository(req.Repository)

	responseData := gin.H{
		"repository": req.Repository,
	}

	response.Success(c, "Repository removed successfully", responseData)
}

// GetRepositories handles listing all repositories
func (h *Handler) GetRepositories(c *gin.Context) {
	repositories := h.pkgManager.GetRepositories()

	responseData := gin.H{
		"repositories": repositories,
		"count":        len(repositories),
	}

	response.Success(c, "Repositories retrieved successfully", responseData)
}

// GetRepositoryInfo handles getting detailed information about a specific repository
func (h *Handler) GetRepositoryInfo(c *gin.Context) {
	repository := c.Param("repository")
	if repository == "" {
		response.BadRequest(c, "Repository path is required")
		return
	}

	repositories := h.pkgManager.GetRepositories()
	
	// Check if repository exists
	found := false
	for _, repo := range repositories {
		if repo == repository {
			found = true
			break
		}
	}
	
	if !found {
		response.NotFound(c, "Repository")
		return
	}

	// For now, just return basic information
	// TODO: Extend this to return more detailed repository information
	responseData := gin.H{
		"repository": repository,
		"type":       "unknown", // TODO: Determine type from repository analysis
		"status":     "active",
	}

	response.Success(c, "Repository information retrieved successfully", responseData)
}

// SearchRepositories handles searching through repositories
func (h *Handler) SearchRepositories(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		response.BadRequest(c, "Query parameter 'q' is required")
		return
	}

	// Use the package manager search functionality across repositories
	arrows, err := h.pkgManager.SearchArrows(query)
	if err != nil {
		h.logger.Error("Failed to search repositories: %v", err)
		response.InternalServerError(c, "Failed to search repositories", err.Error())
		return
	}

	responseData := gin.H{
		"query":   query,
		"results": arrows,
		"count":   len(arrows),
	}

	response.Success(c, "Repository search completed successfully", responseData)
}

// ValidateRepository handles validating a repository
func (h *Handler) ValidateRepository(c *gin.Context) {
	repository := c.Param("repository")
	if repository == "" {
		response.BadRequest(c, "Repository path is required")
		return
	}

	repositories := h.pkgManager.GetRepositories()
	
	// Check if repository exists
	found := false
	for _, repo := range repositories {
		if repo == repository {
			found = true
			break
		}
	}
	
	if !found {
		response.NotFound(c, "Repository")
		return
	}

	// TODO: Implement actual repository validation logic
	// For now, we'll consider it valid if it exists in our list
	responseData := gin.H{
		"repository": repository,
		"valid":      true,
		"message":    "Repository validation passed",
	}

	response.Success(c, "Repository validation completed successfully", responseData)
} 