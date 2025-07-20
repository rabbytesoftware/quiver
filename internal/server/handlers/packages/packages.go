package packages

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rabbytesoftware/quiver/internal/logger"
	"github.com/rabbytesoftware/quiver/internal/packages"
	"github.com/rabbytesoftware/quiver/internal/packages/types"
	"github.com/rabbytesoftware/quiver/internal/server/response"
)

// Handler handles package-related HTTP requests (legacy compatibility)
type Handler struct {
	pkgManager *packages.ArrowsServer
	logger     *logger.Logger
}

// NewHandler creates a new packages handler instance
func NewHandler(pkgManager *packages.ArrowsServer, logger *logger.Logger) *Handler {
	return &Handler{
		pkgManager: pkgManager,
		logger:     logger.WithService("packages-handler"),
	}
}

// ListPackages handles listing all available packages from repositories
func (h *Handler) ListPackages(c *gin.Context) {
	var allFiles []string
	
	// Get repository directories from package manager
	repositories := h.pkgManager.GetRepositories()
	
	// Read each repository directory directly
	for _, repoPath := range repositories {
		files, err := os.ReadDir(repoPath)
		if err != nil {
			h.logger.Warn("Failed to read repository directory %s: %v", repoPath, err)
			continue
		}
		
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			
			// Check if it's a YAML file and remove extension
			filename := file.Name()
			if strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml") {
				nameWithoutExt := strings.TrimSuffix(strings.TrimSuffix(filename, ".yaml"), ".yml")
				allFiles = append(allFiles, nameWithoutExt)
			}
		}
	}
	
	responseData := gin.H{
		"packages": allFiles,
		"count":    len(allFiles),
	}

	response.Success(c, "Available packages retrieved successfully", responseData)
}

// GetPackage handles getting a specific package information
func (h *Handler) GetPackage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Package ID is required")
		return
	}

	// Check if package is installed
	installed := h.pkgManager.GetInstalledArrows()
	pkg, exists := installed[id]
	if !exists {
		response.NotFound(c, "Package")
		return
	}

	response.Success(c, "Package information retrieved successfully", pkg)
}

// StartPackage handles starting a package (arrow execution)
func (h *Handler) StartPackage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Package ID is required")
		return
	}

	// Check if package is installed
	installed := h.pkgManager.GetInstalledArrows()
	_, exists := installed[id]
	if !exists {
		response.NotFound(c, "Package")
		return
	}

	h.logger.Info("Starting package: %s", id)

	// Execute the arrow (equivalent to starting the package)
	err := h.pkgManager.ExecuteArrow(id, nil)
	if err != nil {
		h.logger.Error("Failed to start package %s: %v", id, err)
		response.BadRequest(c, "Failed to start package", err.Error())
		return
	}

	responseData := gin.H{
		"package": id,
		"action":  "started",
	}

	response.Success(c, "Package started successfully", responseData)
}

// StopPackage handles stopping a package
func (h *Handler) StopPackage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Package ID is required")
		return
	}

	// Check if package is installed
	installed := h.pkgManager.GetInstalledArrows()
	_, exists := installed[id]
	if !exists {
		response.NotFound(c, "Package")
		return
	}

	h.logger.Info("Stopping package: %s", id)

	// Parse optional query parameters for stopping behavior
	gracefulStr := c.DefaultQuery("graceful", "true")
	timeoutStr := c.DefaultQuery("timeout", "30")

	graceful, err := strconv.ParseBool(gracefulStr)
	if err != nil {
		graceful = true // Default to graceful shutdown
	}

	timeoutSecs, err := strconv.Atoi(timeoutStr)
	if err != nil {
		timeoutSecs = 30 // Default timeout
	}
	timeout := time.Duration(timeoutSecs) * time.Second

	// Get current status before attempting to stop
	previousStatus, err := h.pkgManager.GetArrowStatus(id)
	if err != nil {
		h.logger.Error("Failed to get package status %s: %v", id, err)
		response.InternalServerError(c, "Failed to get package status", err.Error())
		return
	}

	// Actually stop the arrow processes
	err = h.pkgManager.StopArrow(id, graceful, timeout)
	if err != nil {
		h.logger.Error("Failed to stop package %s: %v", id, err)
		response.InternalServerError(c, "Failed to stop package", err.Error())
		return
	}

	responseData := gin.H{
		"package":         id,
		"action":          "stopped",
		"previous_status": string(previousStatus),
		"graceful":        graceful,
		"timeout_seconds": timeoutSecs,
	}

	response.Success(c, "Package stopped successfully", responseData)
}

