package arrows

import (
	"github.com/gin-gonic/gin"

	"github.com/rabbytesoftware/quiver/internal/logger"
	"github.com/rabbytesoftware/quiver/internal/packages"
	"github.com/rabbytesoftware/quiver/internal/packages/types"
	"github.com/rabbytesoftware/quiver/internal/server/response"
)

// Handler handles arrow-related HTTP requests
type Handler struct {
	pkgManager *packages.ArrowsServer
	logger     *logger.Logger
}

// NewHandler creates a new arrows handler instance
func NewHandler(pkgManager *packages.ArrowsServer, logger *logger.Logger) *Handler {
	return &Handler{
		pkgManager: pkgManager,
		logger:     logger.WithService("arrows-handler"),
	}
}

// InstallArrowRequest represents the request payload for arrow installation
type InstallArrowRequest struct {
	Variables  map[string]string `json:"variables,omitempty"`
	Repository string            `json:"repository,omitempty"`
}

// ExecuteArrowRequest represents the request payload for arrow execution
type ExecuteArrowRequest struct {
	Variables map[string]string `json:"variables,omitempty"`
}

// UpdateArrowRequest represents the request payload for arrow updates
type UpdateArrowRequest struct {
	Repository string `json:"repository,omitempty"`
}

// SearchArrows handles searching for arrows across repositories
func (h *Handler) SearchArrows(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		response.BadRequest(c, "Query parameter 'q' is required")
		return
	}

	arrows, err := h.pkgManager.SearchArrows(query)
	if err != nil {
		h.logger.Error("Failed to search arrows: %v", err)
		response.InternalServerError(c, "Failed to search arrows", err.Error())
		return
	}

	responseData := gin.H{
		"arrows": arrows,
		"count":  len(arrows),
		"query":  query,
	}

	response.Success(c, "Search completed successfully", responseData)
}

// InstallArrow handles arrow installation
func (h *Handler) InstallArrow(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		response.BadRequest(c, "Arrow name is required")
		return
	}

	var req InstallArrowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Variables are optional, so we'll ignore binding errors
		req.Variables = make(map[string]string)
	}

	// Check for repository specification in query parameter (backwards compatibility)
	if repo := c.Query("repository"); repo != "" {
		req.Repository = repo
	}

	// Build full package specification
	fullName := name
	if req.Repository != "" {
		fullName = req.Repository + "@" + name
	}

	h.logger.Info("Installing arrow: %s (full name: %s)", name, fullName)

	err := h.pkgManager.InstallArrow(fullName, req.Variables)
	if err != nil {
		h.logger.Error("Failed to install arrow %s: %v", fullName, err)
		response.BadRequest(c, "Failed to install arrow", err.Error())
		return
	}

	responseData := gin.H{
		"arrow":      name,
		"repository": req.Repository,
		"variables":  req.Variables,
	}

	response.Created(c, "Arrow installed successfully", responseData)
}

// ExecuteArrow handles arrow execution
func (h *Handler) ExecuteArrow(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		response.BadRequest(c, "Arrow name is required")
		return
	}

	var req ExecuteArrowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Variables are optional, so we'll ignore binding errors
		req.Variables = make(map[string]string)
	}

	h.logger.Info("Executing arrow: %s", name)

	err := h.pkgManager.ExecuteArrow(name, req.Variables)
	if err != nil {
		h.logger.Error("Failed to execute arrow %s: %v", name, err)
		response.BadRequest(c, "Failed to execute arrow", err.Error())
		return
	}

	responseData := gin.H{
		"arrow":     name,
		"variables": req.Variables,
	}

	response.Success(c, "Arrow executed successfully", responseData)
}

// UninstallArrow handles arrow uninstallation
func (h *Handler) UninstallArrow(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		response.BadRequest(c, "Arrow name is required")
		return
	}

	h.logger.Info("Uninstalling arrow: %s", name)

	err := h.pkgManager.UninstallArrow(name)
	if err != nil {
		h.logger.Error("Failed to uninstall arrow %s: %v", name, err)
		response.BadRequest(c, "Failed to uninstall arrow", err.Error())
		return
	}

	responseData := gin.H{
		"arrow": name,
	}

	response.Success(c, "Arrow uninstalled successfully", responseData)
}

