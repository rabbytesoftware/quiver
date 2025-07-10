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

// Manager handles package database operations
type Manager struct {
	dbPath   string
	database *types.PackageDatabase
	logger   *logger.Logger
}

// NewManager creates a new database manager
func NewManager(dbPath string, logger *logger.Logger) *Manager {
	return &Manager{
		dbPath: dbPath,
		logger: logger.WithService("database"),
		database: &types.PackageDatabase{
			Installed: make(map[string]*types.InstalledPackage),
			Updated:   time.Now(),
		},
	}
}

// Load loads the database from disk
func (m *Manager) Load() error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(m.dbPath), 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}

	// Check if database file exists
	if _, err := os.Stat(m.dbPath); os.IsNotExist(err) {
		m.logger.Info("Database file doesn't exist, creating new one")
		return m.Save()
	}

	// Read database file
	data, err := os.ReadFile(m.dbPath)
	if err != nil {
		return fmt.Errorf("failed to read database: %w", err)
	}

	// Parse JSON
	if err := json.Unmarshal(data, &m.database); err != nil {
		return fmt.Errorf("failed to parse database: %w", err)
	}

	m.logger.Info("Loaded database with %d installed packages", len(m.database.Installed))
	return nil
}

// Save saves the database to disk
func (m *Manager) Save() error {
	m.database.Updated = time.Now()

	data, err := json.MarshalIndent(m.database, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal database: %w", err)
	}

	if err := os.WriteFile(m.dbPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write database: %w", err)
	}

	return nil
}

// Package operations

// IsInstalled checks if a package is installed
func (m *Manager) IsInstalled(name string) bool {
	_, exists := m.database.Installed[name]
	return exists
}

// AddPackage adds a package to the database
func (m *Manager) AddPackage(pkg *types.InstalledPackage) error {
	m.database.Installed[pkg.Name] = pkg
	return m.Save()
}

// RemovePackage removes a package from the database
func (m *Manager) RemovePackage(name string) error {
	delete(m.database.Installed, name)
	return m.Save()
}

// GetPackage gets a package from the database
func (m *Manager) GetPackage(name string) (*types.InstalledPackage, bool) {
	pkg, exists := m.database.Installed[name]
	return pkg, exists
}

// GetAllPackages returns all installed packages
func (m *Manager) GetAllPackages() map[string]*types.InstalledPackage {
	return m.database.Installed
}

// UpdatePackage updates an existing package
func (m *Manager) UpdatePackage(pkg *types.InstalledPackage) error {
	if _, exists := m.database.Installed[pkg.Name]; !exists {
		return fmt.Errorf("package %s not found", pkg.Name)
	}
	
	pkg.UpdatedAt = time.Now()
	m.database.Installed[pkg.Name] = pkg
	return m.Save()
}

// Status operations

// UpdatePackageStatus updates the status of a package
func (m *Manager) UpdatePackageStatus(name string, status types.PackageStatus) error {
	if pkg, exists := m.database.Installed[name]; exists {
		pkg.Status = string(status)
		pkg.UpdatedAt = time.Now()
		return m.Save()
	}
	return fmt.Errorf("package %s not found", name)
}

// GetPackageStatus returns the status of a package
func (m *Manager) GetPackageStatus(name string) (types.PackageStatus, error) {
	if pkg, exists := m.database.Installed[name]; exists {
		return types.PackageStatus(pkg.Status), nil
	}
	return "", fmt.Errorf("package %s not found", name)
}

// GetPackagesByStatus returns packages with a specific status
func (m *Manager) GetPackagesByStatus(status types.PackageStatus) []*types.InstalledPackage {
	var packages []*types.InstalledPackage
	for _, pkg := range m.database.Installed {
		if pkg.Status == string(status) {
			packages = append(packages, pkg)
		}
	}
	return packages
} 