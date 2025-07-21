package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rabbytesoftware/quiver/internal/database"
	"github.com/rabbytesoftware/quiver/internal/logger"
	"github.com/rabbytesoftware/quiver/internal/packages/execution"
	"github.com/rabbytesoftware/quiver/internal/packages/execution/process"
	"github.com/rabbytesoftware/quiver/internal/packages/manifest"
	"github.com/rabbytesoftware/quiver/internal/packages/repository"
	"github.com/rabbytesoftware/quiver/internal/packages/types"
)

// ExecutionStatus represents the status of an async execution
type ExecutionStatus struct {
	ArrowName   string            `json:"arrow_name"`
	Status      string            `json:"status"`      // "running", "completed", "failed"
	StartedAt   time.Time         `json:"started_at"`
	CompletedAt *time.Time        `json:"completed_at,omitempty"`
	Error       string            `json:"error,omitempty"`
	Variables   map[string]string `json:"variables"`
}

// Manager is the main package manager that orchestrates all package operations
type Manager struct {
	database      database.Database
	repository    *repository.Manager
	execution     *execution.Engine
	arrowProcessor *manifest.Processor
	installDir    string
	logger        *logger.Logger
	executions     map[string]*ExecutionStatus
	executionMutex sync.RWMutex
}

// NewManager creates a new package manager
func NewManager(repositories []string, installDir, dbPath string, logger *logger.Logger) *Manager {
	return &Manager{
		database:       database.NewDefaultDatabase(dbPath, logger),
		repository:     repository.NewManager(repositories, logger),
		execution:      execution.NewEngine(logger),
		arrowProcessor: manifest.NewProcessor(logger),
		installDir:     installDir,
		logger:         logger.WithService("package-manager"),
		executions:     make(map[string]*ExecutionStatus),
	}
}

// Initialize initializes the package manager
func (m *Manager) Initialize() error {
	// Ensure install directory exists
	if err := os.MkdirAll(m.installDir, 0755); err != nil {
		return fmt.Errorf("failed to create install directory: %w", err)
	}

	// Load database
	if err := m.database.Load(); err != nil {
		return fmt.Errorf("failed to load database: %w", err)
	}

	m.logger.Info("Package manager initialized")

	// Load manifests
	registerManifests()

	return nil
}

// Package Operations

// SearchArrows searches for arrows across repositories
func (m *Manager) SearchArrows(query string) ([]*types.ArrowInfo, error) {
	return m.repository.SearchArrows(query)
}

