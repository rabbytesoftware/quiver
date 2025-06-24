package packages

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rabbytesoftware/quiver/internal/logger"
	"github.com/rabbytesoftware/quiver/internal/packages/manifest"
	v0_1 "github.com/rabbytesoftware/quiver/internal/packages/manifest/v0.1"
	"gopkg.in/yaml.v3"
)

// ArrowInfo represents information about an available arrow
type ArrowInfo struct {
	Name       string `json:"name"`
	Repository string `json:"repository"`
	Path       string `json:"path"` // Local path or URL
}

// RepositoryManager handles repository operations
type RepositoryManager struct {
	repositories []string
	httpClient   *http.Client
	logger       *logger.Logger
}

// NewRepositoryManager creates a new repository manager
func NewRepositoryManager(repositories []string, logger *logger.Logger) *RepositoryManager {
	return &RepositoryManager{
		repositories: repositories,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger.WithService("repository"),
	}
}

// SearchArrows searches for arrows across all repositories
func (rm *RepositoryManager) SearchArrows(query string) ([]*ArrowInfo, error) {
	var allArrows []*ArrowInfo

	for _, repo := range rm.repositories {
		arrows, err := rm.searchInRepository(repo, query)
		if err != nil {
			rm.logger.Warn("Failed to search in repository %s: %v", repo, err)
			continue
		}
		allArrows = append(allArrows, arrows...)
	}

	return allArrows, nil
}

// searchInRepository searches for arrows in a specific repository
func (rm *RepositoryManager) searchInRepository(repo, query string) ([]*ArrowInfo, error) {
	if rm.isURL(repo) {
		return rm.searchRemoteRepository(repo, query)
	}
	return rm.searchLocalRepository(repo, query)
}

// searchLocalRepository searches for arrows in a local directory
func (rm *RepositoryManager) searchLocalRepository(dirPath, query string) ([]*ArrowInfo, error) {
	var arrows []*ArrowInfo

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue walking even if there's an error
		}

		// Check if it's a YAML file
		if info.IsDir() || (!strings.HasSuffix(path, ".yaml") && !strings.HasSuffix(path, ".yml")) {
			return nil
		}

		// Extract filename without extension for matching
		filename := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		
		// Apply search filter based on filename only
		if query == "" || rm.matchesFilename(filename, query) {
			arrowInfo := &ArrowInfo{
				Name:       filename, // Use filename as name for search results
				Repository: dirPath,
				Path:       path,
			}
			arrows = append(arrows, arrowInfo)
		}

		return nil
	})

	return arrows, err
}

// searchRemoteRepository searches for arrows in a remote repository
func (rm *RepositoryManager) searchRemoteRepository(repoURL, query string) ([]*ArrowInfo, error) {
	// For remote repositories, we would need to implement a discovery mechanism
	// This is a simplified implementation that assumes a known structure
	
	// TODO: Implement proper remote repository discovery
	// For now, we'll return an empty list
	rm.logger.Warn("Remote repository search not fully implemented for %s", repoURL)
	return []*ArrowInfo{}, nil
}

// GetArrow fetches an arrow by name from repositories
func (rm *RepositoryManager) GetArrow(name string) (manifest.ArrowInterface, string, error) {
	for _, repo := range rm.repositories {
		arrow, path, err := rm.getArrowFromRepository(repo, name)
		if err == nil {
			return arrow, path, nil
		}
		rm.logger.Debug("Arrow %s not found in repository %s: %v", name, repo, err)
	}

	return nil, "", fmt.Errorf("arrow %s not found in any repository", name)
}

// getArrowFromRepository gets an arrow from a specific repository
func (rm *RepositoryManager) getArrowFromRepository(repo, name string) (manifest.ArrowInterface, string, error) {
	if rm.isURL(repo) {
		return rm.getArrowFromRemoteRepository(repo, name)
	}
	return rm.getArrowFromLocalRepository(repo, name)
}

// getArrowFromLocalRepository gets an arrow from a local repository
func (rm *RepositoryManager) getArrowFromLocalRepository(dirPath, name string) (manifest.ArrowInterface, string, error) {
	// Try different file patterns
	patterns := []string{
		filepath.Join(dirPath, name+".yaml"),
		filepath.Join(dirPath, name+".yml"),
		filepath.Join(dirPath, name, "arrow.yaml"),
		filepath.Join(dirPath, name, "arrow.yml"),
	}

	for _, path := range patterns {
		if _, err := os.Stat(path); err == nil {
			arrow, err := rm.loadArrowFromFile(path)
			if err != nil {
				continue
			}
			// Use filename-based matching instead of metadata name
			// This ensures consistency with search behavior
			return arrow, path, nil
		}
	}

	return nil, "", fmt.Errorf("arrow %s not found in local repository %s", name, dirPath)
}

