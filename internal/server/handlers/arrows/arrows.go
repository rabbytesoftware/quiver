package arrows

import (
	"fmt"
	"runtime"

	"github.com/gin-gonic/gin"

	"github.com/rabbytesoftware/quiver/internal/logger"
	"github.com/rabbytesoftware/quiver/internal/packages"
	"github.com/rabbytesoftware/quiver/internal/packages/execution"
	"github.com/rabbytesoftware/quiver/internal/packages/manifest"
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

// ExecuteArrow handles executing an arrow's execute method
func (h *Handler) ExecuteArrow(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		response.BadRequest(c, "Arrow name is required")
		return
	}

	h.logger.Info("Executing arrow: %s", name)

	// Get optional variables from request
	var req struct {
		Variables map[string]string `json:"variables,omitempty"`
	}
	c.ShouldBindJSON(&req)

	err := h.pkgManager.ExecuteArrow(name, req.Variables)
	if err != nil {
		h.logger.Error("Failed to execute arrow %s: %v", name, err)
		response.BadRequest(c, "Arrow execution failed", err.Error())
		return
	}

	responseData := gin.H{
		"arrow": name,
	}

	response.Success(c, "Arrow execution started successfully", responseData)
}

// InitializeArrowMethod handles initializing an arrow method with netbridge processing
func (h *Handler) InitializeArrowMethod(c *gin.Context) {
	name := c.Param("name")
	method := c.Param("method")
	
	if name == "" || method == "" {
		response.BadRequest(c, "Arrow name and method are required")
		return
	}

	h.logger.Info("Initializing method %s for arrow: %s", method, name)

	// Get optional variables from request
	var req struct {
		Variables map[string]string `json:"variables,omitempty"`
	}
	c.ShouldBindJSON(&req)

	// Check if arrow is installed
	installed := h.pkgManager.GetInstalledArrows()
	pkg, exists := installed[name]
	if !exists {
		response.NotFound(c, "Arrow")
		return
	}

	// Load arrow manifest
	arrow, err := h.getArrowFromInstallation(pkg.InstallPath)
	if err != nil {
		h.logger.Error("Failed to load arrow %s: %v", name, err)
		response.InternalServerError(c, "Failed to load arrow", err.Error())
		return
	}

	// Prepare execution context
	finalVariables := make(map[string]string)
	for k, v := range pkg.Variables {
		finalVariables[k] = v
	}
	for k, v := range req.Variables {
		finalVariables[k] = v
	}

	ctx := &types.ExecutionContext{
		ArrowName:   name,
		InstallPath: pkg.InstallPath,
		Variables:   finalVariables,
	}

	// Get netbridge processing results for this arrow
	netbridgeResults, err := h.getNetbridgeResults(arrow, ctx)
	if err != nil {
		h.logger.Error("Failed to get netbridge results for %s: %v", name, err)
		response.InternalServerError(c, "Netbridge processing failed", err.Error())
		return
	}

	// Validate method exists and is supported
	if err := h.validateMethodSupport(arrow, method); err != nil {
		h.logger.Error("Method validation failed for %s.%s: %v", name, method, err)
		response.BadRequest(c, "Method not supported", err.Error())
		return
	}

	responseData := map[string]interface{}{
		"arrow":   name,
		"method":  method,
		"message": fmt.Sprintf("Method %s initialized successfully", method),
		"netbridge": map[string]interface{}{
			"variables":   len(netbridgeResults),
			"results":     netbridgeResults,
			"has_failures": func() bool {
				for _, result := range netbridgeResults {
					if !result.Success {
						return true
					}
				}
				return false
			}(),
		},
	}

	h.logger.Info("Method %s initialization completed for arrow %s", method, name)
	response.Success(c, "Method initialized with netbridge processing", responseData)
}

// GetArrowNetbridgeStatus returns the current netbridge status for an arrow
func (h *Handler) GetArrowNetbridgeStatus(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		response.BadRequest(c, "Arrow name is required")
		return
	}

	// Check if arrow is installed
	installed := h.pkgManager.GetInstalledArrows()
	pkg, exists := installed[name]
	if !exists {
		response.NotFound(c, "Arrow")
		return
	}

	// Load arrow manifest
	arrow, err := h.getArrowFromInstallation(pkg.InstallPath)
	if err != nil {
		h.logger.Error("Failed to load arrow %s: %v", name, err)
		response.InternalServerError(c, "Failed to load arrow", err.Error())
		return
	}

	// Create basic execution context to check netbridge variables
	ctx := &types.ExecutionContext{
		ArrowName:   name,
		InstallPath: pkg.InstallPath,
		Variables:   pkg.Variables,
	}

	// Get netbridge processing results
	netbridgeResults, err := h.getNetbridgeResults(arrow, ctx)
	if err != nil {
		h.logger.Error("Failed to get netbridge status for %s: %v", name, err)
		response.InternalServerError(c, "Failed to get netbridge status", err.Error())
		return
	}

	responseData := map[string]interface{}{
		"arrow": name,
		"netbridge": map[string]interface{}{
			"variables": len(netbridgeResults),
			"results":   netbridgeResults,
			"summary": map[string]int{
				"total":     len(netbridgeResults),
				"success":   0,
				"failed":    0,
			},
		},
	}

	// Calculate summary statistics
	summary := responseData["netbridge"].(map[string]interface{})["summary"].(map[string]int)
	for _, result := range netbridgeResults {
		if result.Success {
			summary["success"]++
		} else {
			summary["failed"]++
		}
	}

	h.logger.Info("Netbridge status retrieved for arrow %s: %d total, %d success, %d failed", 
		name, summary["total"], summary["success"], summary["failed"])
	
	response.Success(c, "Netbridge status retrieved successfully", responseData)
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

// Helper methods

// getArrowFromInstallation loads arrow manifest from installation directory
func (h *Handler) getArrowFromInstallation(installPath string) (manifest.ArrowInterface, error) {
	processor := manifest.NewProcessor(h.logger)
	return processor.LoadFromInstallation(installPath)
}

// getNetbridgeResults gets netbridge processing results using the execution engine
func (h *Handler) getNetbridgeResults(arrow manifest.ArrowInterface, ctx *types.ExecutionContext) ([]*execution.NetbridgeResult, error) {
	engine := execution.NewEngine(h.logger)
	defer engine.Cleanup()
	return engine.GetNetbridgeResults(arrow, ctx)
}

// validateMethodSupport validates that a method is supported on the current platform
func (h *Handler) validateMethodSupport(arrow manifest.ArrowInterface, methodType string) error {
	methods := arrow.GetMethods()
	methodMap := methods.GetMethod(methodType)
	
	if methodMap == nil {
		return fmt.Errorf("method %s not found", methodType)
	}

	osStr := runtime.GOOS
	archStr := runtime.GOARCH

	osMap, exists := methodMap[osStr]
	if !exists {
		return fmt.Errorf("method %s not supported on platform %s", methodType, osStr)
	}

	if _, exists := osMap[archStr]; !exists {
		// Check if any architecture is supported for this OS
		if len(osMap) == 0 {
			return fmt.Errorf("method %s not supported on platform %s", methodType, osStr)
		}
		// Architecture fallback available
		h.logger.Warn("Architecture %s not directly supported for method %s, fallback available", archStr, methodType)
	}

	return nil
} 