package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rabbytesoftware/quiver/internal/logger"
	"github.com/rabbytesoftware/quiver/packages"
)

// Handlers contains all HTTP handlers
type Handlers struct {
	pkgManager *packages.ArrowsServer
	logger     *logger.Logger
}

// NewHandlers creates a new handlers instance
func NewHandlers(pkgManager *packages.ArrowsServer, logger *logger.Logger) *Handlers {
	return &Handlers{
		pkgManager: pkgManager,
		logger:     logger.WithService("handlers"),
	}
}

// HealthHandler handles health check requests
func (h *Handlers) HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":  "ok",
		"service": "quiver",
		"version": "1.0.0",
	}
	h.writeJSON(w, http.StatusOK, response)
}

// ListPackagesHandler handles listing all packages
func (h *Handlers) ListPackagesHandler(w http.ResponseWriter, r *http.Request) {
	packages := h.pkgManager.Packages
	response := map[string]interface{}{
		"packages": packages,
		"count":    len(packages),
	}
	h.writeJSON(w, http.StatusOK, response)
}

// GetPackageHandler handles getting a specific package
func (h *Handlers) GetPackageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	pkg, exists := h.pkgManager.Packages[id]
	if !exists {
		h.writeError(w, http.StatusNotFound, "Package not found")
		return
	}

	h.writeJSON(w, http.StatusOK, pkg)
}

// StartPackageHandler handles starting a package
func (h *Handlers) StartPackageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// TODO: Implement package starting logic
	_, exists := h.pkgManager.Packages[id]
	if !exists {
		h.writeError(w, http.StatusNotFound, "Package not found")
		return
	}

	h.logger.Info("Starting package: %s", id)
	response := map[string]string{
		"message": "Package started successfully",
	}
	h.writeJSON(w, http.StatusOK, response)
}

// StopPackageHandler handles stopping a package
func (h *Handlers) StopPackageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// TODO: Implement package stopping logic
	_, exists := h.pkgManager.Packages[id]
	if !exists {
		h.writeError(w, http.StatusNotFound, "Package not found")
		return
	}

	h.logger.Info("Stopping package: %s", id)
	response := map[string]string{
		"message": "Package stopped successfully",
	}
	h.writeJSON(w, http.StatusOK, response)
}

// PackageStatusHandler handles getting package status
func (h *Handlers) PackageStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, exists := h.pkgManager.Packages[id]
	if !exists {
		h.writeError(w, http.StatusNotFound, "Package not found")
		return
	}

	// TODO: Implement proper status checking
	status := map[string]interface{}{
		"id":     id,
		"status": "stopped", // Default status
	}
	h.writeJSON(w, http.StatusOK, status)
}

// ServerInfoHandler handles server info requests
func (h *Handlers) ServerInfoHandler(w http.ResponseWriter, r *http.Request) {
	info := map[string]interface{}{
		"name":    "Quiver",
		"version": "1.0.0",
		"description": "Game Server Management Platform",
	}
	h.writeJSON(w, http.StatusOK, info)
}

// ServerStatusHandler handles server status requests
func (h *Handlers) ServerStatusHandler(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"status":           "running",
		"packages_loaded": len(h.pkgManager.Packages),
	}
	h.writeJSON(w, http.StatusOK, status)
}

// writeJSON writes a JSON response
func (h *Handlers) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// writeError writes an error response
func (h *Handlers) writeError(w http.ResponseWriter, status int, message string) {
	errorResponse := map[string]string{
		"error": message,
	}
	h.writeJSON(w, status, errorResponse)
} 