// getArrowFromRemoteRepository gets an arrow from a remote repository
func (rm *RepositoryManager) getArrowFromRemoteRepository(repoURL, name string) (manifest.ArrowInterface, string, error) {
	// Try different URL patterns
	patterns := []string{
		fmt.Sprintf("%s/%s.yaml", repoURL, name),
		fmt.Sprintf("%s/%s.yml", repoURL, name),
		fmt.Sprintf("%s/%s/arrow.yaml", repoURL, name),
		fmt.Sprintf("%s/%s/arrow.yml", repoURL, name),
	}

	for _, url := range patterns {
		resp, err := rm.httpClient.Get(url)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			continue
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			continue
		}

		arrow, err := rm.loadArrowFromData(data)
		if err != nil {
			continue
		}

		// Use URL-based matching instead of metadata name
		// This ensures consistency with search behavior
		return arrow, url, nil
	}

	return nil, "", fmt.Errorf("arrow %s not found in remote repository %s", name, repoURL)
}

// DownloadArrow downloads an arrow to local storage
func (rm *RepositoryManager) DownloadArrow(arrowInfo *ArrowInfo, targetPath string) error {
	if rm.isURL(arrowInfo.Path) {
		return rm.downloadRemoteArrow(arrowInfo.Path, targetPath)
	}
	return rm.copyLocalArrow(arrowInfo.Path, targetPath)
}

// downloadRemoteArrow downloads an arrow from a remote URL
func (rm *RepositoryManager) downloadRemoteArrow(url, targetPath string) error {
	// Ensure target directory exists
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Download the file
	resp, err := rm.httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download arrow: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download arrow: HTTP %d", resp.StatusCode)
	}

	// Create target file
	file, err := os.Create(targetPath)
	if err != nil {
		return fmt.Errorf("failed to create target file: %w", err)
	}
	defer file.Close()

	// Copy data
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write arrow data: %w", err)
	}

	return nil
}

// copyLocalArrow copies an arrow from local storage
func (rm *RepositoryManager) copyLocalArrow(sourcePath, targetPath string) error {
	// Ensure target directory exists
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Read source file
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read source arrow: %w", err)
	}

	// Write to target
	if err := os.WriteFile(targetPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write target arrow: %w", err)
	}

	return nil
}

// AddRepository adds a new repository to the list
func (rm *RepositoryManager) AddRepository(repo string) {
	// Check if repository already exists
	for _, existing := range rm.repositories {
		if existing == repo {
			return
		}
	}

	rm.repositories = append(rm.repositories, repo)
	rm.logger.Info("Added repository: %s", repo)
}

// RemoveRepository removes a repository from the list
func (rm *RepositoryManager) RemoveRepository(repo string) {
	for i, existing := range rm.repositories {
		if existing == repo {
			rm.repositories = append(rm.repositories[:i], rm.repositories[i+1:]...)
			rm.logger.Info("Removed repository: %s", repo)
			return
		}
	}
}

// GetRepositories returns the list of repositories
func (rm *RepositoryManager) GetRepositories() []string {
	return rm.repositories
}

// Helper methods

// isURL checks if a path is a URL
func (rm *RepositoryManager) isURL(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

// loadArrowFromFile loads an arrow from a file
func (rm *RepositoryManager) loadArrowFromFile(path string) (manifest.ArrowInterface, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return rm.loadArrowFromData(data)
}

// loadArrowFromData loads an arrow from data
func (rm *RepositoryManager) loadArrowFromData(data []byte) (manifest.ArrowInterface, error) {
	// Parse version first
	var versionInfo struct {
		Version string `yaml:"version"`
	}

	if err := yaml.Unmarshal(data, &versionInfo); err != nil {
		return nil, fmt.Errorf("failed to parse version: %w", err)
	}

	// Load based on version
	switch versionInfo.Version {
	case "0.1":
		var arrow v0_1.Arrow
		if err := yaml.Unmarshal(data, &arrow); err != nil {
			return nil, fmt.Errorf("failed to unmarshal v0.1 arrow: %w", err)
		}
		return &arrow, nil
	default:
		return nil, fmt.Errorf("unsupported version: %s", versionInfo.Version)
	}
}

// matchesFilename checks if a filename matches a search query
func (rm *RepositoryManager) matchesFilename(filename, query string) bool {
	if query == "" {
		return true
	}
	
	filename = strings.ToLower(filename)
	query = strings.ToLower(query)
	
	// Direct substring match
	if strings.Contains(filename, query) {
		return true
	}
	
	// Fuzzy matching - remove common separators
	fuzzyFilename := strings.ReplaceAll(strings.ReplaceAll(filename, "-", ""), "_", "")
	fuzzyQuery := strings.ReplaceAll(strings.ReplaceAll(query, "-", ""), "_", "")
	
	if strings.Contains(fuzzyFilename, fuzzyQuery) {
		return true
	}
	
	return false
} 