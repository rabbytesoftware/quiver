package packages

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rabbytesoftware/quiver/internal/logger"
)

// PackageDatabase represents the local package database
type PackageDatabase struct {
	Installed map[string]*InstalledPackage `json:"installed"`
	Updated   time.Time                    `json:"updated"`
}

// InstalledPackage represents an installed arrow package
type InstalledPackage struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Repository   string            `json:"repository"`
	InstallPath  string            `json:"install_path"`
	Dependencies []string          `json:"dependencies"`
	DependentBy  []string          `json:"dependent_by"`  // Arrows that depend on this one
	Variables    map[string]string `json:"variables"`     // User-configured variables
	Status       string            `json:"status"`        // installed, running, stopped, error
	InstalledAt  time.Time         `json:"installed_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

// DatabaseManager handles package database operations
type DatabaseManager struct {
	dbPath   string
	database *PackageDatabase
	logger   *logger.Logger
}

// NewDatabaseManager creates a new database manager
func NewDatabaseManager(dbPath string, logger *logger.Logger) *DatabaseManager {
	return &DatabaseManager{
		dbPath: dbPath,
		logger: logger.WithService("database"),
		database: &PackageDatabase{
			Installed: make(map[string]*InstalledPackage),
			Updated:   time.Now(),
		},
	}
}

// Load loads the database from disk
func (dm *DatabaseManager) Load() error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(dm.dbPath), 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	// Check if database file exists
	if _, err := os.Stat(dm.dbPath); os.IsNotExist(err) {
		dm.logger.Info("Database file doesn't exist, creating new one")
		return dm.Save()
	}

	// Read database file
	data, err := os.ReadFile(dm.dbPath)
	if err != nil {
		return fmt.Errorf("failed to read database: %w", err)
	}

	// Parse JSON
	if err := json.Unmarshal(data, &dm.database); err != nil {
		return fmt.Errorf("failed to parse database: %w", err)
	}

	dm.logger.Info("Loaded database with %d installed packages", len(dm.database.Installed))
	return nil
}

// Save saves the database to disk
func (dm *DatabaseManager) Save() error {
	dm.database.Updated = time.Now()

	data, err := json.MarshalIndent(dm.database, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal database: %w", err)
	}

	if err := os.WriteFile(dm.dbPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write database: %w", err)
	}

	return nil
}

// IsInstalled checks if a package is installed
func (dm *DatabaseManager) IsInstalled(name string) bool {
	_, exists := dm.database.Installed[name]
	return exists
}

// AddPackage adds a package to the database
func (dm *DatabaseManager) AddPackage(pkg *InstalledPackage) error {
	dm.database.Installed[pkg.Name] = pkg
	return dm.Save()
}

// RemovePackage removes a package from the database
func (dm *DatabaseManager) RemovePackage(name string) error {
	delete(dm.database.Installed, name)
	return dm.Save()
}

// GetPackage gets a package from the database
func (dm *DatabaseManager) GetPackage(name string) (*InstalledPackage, bool) {
	pkg, exists := dm.database.Installed[name]
	return pkg, exists
}

// GetAllPackages returns all installed packages
func (dm *DatabaseManager) GetAllPackages() map[string]*InstalledPackage {
	return dm.database.Installed
}

// UpdatePackageStatus updates the status of a package
func (dm *DatabaseManager) UpdatePackageStatus(name, status string) error {
	if pkg, exists := dm.database.Installed[name]; exists {
		pkg.Status = status
		pkg.UpdatedAt = time.Now()
		return dm.Save()
	}
	return fmt.Errorf("package %s not found", name)
}

// AddDependency adds a dependency relationship
func (dm *DatabaseManager) AddDependency(packageName, dependencyName string) error {
	// Add to package's dependencies
	if pkg, exists := dm.database.Installed[packageName]; exists {
		for _, dep := range pkg.Dependencies {
			if dep == dependencyName {
				return nil // Already exists
			}
		}
		pkg.Dependencies = append(pkg.Dependencies, dependencyName)
	}

	// Add to dependency's dependent_by
	if dep, exists := dm.database.Installed[dependencyName]; exists {
		for _, dependent := range dep.DependentBy {
			if dependent == packageName {
				return dm.Save() // Already exists
			}
		}
		dep.DependentBy = append(dep.DependentBy, packageName)
	}

	return dm.Save()
}

// RemoveDependency removes a dependency relationship
func (dm *DatabaseManager) RemoveDependency(packageName, dependencyName string) error {
	// Remove from package's dependencies
	if pkg, exists := dm.database.Installed[packageName]; exists {
		for i, dep := range pkg.Dependencies {
			if dep == dependencyName {
				pkg.Dependencies = append(pkg.Dependencies[:i], pkg.Dependencies[i+1:]...)
				break
			}
		}
	}

	// Remove from dependency's dependent_by
	if dep, exists := dm.database.Installed[dependencyName]; exists {
		for i, dependent := range dep.DependentBy {
			if dependent == packageName {
				dep.DependentBy = append(dep.DependentBy[:i], dep.DependentBy[i+1:]...)
				break
			}
		}
	}

	return dm.Save()
}

// GetDependents returns packages that depend on the given package
func (dm *DatabaseManager) GetDependents(packageName string) []string {
	if pkg, exists := dm.database.Installed[packageName]; exists {
		return pkg.DependentBy
	}
	return []string{}
}

// GetDependencies returns the dependencies of a package
func (dm *DatabaseManager) GetDependencies(packageName string) []string {
	if pkg, exists := dm.database.Installed[packageName]; exists {
		return pkg.Dependencies
	}
	return []string{}
} 