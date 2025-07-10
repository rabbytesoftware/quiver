package server

import (
	"fmt"

	"github.com/rabbytesoftware/quiver/internal/logger"
	"github.com/rabbytesoftware/quiver/internal/packages/core"
	"github.com/rabbytesoftware/quiver/internal/packages/types"
)

// ArrowsServer provides the main server interface for package operations
type ArrowsServer struct {
	packageManager *core.Manager
	logger         *logger.Logger
}

// NewArrowsServer creates a new arrows server
func NewArrowsServer(repositories []string, installDir, dbPath string, logger *logger.Logger) *ArrowsServer {
	pm := core.NewManager(repositories, installDir, dbPath, logger)
	
	return &ArrowsServer{
		packageManager: pm,
		logger:         logger.WithService("arrows-server"),
	}
}

// Initialize initializes the arrows server
func (as *ArrowsServer) Initialize() error {
	if err := as.packageManager.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize package manager: %w", err)
	}

	as.logger.Info("Arrows server initialized successfully")
	return nil
}

// Package Management Operations

// SearchArrows searches for arrows across repositories
func (as *ArrowsServer) SearchArrows(query string) ([]*types.ArrowInfo, error) {
	return as.packageManager.SearchArrows(query)
}

// InstallArrow installs an arrow with dependency resolution
func (as *ArrowsServer) InstallArrow(name string, variables map[string]string) error {
	return as.packageManager.InstallArrow(name, variables)
}

// UninstallArrow uninstalls an arrow and handles dependency cleanup
func (as *ArrowsServer) UninstallArrow(name string) error {
	return as.packageManager.UninstallArrow(name)
}

// ExecuteArrow executes an arrow's execute method
func (as *ArrowsServer) ExecuteArrow(name string, variables map[string]string) error {
	return as.packageManager.ExecuteArrow(name, variables)
}

// UpdateArrow updates an arrow to the latest version
func (as *ArrowsServer) UpdateArrow(name string) error {
	return as.packageManager.UpdateArrow(name)
}

// ValidateArrow validates an arrow installation
func (as *ArrowsServer) ValidateArrow(name string) error {
	return as.packageManager.ValidateArrow(name)
}

// Information Operations

// GetInstalledArrows returns all installed arrows
func (as *ArrowsServer) GetInstalledArrows() map[string]*types.InstalledPackage {
	return as.packageManager.GetInstalledArrows()
}

// GetArrowStatus returns the status of an arrow
func (as *ArrowsServer) GetArrowStatus(name string) (types.PackageStatus, error) {
	return as.packageManager.GetArrowStatus(name)
}

// GetArrowsByStatus returns arrows with a specific status
func (as *ArrowsServer) GetArrowsByStatus(status types.PackageStatus) []*types.InstalledPackage {
	installed := as.packageManager.GetInstalledArrows()
	var filtered []*types.InstalledPackage
	
	for _, pkg := range installed {
		if types.PackageStatus(pkg.Status) == status {
			filtered = append(filtered, pkg)
		}
	}
	
	return filtered
}

// Repository Management

// AddRepository adds a new repository
func (as *ArrowsServer) AddRepository(repo string) {
	as.packageManager.AddRepository(repo)
	as.logger.Info("Added repository: %s", repo)
}

// RemoveRepository removes a repository
func (as *ArrowsServer) RemoveRepository(repo string) {
	as.packageManager.RemoveRepository(repo)
	as.logger.Info("Removed repository: %s", repo)
}

// GetRepositories returns all repositories
func (as *ArrowsServer) GetRepositories() []string {
	return as.packageManager.GetRepositories()
}

// Health and Status

// GetServerStatus returns the server status information
func (as *ArrowsServer) GetServerStatus() *ServerStatus {
	installed := as.packageManager.GetInstalledArrows()
	repositories := as.packageManager.GetRepositories()
	
	var runningCount, stoppedCount, errorCount int
	for _, pkg := range installed {
		switch types.PackageStatus(pkg.Status) {
		case types.StatusRunning:
			runningCount++
		case types.StatusStopped:
			stoppedCount++
		case types.StatusError:
			errorCount++
		}
	}
	
	return &ServerStatus{
		TotalPackages:     len(installed),
		RunningPackages:   runningCount,
		StoppedPackages:   stoppedCount,
		ErrorPackages:     errorCount,
		TotalRepositories: len(repositories),
		Repositories:      repositories,
	}
}

// ServerStatus contains server status information
type ServerStatus struct {
	TotalPackages     int      `json:"total_packages"`
	RunningPackages   int      `json:"running_packages"`
	StoppedPackages   int      `json:"stopped_packages"`
	ErrorPackages     int      `json:"error_packages"`
	TotalRepositories int      `json:"total_repositories"`
	Repositories      []string `json:"repositories"`
}

// Shutdown gracefully shuts down the server
func (as *ArrowsServer) Shutdown() error {
	as.logger.Info("Shutting down arrows server")
	// Perform any cleanup operations here
	return nil
} 