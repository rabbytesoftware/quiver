package repository

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/rabbytesoftware/quiver/internal/logger"
	"github.com/rabbytesoftware/quiver/internal/packages/manifest"
	"github.com/rabbytesoftware/quiver/internal/packages/types"
)

// Manager handles repository operations and coordination
type Manager struct {
	repositories []string
	httpClient   *http.Client
	searcher     *Searcher
	downloader   *Downloader
	logger       *logger.Logger
}

// NewManager creates a new repository manager
func NewManager(repositories []string, logger *logger.Logger) *Manager {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	manager := &Manager{
		repositories: repositories,
		httpClient:   httpClient,
		logger:       logger.WithService("repository"),
	}

	// Initialize sub-components
	manager.searcher = NewSearcher(httpClient, logger)
	manager.downloader = NewDownloader(httpClient, logger)

	return manager
}

// Repository operations

// AddRepository adds a new repository to the list
func (m *Manager) AddRepository(repo string) {
	// Check if repository already exists
	for _, existing := range m.repositories {
		if existing == repo {
			return
		}
	}

	m.repositories = append(m.repositories, repo)
	m.logger.Info("Added repository: %s", repo)
}

// RemoveRepository removes a repository from the list
func (m *Manager) RemoveRepository(repo string) {
	for i, existing := range m.repositories {
		if existing == repo {
			m.repositories = append(m.repositories[:i], m.repositories[i+1:]...)
			m.logger.Info("Removed repository: %s", repo)
			return
		}
	}
}

// GetRepositories returns the list of repositories
func (m *Manager) GetRepositories() []string {
	return m.repositories
}

// SetRepositories replaces the current repository list
func (m *Manager) SetRepositories(repositories []string) {
	m.repositories = repositories
	m.logger.Info("Updated repositories list with %d repositories", len(repositories))
}

// Search operations

// SearchArrows searches for arrows across all repositories
func (m *Manager) SearchArrows(query string) ([]*types.ArrowInfo, error) {
	// Check if query contains repository specification (repo@package)
	if strings.Contains(query, "@") {
		return m.searchWithRepositorySpec(query)
	}

	// Search in all repositories
	var allArrows []*types.ArrowInfo
	for _, repo := range m.repositories {
		arrows, err := m.searcher.SearchInRepository(repo, query)
		if err != nil {
			m.logger.Warn("Failed to search in repository %s: %v", repo, err)
			continue
		}
		allArrows = append(allArrows, arrows...)
	}

	return allArrows, nil
}

// searchWithRepositorySpec handles repository-specific search
func (m *Manager) searchWithRepositorySpec(query string) ([]*types.ArrowInfo, error) {
	repoPath, packageName, _ := m.ParseRepositorySpec(query)
	
	// Find the specified repository
	targetRepo := m.findRepository(repoPath)
	if targetRepo == "" {
		return nil, fmt.Errorf("repository %s not found", repoPath)
	}
	
	// Search only in the specified repository
	arrows, err := m.searcher.SearchInRepository(targetRepo, packageName)
	if err != nil {
		return nil, fmt.Errorf("failed to search in repository %s: %v", targetRepo, err)
	}
	return arrows, nil
}

// GetArrow fetches an arrow by name from repositories
func (m *Manager) GetArrow(name string) (manifest.ArrowInterface, string, error) {
	// Check if name contains repository specification (repo@package)
	if strings.Contains(name, "@") {
		return m.getArrowWithRepositorySpec(name)
	}

	// Search all repositories
	for _, repo := range m.repositories {
		arrow, path, err := m.searcher.GetArrowFromRepository(repo, name)
		if err == nil {
			return arrow, path, nil
		}
		m.logger.Debug("Arrow %s not found in repository %s: %v", name, repo, err)
	}

	return nil, "", fmt.Errorf("arrow %s not found in any repository", name)
}

// getArrowWithRepositorySpec handles repository-specific arrow retrieval
func (m *Manager) getArrowWithRepositorySpec(name string) (manifest.ArrowInterface, string, error) {
	repoPath, packageName, _ := m.ParseRepositorySpec(name)
	
	// Find the specified repository
	targetRepo := m.findRepository(repoPath)
	if targetRepo == "" {
		return nil, "", fmt.Errorf("repository %s not found", repoPath)
	}
	
	// Get arrow from the specified repository
	arrow, path, err := m.searcher.GetArrowFromRepository(targetRepo, packageName)
	if err != nil {
		return nil, "", fmt.Errorf("arrow %s not found in repository %s: %v", packageName, targetRepo, err)
	}
	return arrow, path, nil
}

// Download operations

// DownloadArrow downloads an arrow to local storage
func (m *Manager) DownloadArrow(arrowInfo *types.ArrowInfo, targetPath string) error {
	return m.downloader.DownloadArrow(arrowInfo, targetPath)
}

// Utility methods

// ParseRepositorySpec parses repository specification syntax
func (m *Manager) ParseRepositorySpec(spec string) (repositoryPath, packageName string, hasRepoSpec bool) {
	if strings.Contains(spec, "@") {
		parts := strings.SplitN(spec, "@", 2)
		if len(parts) == 2 {
			return parts[0], parts[1], true
		}
	}
	return "", spec, false
}

// findRepository finds a repository by path or partial match
func (m *Manager) findRepository(repoPath string) string {
	for _, repo := range m.repositories {
		if strings.Contains(repo, repoPath) || repo == repoPath {
			return repo
		}
	}
	return ""
}

// isURL checks if a path is a URL
func (m *Manager) isURL(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

// ValidateRepository checks if a repository is accessible
func (m *Manager) ValidateRepository(repoPath string) error {
	return m.searcher.ValidateRepository(repoPath)
} 