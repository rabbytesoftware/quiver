package packages

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/rabbytesoftware/quiver/internal/logger"
	"github.com/rabbytesoftware/quiver/internal/packages/manifest"
	v0_1 "github.com/rabbytesoftware/quiver/internal/packages/manifest/v0.1"
	"gopkg.in/yaml.v3"
)

// PackageManager handles all package management operations
type PackageManager struct {
	database   *DatabaseManager
	repository *RepositoryManager
	installDir string
	logger     *logger.Logger
}

// NewPackageManager creates a new package manager
func NewPackageManager(repositories []string, installDir, dbPath string, logger *logger.Logger) *PackageManager {
	return &PackageManager{
		database:   NewDatabaseManager(dbPath, logger),
		repository: NewRepositoryManager(repositories, logger),
		installDir: installDir,
		logger:     logger.WithService("package-manager"),
	}
}

// Initialize initializes the package manager
func (pm *PackageManager) Initialize() error {
	// Ensure install directory exists
	if err := os.MkdirAll(pm.installDir, 0755); err != nil {
		return fmt.Errorf("failed to create install directory: %w", err)
	}

	// Load database
	if err := pm.database.Load(); err != nil {
		return fmt.Errorf("failed to load database: %w", err)
	}

	pm.logger.Info("Package manager initialized")
	return nil
}

// SearchArrows searches for arrows across repositories
func (pm *PackageManager) SearchArrows(query string) ([]*ArrowInfo, error) {
	return pm.repository.SearchArrows(query)
}