// UpdateArrow handles arrow updates
func (h *Handler) UpdateArrow(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		response.BadRequest(c, "Arrow name is required")
		return
	}

	var req UpdateArrowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Repository is optional, so we'll ignore binding errors
	}

	// Check for repository specification in query parameter (backwards compatibility)
	if repo := c.Query("repository"); repo != "" {
		req.Repository = repo
	}

	// Build full package specification
	fullName := name
	if req.Repository != "" {
		fullName = req.Repository + "@" + name
	}

	h.logger.Info("Updating arrow: %s (full name: %s)", name, fullName)

	err := h.pkgManager.UpdateArrow(fullName)
	if err != nil {
		h.logger.Error("Failed to update arrow %s: %v", fullName, err)
		response.BadRequest(c, "Failed to update arrow", err.Error())
		return
	}

	responseData := gin.H{
		"arrow":      name,
		"repository": req.Repository,
	}

	response.Success(c, "Arrow updated successfully", responseData)
}

// ValidateArrow handles arrow validation
func (h *Handler) ValidateArrow(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		response.BadRequest(c, "Arrow name is required")
		return
	}

	h.logger.Info("Validating arrow: %s", name)

	err := h.pkgManager.ValidateArrow(name)
	if err != nil {
		h.logger.Error("Failed to validate arrow %s: %v", name, err)
		response.BadRequest(c, "Arrow validation failed", err.Error())
		return
	}

	responseData := gin.H{
		"arrow": name,
	}

	response.Success(c, "Arrow validation completed successfully", responseData)
}

// GetInstalledArrows handles listing installed arrows
func (h *Handler) GetInstalledArrows(c *gin.Context) {
	installed := h.pkgManager.GetInstalledArrows()

	responseData := gin.H{
		"installed": installed,
		"count":     len(installed),
	}

	response.Success(c, "Installed arrows retrieved successfully", responseData)
}

// GetArrowStatus handles getting arrow status
func (h *Handler) GetArrowStatus(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		response.BadRequest(c, "Arrow name is required")
		return
	}

	status, err := h.pkgManager.GetArrowStatus(name)
	if err != nil {
		h.logger.Error("Failed to get arrow status %s: %v", name, err)
		response.NotFound(c, "Arrow")
		return
	}

	responseData := gin.H{
		"arrow":  name,
		"status": status,
	}

	response.Success(c, "Arrow status retrieved successfully", responseData)
}

// GetArrowsByStatus handles getting arrows filtered by status
func (h *Handler) GetArrowsByStatus(c *gin.Context) {
	statusParam := c.Query("status")
	if statusParam == "" {
		response.BadRequest(c, "Status query parameter is required")
		return
	}

	status := types.PackageStatus(statusParam)
	arrows := h.pkgManager.GetArrowsByStatus(status)

	responseData := gin.H{
		"arrows": arrows,
		"count":  len(arrows),
		"status": status,
	}

	response.Success(c, "Arrows filtered by status retrieved successfully", responseData)
}

// GetArrowInfo handles getting detailed arrow information
func (h *Handler) GetArrowInfo(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		response.BadRequest(c, "Arrow name is required")
		return
	}

	installed := h.pkgManager.GetInstalledArrows()
	arrow, exists := installed[name]
	if !exists {
		response.NotFound(c, "Arrow")
		return
	}

	response.Success(c, "Arrow information retrieved successfully", arrow)
}

// ListArrowStatuses handles listing all arrow statuses
func (h *Handler) ListArrowStatuses(c *gin.Context) {
	installed := h.pkgManager.GetInstalledArrows()
	statuses := make(map[string]types.PackageStatus)

	for name := range installed {
		status, err := h.pkgManager.GetArrowStatus(name)
		if err != nil {
			h.logger.Warn("Failed to get status for arrow %s: %v", name, err)
			statuses[name] = types.StatusError
		} else {
			statuses[name] = status
		}
	}

	responseData := gin.H{
		"statuses": statuses,
		"count":    len(statuses),
	}

	response.Success(c, "Arrow statuses retrieved successfully", responseData)
} 