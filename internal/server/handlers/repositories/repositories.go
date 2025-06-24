package repositories

import (
	"encoding/json"
	"net/http"

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

// AddRepositoryHandler handles adding a new repository
func (h *Handler) AddRepositoryHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Repository string `json:"repository"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if requestBody.Repository == "" {
		response.WriteError(w, http.StatusBadRequest, "Repository URL/path is required")
		return
	}

	h.pkgManager.PackageManager.AddRepository(requestBody.Repository)

	responseData := map[string]string{
		"message":    "Repository added successfully",
		"repository": requestBody.Repository,
	}
	response.WriteJSON(w, http.StatusOK, responseData)
}

// RemoveRepositoryHandler handles removing a repository
func (h *Handler) RemoveRepositoryHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Repository string `json:"repository"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if requestBody.Repository == "" {
		response.WriteError(w, http.StatusBadRequest, "Repository URL/path is required")
		return
	}

	h.pkgManager.PackageManager.RemoveRepository(requestBody.Repository)

	responseData := map[string]string{
		"message":    "Repository removed successfully",
		"repository": requestBody.Repository,
	}
	response.WriteJSON(w, http.StatusOK, responseData)
}

// GetRepositoriesHandler handles listing all repositories
func (h *Handler) GetRepositoriesHandler(w http.ResponseWriter, r *http.Request) {
	repositories := h.pkgManager.PackageManager.GetRepositories()

	responseData := map[string]interface{}{
		"repositories": repositories,
		"count":        len(repositories),
	}
	response.WriteJSON(w, http.StatusOK, responseData)
} 