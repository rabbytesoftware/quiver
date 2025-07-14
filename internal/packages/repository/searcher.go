package repository

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/rabbytesoftware/quiver/internal/logger"
	"github.com/rabbytesoftware/quiver/internal/packages/manifest"
	"github.com/rabbytesoftware/quiver/internal/packages/types"
)

// Searcher handles search operations in repositories
type Searcher struct {
	httpClient *http.Client
	processor  *manifest.Processor
	logger     *logger.Logger
}

// NewSearcher creates a new searcher
func NewSearcher(httpClient *http.Client, logger *logger.Logger) *Searcher {
	return &Searcher{
		httpClient: httpClient,
		processor:  manifest.NewProcessor(logger),
		logger:     logger.WithService("repository-searcher"),
	}
}

// SearchInRepository searches for arrows in a specific repository
func (s *Searcher) SearchInRepository(repo, query string) ([]*types.ArrowInfo, error) {
	if s.isURL(repo) {
		return s.searchRemoteRepository(repo, query)
	}
	return s.searchLocalRepository(repo, query)
}

// searchLocalRepository searches for arrows in a local directory
func (s *Searcher) searchLocalRepository(dirPath, query string) ([]*types.ArrowInfo, error) {
	var arrows []*types.ArrowInfo

	// Try different file patterns for the query
	patterns := []string{
		filepath.Join(dirPath, query+".yaml"),
		filepath.Join(dirPath, query+".yml"),
		filepath.Join(dirPath, query, "arrow.yaml"),
		filepath.Join(dirPath, query, "arrow.yml"),
	}

	for _, path := range patterns {
		if _, err := os.Stat(path); err == nil {
			// File exists, verify it's a valid arrow
			if s.processor.IsValidArrowFile(path) {
				arrowInfo := &types.ArrowInfo{
					Name:           query,
					PackageName:    filepath.Base(path),
					RepositoryPath: dirPath,
					Path:           path,
				}
				arrows = append(arrows, arrowInfo)
				s.logger.Debug("Found arrow %s in local repository %s", query, dirPath)
				break // Found the file, no need to try other patterns
			}
		}
	}

	return arrows, nil
}

// searchRemoteRepository searches for arrows in a remote repository
func (s *Searcher) searchRemoteRepository(repoURL, query string) ([]*types.ArrowInfo, error) {
	// For remote repositories, try to directly fetch the specific yaml file
	patterns := []string{
		fmt.Sprintf("%s/%s.yaml", repoURL, query),
		fmt.Sprintf("%s/%s.yml", repoURL, query),
		fmt.Sprintf("%s/%s/arrow.yaml", repoURL, query),
		fmt.Sprintf("%s/%s/arrow.yml", repoURL, query),
	}

	var arrows []*types.ArrowInfo

	for _, url := range patterns {
		resp, err := s.httpClient.Get(url)
		if err != nil {
			s.logger.Debug("Failed to fetch %s: %v", url, err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			// File exists, verify it's a valid arrow
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				s.logger.Debug("Failed to read response from %s: %v", url, err)
				continue
			}

			// Try to parse it as a valid arrow manifest
			_, err = s.processor.LoadFromData(data)
			if err != nil {
				s.logger.Debug("Invalid arrow manifest at %s: %v", url, err)
				continue
			}

			arrowInfo := &types.ArrowInfo{
				Name:           query,
				PackageName:    filepath.Base(url),
				RepositoryPath: repoURL,
				Path:           url,
			}
			arrows = append(arrows, arrowInfo)
			s.logger.Debug("Found arrow %s in remote repository %s", query, repoURL)
			break // Found the file, no need to try other patterns
		}
	}

	return arrows, nil
}

// GetArrowFromRepository gets an arrow from a specific repository
func (s *Searcher) GetArrowFromRepository(repo, name string) (manifest.ArrowInterface, string, error) {
	if s.isURL(repo) {
		return s.getArrowFromRemoteRepository(repo, name)
	}
	return s.getArrowFromLocalRepository(repo, name)
}

// getArrowFromLocalRepository gets an arrow from a local repository
func (s *Searcher) getArrowFromLocalRepository(dirPath, name string) (manifest.ArrowInterface, string, error) {
	// Try different file patterns
	patterns := []string{
		filepath.Join(dirPath, name+".yaml"),
		filepath.Join(dirPath, name+".yml"),
		filepath.Join(dirPath, name, "arrow.yaml"),
		filepath.Join(dirPath, name, "arrow.yml"),
	}

	for _, path := range patterns {
		if _, err := os.Stat(path); err == nil {
			arrow, err := s.processor.LoadFromFile(path)
			if err != nil {
				continue
			}
			return arrow, path, nil
		}
	}

	return nil, "", fmt.Errorf("arrow %s not found in local repository %s", name, dirPath)
}

// getArrowFromRemoteRepository gets an arrow from a remote repository
func (s *Searcher) getArrowFromRemoteRepository(repoURL, name string) (manifest.ArrowInterface, string, error) {
	// Try different URL patterns
	patterns := []string{
		fmt.Sprintf("%s/%s.yaml", repoURL, name),
		fmt.Sprintf("%s/%s.yml", repoURL, name),
		fmt.Sprintf("%s/%s/arrow.yaml", repoURL, name),
		fmt.Sprintf("%s/%s/arrow.yml", repoURL, name),
	}

	for _, url := range patterns {
		resp, err := s.httpClient.Get(url)
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

		arrow, err := s.processor.LoadFromData(data)
		if err != nil {
			continue
		}

		return arrow, url, nil
	}

	return nil, "", fmt.Errorf("arrow %s not found in remote repository %s", name, repoURL)
}

// ValidateRepository checks if a repository is accessible
func (s *Searcher) ValidateRepository(repoPath string) error {
	if s.isURL(repoPath) {
		// For remote repositories, try a simple GET request
		resp, err := s.httpClient.Get(repoPath)
		if err != nil {
			return fmt.Errorf("failed to access remote repository: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			return fmt.Errorf("remote repository returned status %d", resp.StatusCode)
		}
		return nil
	}

	// For local repositories, check if directory exists
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return fmt.Errorf("local repository directory does not exist: %s", repoPath)
	}

	return nil
}

// isURL checks if a path is a URL
func (s *Searcher) isURL(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
} 