// InstallArrow installs an arrow with dependency resolution
func (m *Manager) InstallArrow(name string, variables map[string]string) error {
	m.logger.Info("Installing arrow: %s", name)

	// Parse repository specification if present
	repoPath, packageName, hasRepoSpec := m.repository.ParseRepositorySpec(name)
	displayName := name
	if hasRepoSpec {
		m.logger.Info("Installing arrow %s from repository %s", packageName, repoPath)
		displayName = packageName
	}

	// Check if already installed
	if m.database.IsInstalled(displayName) {
		return fmt.Errorf("arrow %s is already installed", displayName)
	}

	// Get arrow from repository
	arrow, sourcePath, err := m.repository.GetArrow(name)
	if err != nil {
		return fmt.Errorf("failed to get arrow %s: %w", name, err)
	}

	// Validate arrow
	if err := m.arrowProcessor.ValidateArrow(arrow); err != nil {
		return fmt.Errorf("invalid arrow %s: %w", name, err)
	}

	// Install dependencies first
	dependencies := arrow.GetDependencies()
	for _, dep := range dependencies {
		if !m.database.IsInstalled(dep) {
			m.logger.Info("Installing dependency: %s", dep)
			if err := m.InstallArrow(dep, nil); err != nil {
				return fmt.Errorf("failed to install dependency %s: %w", dep, err)
			}
		}
		// Add dependency relationship
		if err := m.database.AddDependency(displayName, dep); err != nil {
			m.logger.Warn("Failed to add dependency relationship: %v", err)
		}
	}

	// Create install path
	arrowInstallPath := filepath.Join(m.installDir, displayName)
	arrowFile := filepath.Join(arrowInstallPath, "arrow.yaml")

	// Download/copy arrow
	arrowInfo := &types.ArrowInfo{
		Name: displayName,
		Path: sourcePath,
	}
	if err := m.repository.DownloadArrow(arrowInfo, arrowFile); err != nil {
		return fmt.Errorf("failed to download arrow: %w", err)
	}

	// Execute install method
	ctx := &types.ExecutionContext{
		ArrowName:   displayName,
		InstallPath: arrowInstallPath,
		Variables:   variables,
	}
	if err := m.execution.ExecuteMethod(arrow, types.MethodInstall, ctx); err != nil {
		// Cleanup on failure
		os.RemoveAll(arrowInstallPath)
		return fmt.Errorf("failed to execute install method: %w", err)
	}

	// Add to database
	installedPkg := &types.InstalledPackage{
		Name:         displayName,
		Version:      arrow.ArrowVersion(),
		Repository:   sourcePath,
		InstallPath:  arrowInstallPath,
		Dependencies: dependencies,
		Variables:    variables,
		Status:       string(types.StatusInstalled),
		InstalledAt:  time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := m.database.AddPackage(installedPkg); err != nil {
		return fmt.Errorf("failed to add package to database: %w", err)
	}

	m.logger.Info("Successfully installed arrow: %s", displayName)
	return nil
}

// UninstallArrow uninstalls an arrow and handles dependency cleanup
func (m *Manager) UninstallArrow(name string) error {
	m.logger.Info("Uninstalling arrow: %s", name)

	// Check if installed
	pkg, exists := m.database.GetPackage(name)
	if !exists {
		return fmt.Errorf("arrow %s is not installed", name)
	}

	// Check if other packages depend on this one
	if m.database.HasDependents(name) {
		dependents := m.database.GetDependents(name)
		return fmt.Errorf("cannot uninstall %s: packages %v depend on it", name, dependents)
	}

	// Load arrow for uninstall execution
	arrow, err := m.arrowProcessor.LoadFromInstallation(pkg.InstallPath)
	if err != nil {
		m.logger.Warn("Failed to load arrow for uninstall, proceeding anyway: %v", err)
	} else {
		// Execute uninstall method
		ctx := &types.ExecutionContext{
			ArrowName:   name,
			InstallPath: pkg.InstallPath,
			Variables:   pkg.Variables,
		}
		if err := m.execution.ExecuteMethod(arrow, types.MethodUninstall, ctx); err != nil {
			m.logger.Warn("Failed to execute uninstall method: %v", err)
		}
	}

	// Remove installation directory
	if err := os.RemoveAll(pkg.InstallPath); err != nil {
		m.logger.Warn("Failed to remove install directory: %v", err)
	}

	// Remove from database
	if err := m.database.RemovePackage(name); err != nil {
		return fmt.Errorf("failed to remove package from database: %w", err)
	}

	// Clean up unused dependencies
	for _, dep := range pkg.Dependencies {
		if !m.database.HasDependents(dep) {
			m.logger.Info("Dependency %s is no longer needed, consider removing it", dep)
		}
	}

	m.logger.Info("Successfully uninstalled arrow: %s", name)
	return nil
}

// ExecuteArrow executes an arrow's execute method
func (m *Manager) ExecuteArrow(name string, variables map[string]string) error {
	m.logger.Info("Executing arrow: %s", name)

	// Check if installed
	pkg, exists := m.database.GetPackage(name)
	if !exists {
		return fmt.Errorf("arrow %s is not installed", name)
	}

	// Merge variables (user provided override defaults)
	finalVariables := make(map[string]string)
	for k, v := range pkg.Variables {
		finalVariables[k] = v
	}
	for k, v := range variables {
		finalVariables[k] = v
	}

	// Load arrow
	arrow, err := m.arrowProcessor.LoadFromInstallation(pkg.InstallPath)
	if err != nil {
		return fmt.Errorf("failed to load arrow: %w", err)
	}

	// Track execution start
	m.trackExecutionStart(name, finalVariables)

	// Update status to running
	if err := m.database.UpdatePackageStatus(name, types.StatusRunning); err != nil {
		m.logger.Warn("Failed to update package status: %v", err)
	}

	// Execute method
	ctx := &types.ExecutionContext{
		ArrowName:   name,
		InstallPath: pkg.InstallPath,
		Variables:   finalVariables,
	}
	err = m.execution.ExecuteMethod(arrow, types.MethodExecute, ctx)
	
	// Update status based on result
	status := types.StatusStopped
	if err != nil {
		status = types.StatusError
		m.trackExecutionComplete(name, err)
	} else {
		m.trackExecutionComplete(name, nil)
	}
	
	if statusErr := m.database.UpdatePackageStatus(name, status); statusErr != nil {
		m.logger.Warn("Failed to update package status: %v", statusErr)
	}

	return err
}

// ExecuteArrowAsync executes an arrow asynchronously and returns immediately
func (m *Manager) ExecuteArrowAsync(name string, variables map[string]string) error {
	// Validate that arrow exists before starting async execution
	_, exists := m.database.GetPackage(name)
	if !exists {
		return fmt.Errorf("arrow %s is not installed", name)
	}

	// Check if already running
	if m.IsExecutionRunning(name) {
		return fmt.Errorf("arrow %s is already being executed", name)
	}

	// Start async execution
	go func() {
		err := m.ExecuteArrow(name, variables)
		if err != nil {
			m.logger.Error("Async execution of arrow %s failed: %v", name, err)
		} else {
			m.logger.Info("Async execution of arrow %s completed successfully", name)
		}
	}()

	return nil
}

// Execution tracking methods

// trackExecutionStart records the start of an execution
func (m *Manager) trackExecutionStart(arrowName string, variables map[string]string) {
	m.executionMutex.Lock()
	defer m.executionMutex.Unlock()
	
	m.executions[arrowName] = &ExecutionStatus{
		ArrowName: arrowName,
		Status:    "running",
		StartedAt: time.Now(),
		Variables: variables,
	}
	
	m.logger.Debug("Started tracking execution for arrow: %s", arrowName)
}

// trackExecutionComplete records the completion of an execution
func (m *Manager) trackExecutionComplete(arrowName string, err error) {
	m.executionMutex.Lock()
	defer m.executionMutex.Unlock()
	
	if execution, exists := m.executions[arrowName]; exists {
		now := time.Now()
		execution.CompletedAt = &now
		
		if err != nil {
			execution.Status = "failed"
			execution.Error = err.Error()
		} else {
			execution.Status = "completed"
		}
		
		m.logger.Debug("Completed tracking execution for arrow: %s (status: %s)", arrowName, execution.Status)
	}
}

// GetExecutionStatus returns the current execution status for an arrow
func (m *Manager) GetExecutionStatus(arrowName string) (*ExecutionStatus, bool) {
	m.executionMutex.RLock()
	defer m.executionMutex.RUnlock()
	
	execution, exists := m.executions[arrowName]
	return execution, exists
}

// IsExecutionRunning checks if an arrow is currently being executed
func (m *Manager) IsExecutionRunning(arrowName string) bool {
	m.executionMutex.RLock()
	defer m.executionMutex.RUnlock()
	
	execution, exists := m.executions[arrowName]
	return exists && execution.Status == "running"
}

// GetAllExecutions returns all tracked executions
func (m *Manager) GetAllExecutions() map[string]*ExecutionStatus {
	m.executionMutex.RLock()
	defer m.executionMutex.RUnlock()
	
	result := make(map[string]*ExecutionStatus)
	for k, v := range m.executions {
		result[k] = v
	}
	return result
}

// CleanupOldExecutions removes execution records older than specified duration
func (m *Manager) CleanupOldExecutions(maxAge time.Duration) {
	m.executionMutex.Lock()
	defer m.executionMutex.Unlock()
	
	cutoff := time.Now().Add(-maxAge)
	for arrowName, execution := range m.executions {
		if execution.CompletedAt != nil && execution.CompletedAt.Before(cutoff) {
			delete(m.executions, arrowName)
			m.logger.Debug("Cleaned up old execution record for arrow: %s", arrowName)
		}
	}
}

// UpdateArrow updates an arrow to the latest version
func (m *Manager) UpdateArrow(name string) error {
	m.logger.Info("Updating arrow: %s", name)

	// Parse repository specification if present
	repoPath, packageName, hasRepoSpec := m.repository.ParseRepositorySpec(name)
	displayName := name
	if hasRepoSpec {
		m.logger.Info("Updating arrow %s from repository %s", packageName, repoPath)
		displayName = packageName
	}

	// Check if installed
	pkg, exists := m.database.GetPackage(displayName)
	if !exists {
		return fmt.Errorf("arrow %s is not installed", displayName)
	}

	// Get latest version from repository
	arrow, sourcePath, err := m.repository.GetArrow(name)
	if err != nil {
		return fmt.Errorf("failed to get arrow %s: %w", name, err)
	}

	// Check if update is needed
	if arrow.ArrowVersion() == pkg.Version {
		m.logger.Info("Arrow %s is already up to date", displayName)
		return nil
	}

	// Execute update method if available
	ctx := &types.ExecutionContext{
		ArrowName:   displayName,
		InstallPath: pkg.InstallPath,
		Variables:   pkg.Variables,
	}
	if err := m.execution.ExecuteMethod(arrow, types.MethodUpdate, ctx); err != nil {
		return fmt.Errorf("failed to execute update method: %w", err)
	}

	// Update arrow file
	arrowFile := filepath.Join(pkg.InstallPath, "arrow.yaml")
	arrowInfo := &types.ArrowInfo{
		Name: displayName,
		Path: sourcePath,
	}
	if err := m.repository.DownloadArrow(arrowInfo, arrowFile); err != nil {
		return fmt.Errorf("failed to download updated arrow: %w", err)
	}

	// Update database
	pkg.Version = arrow.ArrowVersion()
	pkg.Repository = sourcePath
	pkg.UpdatedAt = time.Now()
	if err := m.database.UpdatePackage(pkg); err != nil {
		return fmt.Errorf("failed to update package in database: %w", err)
	}

	m.logger.Info("Successfully updated arrow: %s to version %s", displayName, arrow.ArrowVersion())
	return nil
}

// StopArrow stops all running processes for an arrow
func (m *Manager) StopArrow(name string, graceful bool, timeout time.Duration) error {
	m.logger.Info("Stopping arrow: %s (graceful: %t, timeout: %v)", name, graceful, timeout)

	// Check if installed
	_, exists := m.database.GetPackage(name)
	if !exists {
		return fmt.Errorf("arrow %s is not installed", name)
	}

	// Stop all processes for this arrow
	err := m.execution.StopArrowProcesses(name, graceful, timeout)
	if err != nil {
		m.logger.Error("Failed to stop processes for arrow %s: %v", name, err)
		// Update status to error since we couldn't stop cleanly
		if statusErr := m.database.UpdatePackageStatus(name, types.StatusError); statusErr != nil {
			m.logger.Warn("Failed to update package status to error: %v", statusErr)
		}
		return fmt.Errorf("failed to stop arrow processes: %w", err)
	}

	// Update status to stopped
	if err := m.database.UpdatePackageStatus(name, types.StatusStopped); err != nil {
		m.logger.Warn("Failed to update package status to stopped: %v", err)
	}

	m.logger.Info("Successfully stopped arrow: %s", name)
	return nil
}

// GetArrowProcesses returns all running processes for an arrow
func (m *Manager) GetArrowProcesses(name string) ([]*process.ProcessInfo, error) {
	// Check if installed
	_, exists := m.database.GetPackage(name)
	if !exists {
		return nil, fmt.Errorf("arrow %s is not installed", name)
	}

	return m.execution.GetArrowProcesses(name), nil
}

// HasRunningProcesses checks if an arrow has any running processes
func (m *Manager) HasRunningProcesses(name string) bool {
	return m.execution.HasRunningProcesses(name)
}

// ValidateArrow validates an arrow installation
func (m *Manager) ValidateArrow(name string) error {
	m.logger.Info("Validating arrow: %s", name)

	// Check if installed
	pkg, exists := m.database.GetPackage(name)
	if !exists {
		return fmt.Errorf("arrow %s is not installed", name)
	}

	// Load arrow
	arrow, err := m.arrowProcessor.LoadFromInstallation(pkg.InstallPath)
	if err != nil {
		return fmt.Errorf("failed to load arrow: %w", err)
	}

	// Execute validate method
	ctx := &types.ExecutionContext{
		ArrowName:   name,
		InstallPath: pkg.InstallPath,
		Variables:   pkg.Variables,
	}
	return m.execution.ExecuteMethod(arrow, types.MethodValidate, ctx)
}

// Information Methods

// GetInstalledArrows returns all installed arrows
func (m *Manager) GetInstalledArrows() map[string]*types.InstalledPackage {
	return m.database.GetAllPackages()
}

// GetArrowStatus returns the status of an arrow
func (m *Manager) GetArrowStatus(name string) (types.PackageStatus, error) {
	return m.database.GetPackageStatus(name)
}

// GetArrowsByStatus returns arrows with a specific status
func (m *Manager) GetArrowsByStatus(status types.PackageStatus) []*types.InstalledPackage {
	installed := m.database.GetAllPackages()
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
func (m *Manager) AddRepository(repo string) {
	m.repository.AddRepository(repo)
}

// RemoveRepository removes a repository
func (m *Manager) RemoveRepository(repo string) {
	m.repository.RemoveRepository(repo)
}

// GetRepositories returns all repositories
func (m *Manager) GetRepositories() []string {
	return m.repository.GetRepositories()
} 