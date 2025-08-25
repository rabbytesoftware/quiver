package repository

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/rabbytesoftware/quiver/internal/logger"
	"github.com/rabbytesoftware/quiver/internal/packages/types"
)

// Downloader handles download operations for arrows
type Downloader struct {
	httpClient *http.Client
	logger     *logger.Logger
}

// NewDownloader creates a new downloader
func NewDownloader(httpClient *http.Client, logger *logger.Logger) *Downloader {
	return &Downloader{
		httpClient: httpClient,
		logger:     logger.WithService("repository-downloader"),
	}
}

// DownloadArrow downloads an arrow to local storage
func (d *Downloader) DownloadArrow(arrowInfo *types.ArrowInfo, targetPath string) error {
	if d.isURL(arrowInfo.Path) {
		return d.downloadRemoteArrow(arrowInfo.Path, targetPath)
	}
	return d.copyLocalArrow(arrowInfo.Path, targetPath)
}

// downloadRemoteArrow downloads an arrow from a remote URL
func (d *Downloader) downloadRemoteArrow(url, targetPath string) error {
	d.logger.Debug("Downloading arrow from URL: %s to %s", url, targetPath)

	// Ensure target directory exists
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}

	// Download the file
	resp, err := d.httpClient.Get(url)
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

	d.logger.Debug("Successfully downloaded arrow to: %s", targetPath)
	return nil
}

// copyLocalArrow copies an arrow from local storage
func (d *Downloader) copyLocalArrow(sourcePath, targetPath string) error {
	d.logger.Debug("Copying arrow from local path: %s to %s", sourcePath, targetPath)

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

	d.logger.Debug("Successfully copied arrow to: %s", targetPath)
	return nil
}

// DownloadArrowToDirectory downloads an arrow and creates the installation directory structure
func (d *Downloader) DownloadArrowToDirectory(arrowInfo *types.ArrowInfo, installDir string) error {
	// Create install directory
	if err := os.MkdirAll(installDir, 0755); err != nil {
		return fmt.Errorf("failed to create install directory: %w", err)
	}

	// Determine target file name
	targetPath := filepath.Join(installDir, "arrow.yaml")

	// Download the arrow
	return d.DownloadArrow(arrowInfo, targetPath)
}

// ValidateDownload verifies that a downloaded arrow is valid
func (d *Downloader) ValidateDownload(targetPath string) error {
	// Check if file exists
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		return fmt.Errorf("downloaded file does not exist: %s", targetPath)
	}

	// Check if file is not empty
	info, err := os.Stat(targetPath)
	if err != nil {
		return fmt.Errorf("failed to stat downloaded file: %w", err)
	}

	if info.Size() == 0 {
		return fmt.Errorf("downloaded file is empty: %s", targetPath)
	}

	d.logger.Debug("Download validation passed for: %s", targetPath)
	return nil
}

// isURL checks if a path is a URL
func (d *Downloader) isURL(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}

// GetDownloadSize returns the size of a download (for remote URLs only)
func (d *Downloader) GetDownloadSize(url string) (int64, error) {
	if !d.isURL(url) {
		// For local files, get the file size directly
		info, err := os.Stat(url)
		if err != nil {
			return 0, fmt.Errorf("failed to get local file size: %w", err)
		}
		return info.Size(), nil
	}

	// For remote URLs, make a HEAD request
	resp, err := d.httpClient.Head(url)
	if err != nil {
		return 0, fmt.Errorf("failed to get download size: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to get download size: HTTP %d", resp.StatusCode)
	}

	return resp.ContentLength, nil
} 