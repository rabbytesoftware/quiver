package database

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/rabbytesoftware/quiver/internal/logger"
	"github.com/rabbytesoftware/quiver/internal/packages/types"
)

// JSONDatabase implements the Database interface using JSON file storage
type JSONDatabase struct {
	dbPath   string
	database *types.PackageDatabase
	logger   *logger.Logger
}

// NewJSONDatabase creates a new JSON database implementation
func NewJSONDatabase(dbPath string, logger *logger.Logger) Database {
	return &JSONDatabase{
		dbPath: dbPath,
		logger: logger.WithService("json-database"),
		database: &types.PackageDatabase{
			Installed: make(map[string]*types.InstalledPackage),
			Updated:   time.Now(),
		},
	}
}

// Load loads the database from disk
func (db *JSONDatabase) Load() error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(db.dbPath), 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	// Check if database file exists
	if _, err := os.Stat(db.dbPath); os.IsNotExist(err) {
		db.logger.Info("Database file doesn't exist, creating new one")
		return db.Save()
	}

	// Read database file
	data, err := os.ReadFile(db.dbPath)
	if err != nil {
		return fmt.Errorf("failed to read database: %w", err)
	}

	// Parse JSON
	if err := json.Unmarshal(data, &db.database); err != nil {
		return fmt.Errorf("failed to parse database: %w", err)
	}

	db.logger.Info("Loaded database with %d installed packages", len(db.database.Installed))
	return nil
}

// Save saves the database to disk
func (db *JSONDatabase) Save() error {
	db.database.Updated = time.Now()

	data, err := json.MarshalIndent(db.database, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal database: %w", err)
	}

	if err := os.WriteFile(db.dbPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write database: %w", err)
	}

	return nil
}

// Package operations

// IsInstalled checks if a package is installed
func (db *JSONDatabase) IsInstalled(name string) bool {
	_, exists := db.database.Installed[name]
	return exists
}

// AddPackage adds a package to the database
func (db *JSONDatabase) AddPackage(pkg *types.InstalledPackage) error {
	db.database.Installed[pkg.Name] = pkg
	return db.Save()
}

// RemovePackage removes a package from the database
func (db *JSONDatabase) RemovePackage(name string) error {
	delete(db.database.Installed, name)
	return db.Save()
}

// GetPackage gets a package from the database
func (db *JSONDatabase) GetPackage(name string) (*types.InstalledPackage, bool) {
	pkg, exists := db.database.Installed[name]
	return pkg, exists
}

// GetAllPackages returns all installed packages
func (db *JSONDatabase) GetAllPackages() map[string]*types.InstalledPackage {
	return db.database.Installed
}

// UpdatePackage updates an existing package
func (db *JSONDatabase) UpdatePackage(pkg *types.InstalledPackage) error {
	if _, exists := db.database.Installed[pkg.Name]; !exists {
		return fmt.Errorf("package %s not found", pkg.Name)
	}
	
	pkg.UpdatedAt = time.Now()
	db.database.Installed[pkg.Name] = pkg
	return db.Save()
}

// Status operations

// UpdatePackageStatus updates the status of a package
func (db *JSONDatabase) UpdatePackageStatus(name string, status types.PackageStatus) error {
	if pkg, exists := db.database.Installed[name]; exists {
		pkg.Status = string(status)
		pkg.UpdatedAt = time.Now()
		return db.Save()
	}
	return fmt.Errorf("package %s not found", name)
}

// GetPackageStatus returns the status of a package
func (db *JSONDatabase) GetPackageStatus(name string) (types.PackageStatus, error) {
	if pkg, exists := db.database.Installed[name]; exists {
		return types.PackageStatus(pkg.Status), nil
	}
	return "", fmt.Errorf("package %s not found", name)
}

// GetPackagesByStatus returns packages with a specific status
func (db *JSONDatabase) GetPackagesByStatus(status types.PackageStatus) []*types.InstalledPackage {
	var packages []*types.InstalledPackage
	for _, pkg := range db.database.Installed {
		if pkg.Status == string(status) {
			packages = append(packages, pkg)
		}
	}
	return packages
}

// Dependency operations

