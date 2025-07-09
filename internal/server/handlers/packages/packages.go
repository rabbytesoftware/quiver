package packages

import (
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rabbytesoftware/quiver/internal/logger"
	"github.com/rabbytesoftware/quiver/internal/packages"
	"github.com/rabbytesoftware/quiver/internal/server/response"
)

// Handler handles package-related HTTP requests
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

// ListPackagesHandler handles listing all packages
func (h *Handler) ListPackagesHandler(w http.ResponseWriter, r *http.Request) {
	var allFiles []string
	
	// Get repository directories from package manager
	repositories := h.pkgManager.PackageManager.GetRepositories()
	
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
	
	responseData := map[string]interface{}{
		"packages": allFiles,
		"count":    len(allFiles),
	}
	response.WriteJSON(w, http.StatusOK, responseData)
}

// GetPackageHandler handles getting a specific package
func (h *Handler) GetPackageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	pkg, exists := h.pkgManager.Packages[id]
	if !exists {
		response.WriteError(w, http.StatusNotFound, "Package not found")
		return
	}

	response.WriteJSON(w, http.StatusOK, pkg)
}

// StartPackageHandler handles starting a package
func (h *Handler) StartPackageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// TODO: Implement package starting logic
	_, exists := h.pkgManager.Packages[id]
	if !exists {
		response.WriteError(w, http.StatusNotFound, "Package not found")
		return
	}

	h.logger.Info("Starting package: %s", id)
	responseData := map[string]string{
		"message": "Package started successfully",
	}
	response.WriteJSON(w, http.StatusOK, responseData)
}

// StopPackageHandler handles stopping a package
func (h *Handler) StopPackageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// TODO: Implement package stopping logic
	_, exists := h.pkgManager.Packages[id]
	if !exists {
		response.WriteError(w, http.StatusNotFound, "Package not found")
		return
	}

	h.logger.Info("Stopping package: %s", id)
	responseData := map[string]string{
		"message": "Package stopped successfully",
	}
	response.WriteJSON(w, http.StatusOK, responseData)
}

// PackageStatusHandler handles getting package status
func (h *Handler) PackageStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, exists := h.pkgManager.Packages[id]
	if !exists {
		response.WriteError(w, http.StatusNotFound, "Package not found")
		return
	}

	// TODO: Implement proper status checking
	status := map[string]interface{}{
		"id":     id,
		"status": "stopped", // Default status
	}
	response.WriteJSON(w, http.StatusOK, status)
} 