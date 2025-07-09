package arrows

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rabbytesoftware/quiver/internal/logger"
	"github.com/rabbytesoftware/quiver/internal/packages"
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

// SearchArrowsHandler handles searching for arrows across repositories
func (h *Handler) SearchArrowsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	
	arrows, err := h.pkgManager.PackageManager.SearchArrows(query)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to search arrows: "+err.Error())
		return
	}

	responseData := map[string]interface{}{
		"arrows": arrows,
		"count":  len(arrows),
		"query":  query,
	}
	response.WriteJSON(w, http.StatusOK, responseData)
}

// InstallArrowHandler handles arrow installation
func (h *Handler) InstallArrowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	// Parse request body for variables and optional repository
	var requestBody struct {
		Variables  map[string]string `json:"variables"`
		Repository string            `json:"repository,omitempty"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		// Variables are optional, so ignore decode errors
		requestBody.Variables = make(map[string]string)
	}

	// Check for repository specification in query parameter
	if repo := r.URL.Query().Get("repository"); repo != "" {
		requestBody.Repository = repo
	}

	// Build full package specification
	fullName := name
	if requestBody.Repository != "" {
		fullName = requestBody.Repository + "@" + name
	}

	err := h.pkgManager.PackageManager.InstallArrow(fullName, requestBody.Variables)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Failed to install arrow: "+err.Error())
		return
	}

	responseData := map[string]string{
		"message":    "Arrow installed successfully",
		"arrow":      name,
		"repository": requestBody.Repository,
	}
	response.WriteJSON(w, http.StatusOK, responseData)
}

// ExecuteArrowHandler handles arrow execution
func (h *Handler) ExecuteArrowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	// Parse request body for variables
	var requestBody struct {
		Variables map[string]string `json:"variables"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		// Variables are optional, so ignore decode errors
		requestBody.Variables = make(map[string]string)
	}

	err := h.pkgManager.PackageManager.ExecuteArrow(name, requestBody.Variables)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Failed to execute arrow: "+err.Error())
		return
	}

	responseData := map[string]string{
		"message": "Arrow executed successfully",
		"arrow":   name,
	}
	response.WriteJSON(w, http.StatusOK, responseData)
}

// UninstallArrowHandler handles arrow uninstallation
func (h *Handler) UninstallArrowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	err := h.pkgManager.PackageManager.UninstallArrow(name)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Failed to uninstall arrow: "+err.Error())
		return
	}

	responseData := map[string]string{
		"message": "Arrow uninstalled successfully",
		"arrow":   name,
	}
	response.WriteJSON(w, http.StatusOK, responseData)
}

// UpdateArrowHandler handles arrow updates
func (h *Handler) UpdateArrowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	// Parse request body for optional repository
	var requestBody struct {
		Repository string `json:"repository,omitempty"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		// Repository is optional, so ignore decode errors
	}

	// Check for repository specification in query parameter
	if repo := r.URL.Query().Get("repository"); repo != "" {
		requestBody.Repository = repo
	}

	// Build full package specification
	fullName := name
	if requestBody.Repository != "" {
		fullName = requestBody.Repository + "@" + name
	}

	err := h.pkgManager.PackageManager.UpdateArrow(fullName)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Failed to update arrow: "+err.Error())
		return
	}

	responseData := map[string]string{
		"message":    "Arrow updated successfully",
		"arrow":      name,
		"repository": requestBody.Repository,
	}
	response.WriteJSON(w, http.StatusOK, responseData)
}

// ValidateArrowHandler handles arrow validation
func (h *Handler) ValidateArrowHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	err := h.pkgManager.PackageManager.ValidateArrow(name)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Failed to validate arrow: "+err.Error())
		return
	}

	responseData := map[string]string{
		"message": "Arrow validation completed successfully",
		"arrow":   name,
	}
	response.WriteJSON(w, http.StatusOK, responseData)
}

// GetInstalledArrowsHandler handles listing installed arrows
func (h *Handler) GetInstalledArrowsHandler(w http.ResponseWriter, r *http.Request) {
	installed := h.pkgManager.PackageManager.GetInstalledArrows()

	responseData := map[string]interface{}{
		"installed": installed,
		"count":     len(installed),
	}
	response.WriteJSON(w, http.StatusOK, responseData)
}

// GetArrowStatusHandler handles getting arrow status
func (h *Handler) GetArrowStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	status, err := h.pkgManager.PackageManager.GetArrowStatus(name)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, "Arrow not found: "+err.Error())
		return
	}

	responseData := map[string]interface{}{
		"arrow":  name,
		"status": status,
	}
	response.WriteJSON(w, http.StatusOK, responseData)
} 