// InstallArrow installs an arrow with dependency resolution
func (pm *PackageManager) InstallArrow(name string, variables map[string]string) error {
	pm.logger.Info("Installing arrow: %s", name)

	// Check if already installed
	if pm.database.IsInstalled(name) {
		return fmt.Errorf("arrow %s is already installed", name)
	}

	// Get arrow from repository
	arrow, sourcePath, err := pm.repository.GetArrow(name)
	if err != nil {
		return fmt.Errorf("failed to get arrow %s: %w", name, err)
	}

	// Install dependencies first
	dependencies := arrow.GetDependencies()
	for _, dep := range dependencies {
		if !pm.database.IsInstalled(dep) {
			pm.logger.Info("Installing dependency: %s", dep)
			if err := pm.InstallArrow(dep, nil); err != nil {
				return fmt.Errorf("failed to install dependency %s: %w", dep, err)
			}
		}
		// Add dependency relationship
		if err := pm.database.AddDependency(name, dep); err != nil {
			pm.logger.Warn("Failed to add dependency relationship: %v", err)
		}
	}

	// Create install path
	arrowInstallPath := filepath.Join(pm.installDir, name)
	arrowFile := filepath.Join(arrowInstallPath, "arrow.yaml")

	// Download/copy arrow
	arrowInfo := &ArrowInfo{
		Name: name,
		Path: sourcePath,
	}
	if err := pm.repository.DownloadArrow(arrowInfo, arrowFile); err != nil {
		return fmt.Errorf("failed to download arrow: %w", err)
	}

	// Execute install method
	if err := pm.executeMethod(arrow, "install", arrowInstallPath, variables); err != nil {
		// Cleanup on failure
		os.RemoveAll(arrowInstallPath)
		return fmt.Errorf("failed to execute install method: %w", err)
	}

	// Add to database
	installedPkg := &InstalledPackage{
		Name:         name,
		Version:      arrow.ArrowVersion(),
		Repository:   sourcePath,
		InstallPath:  arrowInstallPath,
		Dependencies: dependencies,
		Variables:    variables,
		Status:       "installed",
		InstalledAt:  time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := pm.database.AddPackage(installedPkg); err != nil {
		return fmt.Errorf("failed to add package to database: %w", err)
	}

	pm.logger.Info("Successfully installed arrow: %s", name)
	return nil
}

// UninstallArrow uninstalls an arrow and handles dependency cleanup
func (pm *PackageManager) UninstallArrow(name string) error {
	pm.logger.Info("Uninstalling arrow: %s", name)

	// Check if installed
	pkg, exists := pm.database.GetPackage(name)
	if !exists {
		return fmt.Errorf("arrow %s is not installed", name)
	}

	// Check if other packages depend on this one
	dependents := pm.database.GetDependents(name)
	if len(dependents) > 0 {
		return fmt.Errorf("cannot uninstall %s: packages %v depend on it", name, dependents)
	}

	// Load arrow
	arrow, err := pm.loadInstalledArrow(name)
	if err != nil {
		pm.logger.Warn("Failed to load arrow for uninstall, proceeding anyway: %v", err)
	} else {
		// Execute uninstall method
		if err := pm.executeMethod(arrow, "uninstall", pkg.InstallPath, pkg.Variables); err != nil {
			pm.logger.Warn("Failed to execute uninstall method: %v", err)
		}
	}

	// Remove installation directory
	if err := os.RemoveAll(pkg.InstallPath); err != nil {
		pm.logger.Warn("Failed to remove install directory: %v", err)
	}

	// Remove from database
	if err := pm.database.RemovePackage(name); err != nil {
		return fmt.Errorf("failed to remove package from database: %w", err)
	}

	// Clean up unused dependencies
	for _, dep := range pkg.Dependencies {
		dependents := pm.database.GetDependents(dep)
		if len(dependents) == 0 {
			pm.logger.Info("Dependency %s is no longer needed, consider removing it", dep)
		}
	}

	pm.logger.Info("Successfully uninstalled arrow: %s", name)
	return nil
}

// ExecuteArrow executes an arrow's execute method
func (pm *PackageManager) ExecuteArrow(name string, variables map[string]string) error {
	pm.logger.Info("Executing arrow: %s", name)

	// Check if installed
	pkg, exists := pm.database.GetPackage(name)
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
	arrow, err := pm.loadInstalledArrow(name)
	if err != nil {
		return fmt.Errorf("failed to load arrow: %w", err)
	}

	// Update status to running
	if err := pm.database.UpdatePackageStatus(name, "running"); err != nil {
		pm.logger.Warn("Failed to update package status: %v", err)
	}

	// Execute method
	err = pm.executeMethod(arrow, "execute", pkg.InstallPath, finalVariables)
	
	// Update status based on result
	status := "stopped"
	if err != nil {
		status = "error"
	}
	if statusErr := pm.database.UpdatePackageStatus(name, status); statusErr != nil {
		pm.logger.Warn("Failed to update package status: %v", statusErr)
	}

	return err
}

// UpdateArrow updates an arrow to the latest version
func (pm *PackageManager) UpdateArrow(name string) error {
	pm.logger.Info("Updating arrow: %s", name)

	// Check if installed
	pkg, exists := pm.database.GetPackage(name)
	if !exists {
		return fmt.Errorf("arrow %s is not installed", name)
	}

	// Get latest version from repository
	arrow, sourcePath, err := pm.repository.GetArrow(name)
	if err != nil {
		return fmt.Errorf("failed to get arrow %s: %w", name, err)
	}

	// Check if update is needed
	if arrow.ArrowVersion() == pkg.Version {
		pm.logger.Info("Arrow %s is already up to date", name)
		return nil
	}

	// Execute update method if available
	if err := pm.executeMethod(arrow, "update", pkg.InstallPath, pkg.Variables); err != nil {
		return fmt.Errorf("failed to execute update method: %w", err)
	}

	// Update arrow file
	arrowFile := filepath.Join(pkg.InstallPath, "arrow.yaml")
	arrowInfo := &ArrowInfo{
		Name: name,
		Path: sourcePath,
	}
	if err := pm.repository.DownloadArrow(arrowInfo, arrowFile); err != nil {
		return fmt.Errorf("failed to download updated arrow: %w", err)
	}

	// Update database
	pkg.Version = arrow.ArrowVersion()
	pkg.Repository = sourcePath
	pkg.UpdatedAt = time.Now()
	if err := pm.database.AddPackage(pkg); err != nil {
		return fmt.Errorf("failed to update package in database: %w", err)
	}

	pm.logger.Info("Successfully updated arrow: %s to version %s", name, arrow.ArrowVersion())
	return nil
}

// ValidateArrow validates an arrow installation
func (pm *PackageManager) ValidateArrow(name string) error {
	pm.logger.Info("Validating arrow: %s", name)

	// Check if installed
	pkg, exists := pm.database.GetPackage(name)
	if !exists {
		return fmt.Errorf("arrow %s is not installed", name)
	}

	// Load arrow
	arrow, err := pm.loadInstalledArrow(name)
	if err != nil {
		return fmt.Errorf("failed to load arrow: %w", err)
	}

	// Execute validate method
	return pm.executeMethod(arrow, "validate", pkg.InstallPath, pkg.Variables)
}

// GetInstalledArrows returns all installed arrows
func (pm *PackageManager) GetInstalledArrows() map[string]*InstalledPackage {
	return pm.database.GetAllPackages()
}

// GetArrowStatus returns the status of an arrow
func (pm *PackageManager) GetArrowStatus(name string) (string, error) {
	pkg, exists := pm.database.GetPackage(name)
	if !exists {
		return "", fmt.Errorf("arrow %s is not installed", name)
	}
	return pkg.Status, nil
}

// AddRepository adds a new repository
func (pm *PackageManager) AddRepository(repo string) {
	pm.repository.AddRepository(repo)
}

// RemoveRepository removes a repository
func (pm *PackageManager) RemoveRepository(repo string) {
	pm.repository.RemoveRepository(repo)
}

// GetRepositories returns all repositories
func (pm *PackageManager) GetRepositories() []string {
	return pm.repository.GetRepositories()
}

// Private methods

// loadInstalledArrow loads an arrow from its installation directory
func (pm *PackageManager) loadInstalledArrow(name string) (manifest.ArrowInterface, error) {
	pkg, exists := pm.database.GetPackage(name)
	if !exists {
		return nil, fmt.Errorf("arrow %s is not installed", name)
	}

	arrowFile := filepath.Join(pkg.InstallPath, "arrow.yaml")
	
	// Load arrow from file
	data, err := os.ReadFile(arrowFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read arrow file: %w", err)
	}

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

// executeMethod executes a specific method of an arrow
func (pm *PackageManager) executeMethod(arrow manifest.ArrowInterface, methodName, installPath string, variables map[string]string) error {
	methods := arrow.GetMethods()
	methodMap := methods.GetMethod(methodName)
	
	if methodMap == nil {
		return fmt.Errorf("method %s not found", methodName)
	}

	// Get platform-specific commands
	platform := runtime.GOOS
	commands, exists := methodMap[platform]
	if !exists {
		return fmt.Errorf("method %s not supported on platform %s", methodName, platform)
	}

	// Set up environment variables
	env := os.Environ()
	env = append(env, fmt.Sprintf("INSTALL_PATH=%s", installPath))
	
	// Add arrow variables
	arrowVars := arrow.GetVariables()
	for _, variable := range arrowVars {
		varName := variable.GetName()
		var varValue string
		
		// Use user-provided value if available, otherwise use default
		if userValue, exists := variables[varName]; exists {
			varValue = userValue
		} else if variable.GetDefault() != nil {
			varValue = fmt.Sprintf("%v", variable.GetDefault())
		}
		
		env = append(env, fmt.Sprintf("%s=%s", varName, varValue))
	}

	// Execute commands
	for _, command := range commands {
		// Expand variables in command
		expandedCommand := pm.expandVariables(command, variables, installPath)
		
		pm.logger.Debug("Executing command: %s", expandedCommand)
		
		// Split command for exec
		parts := strings.Fields(expandedCommand)
		if len(parts) == 0 {
			continue
		}

		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Env = env
		cmd.Dir = installPath
		
		// Capture output
		output, err := cmd.CombinedOutput()
		if err != nil {
			pm.logger.Error("Command failed: %s\nOutput: %s\nError: %v", expandedCommand, string(output), err)
			return fmt.Errorf("command failed: %s - %v", expandedCommand, err)
		}
		
		pm.logger.Debug("Command output: %s", string(output))
	}

	return nil
}

// expandVariables expands variables in a command string
func (pm *PackageManager) expandVariables(command string, variables map[string]string, installPath string) string {
	result := command
	
	// Replace INSTALL_PATH
	result = strings.ReplaceAll(result, "${INSTALL_PATH}", installPath)
	
	// Replace user variables
	for name, value := range variables {
		placeholder := fmt.Sprintf("${%s}", name)
		result = strings.ReplaceAll(result, placeholder, value)
	}
	
	return result
} 