// AddDependency adds a dependency relationship
func (db *JSONDatabase) AddDependency(packageName, dependencyName string) error {
	// Add to package's dependencies
	if pkg, exists := db.database.Installed[packageName]; exists {
		for _, dep := range pkg.Dependencies {
			if dep == dependencyName {
				return nil // Already exists
			}
		}
		pkg.Dependencies = append(pkg.Dependencies, dependencyName)
	}

	// Add to dependency's dependent_by
	if dep, exists := db.database.Installed[dependencyName]; exists {
		for _, dependent := range dep.DependentBy {
			if dependent == packageName {
				return db.Save() // Already exists
			}
		}
		dep.DependentBy = append(dep.DependentBy, packageName)
	}

	return db.Save()
}

// RemoveDependency removes a dependency relationship
func (db *JSONDatabase) RemoveDependency(packageName, dependencyName string) error {
	// Remove from package's dependencies
	if pkg, exists := db.database.Installed[packageName]; exists {
		for i, dep := range pkg.Dependencies {
			if dep == dependencyName {
				pkg.Dependencies = append(pkg.Dependencies[:i], pkg.Dependencies[i+1:]...)
				break
			}
		}
	}

	// Remove from dependency's dependent_by
	if dep, exists := db.database.Installed[dependencyName]; exists {
		for i, dependent := range dep.DependentBy {
			if dependent == packageName {
				dep.DependentBy = append(dep.DependentBy[:i], dep.DependentBy[i+1:]...)
				break
			}
		}
	}

	return db.Save()
}

// GetDependents returns packages that depend on the given package
func (db *JSONDatabase) GetDependents(packageName string) []string {
	if pkg, exists := db.database.Installed[packageName]; exists {
		return pkg.DependentBy
	}
	return []string{}
}

// GetDependencies returns the dependencies of a package
func (db *JSONDatabase) GetDependencies(packageName string) []string {
	if pkg, exists := db.database.Installed[packageName]; exists {
		return pkg.Dependencies
	}
	return []string{}
}

// HasDependents checks if a package has any dependents
func (db *JSONDatabase) HasDependents(packageName string) bool {
	dependents := db.GetDependents(packageName)
	return len(dependents) > 0
}

// GetDependencyTree returns the full dependency tree for a package
func (db *JSONDatabase) GetDependencyTree(packageName string) (map[string][]string, error) {
	if !db.IsInstalled(packageName) {
		return nil, fmt.Errorf("package %s is not installed", packageName)
	}

	tree := make(map[string][]string)
	visited := make(map[string]bool)
	
	db.buildDependencyTree(packageName, tree, visited)
	return tree, nil
}

// buildDependencyTree recursively builds the dependency tree
func (db *JSONDatabase) buildDependencyTree(packageName string, tree map[string][]string, visited map[string]bool) {
	if visited[packageName] {
		return // Avoid circular dependencies
	}
	
	visited[packageName] = true
	dependencies := db.GetDependencies(packageName)
	tree[packageName] = dependencies
	
	for _, dep := range dependencies {
		if db.IsInstalled(dep) {
			db.buildDependencyTree(dep, tree, visited)
		}
	}
}

// ValidateDependencies checks if all dependencies are satisfied
func (db *JSONDatabase) ValidateDependencies(packageName string) error {
	dependencies := db.GetDependencies(packageName)
	
	for _, dep := range dependencies {
		if !db.IsInstalled(dep) {
			return fmt.Errorf("dependency %s is not installed", dep)
		}
	}
	
	return nil
}

// GetUnusedDependencies returns dependencies that are not used by any package
func (db *JSONDatabase) GetUnusedDependencies() []string {
	var unused []string
	
	for packageName := range db.database.Installed {
		if len(db.GetDependents(packageName)) == 0 {
			// Check if this package is a root package (not a dependency of anything)
			isRootPackage := true
			for _, otherPkg := range db.database.Installed {
				for _, dep := range otherPkg.Dependencies {
					if dep == packageName {
						isRootPackage = false
						break
					}
				}
				if !isRootPackage {
					break
				}
			}
			
			if !isRootPackage {
				unused = append(unused, packageName)
			}
		}
	}
	
	return unused
} 