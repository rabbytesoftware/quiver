package database

import (
	"github.com/rabbytesoftware/quiver/internal/packages/types"
)

// Database interface defines the contract for package database operations
type Database interface {
	// Database lifecycle
	Load() error
	Save() error

	// Package operations
	IsInstalled(name string) bool
	AddPackage(pkg *types.InstalledPackage) error
	RemovePackage(name string) error
	GetPackage(name string) (*types.InstalledPackage, bool)
	GetAllPackages() map[string]*types.InstalledPackage
	UpdatePackage(pkg *types.InstalledPackage) error

	// Status operations
	UpdatePackageStatus(name string, status types.PackageStatus) error
	GetPackageStatus(name string) (types.PackageStatus, error)
	GetPackagesByStatus(status types.PackageStatus) []*types.InstalledPackage

	// Dependency operations
	AddDependency(packageName, dependencyName string) error
	RemoveDependency(packageName, dependencyName string) error
	GetDependents(packageName string) []string
	GetDependencies(packageName string) []string
	HasDependents(packageName string) bool
	GetDependencyTree(packageName string) (map[string][]string, error)
	ValidateDependencies(packageName string) error
	GetUnusedDependencies() []string
} 