// GetPackageStatus handles getting package status
func (h *Handler) GetPackageStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Package ID is required")
		return
	}

	status, err := h.pkgManager.GetArrowStatus(id)
	if err != nil {
		h.logger.Error("Failed to get package status %s: %v", id, err)
		response.NotFound(c, "Package")
		return
	}

	// Get process information for enhanced status
	processes, err := h.pkgManager.GetArrowProcesses(id)
	if err != nil {
		h.logger.Warn("Failed to get processes for package %s: %v", id, err)
		processes = nil // Continue without process info
	}

	responseData := gin.H{
		"id":                   id,
		"status":               status,
		"has_running_processes": h.pkgManager.HasRunningProcesses(id),
		"process_count":        len(processes),
	}

	// Only include detailed process info if requested
	if c.Query("include_processes") == "true" {
		responseData["processes"] = processes
	}

	response.Success(c, "Package status retrieved successfully", responseData)
}

// GetPackageProcesses handles getting running processes for a package
func (h *Handler) GetPackageProcesses(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "Package ID is required")
		return
	}

	// Check if package is installed
	installed := h.pkgManager.GetInstalledArrows()
	_, exists := installed[id]
	if !exists {
		response.NotFound(c, "Package")
		return
	}

	// Get running processes for this package
	processes, err := h.pkgManager.GetArrowProcesses(id)
	if err != nil {
		h.logger.Error("Failed to get processes for package %s: %v", id, err)
		response.InternalServerError(c, "Failed to get package processes", err.Error())
		return
	}

	responseData := gin.H{
		"package":   id,
		"processes": processes,
		"count":     len(processes),
		"has_running_processes": h.pkgManager.HasRunningProcesses(id),
	}

	response.Success(c, "Package processes retrieved successfully", responseData)
}

// GetInstalledPackages handles listing all installed packages
func (h *Handler) GetInstalledPackages(c *gin.Context) {
	installed := h.pkgManager.GetInstalledArrows()

	responseData := gin.H{
		"packages": installed,
		"count":    len(installed),
	}

	response.Success(c, "Installed packages retrieved successfully", responseData)
}

// GetPackagesByStatus handles getting packages filtered by status
func (h *Handler) GetPackagesByStatus(c *gin.Context) {
	statusParam := c.Query("status")
	if statusParam == "" {
		response.BadRequest(c, "Status query parameter is required")
		return
	}

	status := types.PackageStatus(statusParam)
	packages := h.pkgManager.GetArrowsByStatus(status)

	responseData := gin.H{
		"packages": packages,
		"count":    len(packages),
		"status":   status,
	}

	response.Success(c, "Packages filtered by status retrieved successfully", responseData)
}

// GetAllPackageStatuses handles listing all package statuses
func (h *Handler) GetAllPackageStatuses(c *gin.Context) {
	installed := h.pkgManager.GetInstalledArrows()
	statuses := make(map[string]types.PackageStatus)

	for name := range installed {
		status, err := h.pkgManager.GetArrowStatus(name)
		if err != nil {
			h.logger.Warn("Failed to get status for package %s: %v", name, err)
			statuses[name] = types.StatusError
		} else {
			statuses[name] = status
		}
	}

	responseData := gin.H{
		"statuses": statuses,
		"count":    len(statuses),
	}

	response.Success(c, "Package statuses retrieved successfully", responseData)
} 