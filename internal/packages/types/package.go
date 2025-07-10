package types

import (
	"time"
)

// ArrowInfo represents information about an available arrow
type ArrowInfo struct {
	Name           string `json:"name"`           // Package name without extension (e.g., "cs2")
	PackageName    string `json:"packagename"`    // Package name with extension (e.g., "cs2.yaml")
	RepositoryPath string `json:"repositorypath"` // Repository path where the package was found
	Path           string `json:"path"`           // Full path or URL to the package file
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

// PackageDatabase represents the local package database
type PackageDatabase struct {
	Installed map[string]*InstalledPackage `json:"installed"`
	Updated   time.Time                    `json:"updated"`
}

// RepositoryType represents the type of repository (local or remote)
type RepositoryType int

const (
	RepositoryTypeLocal RepositoryType = iota
	RepositoryTypeRemote
)

// Repository represents a package repository
type Repository struct {
	Path string
	Type RepositoryType
}

// PackageStatus represents the status of a package
type PackageStatus string

const (
	StatusInstalled PackageStatus = "installed"
	StatusRunning   PackageStatus = "running"
	StatusStopped   PackageStatus = "stopped"
	StatusError     PackageStatus = "error"
)

// ExecutionContext contains context information for command execution
type ExecutionContext struct {
	InstallPath string
	Variables   map[string]string
	Environment []string
}

// MethodType represents the type of arrow method
type MethodType string

const (
	MethodInstall   MethodType = "install"
	MethodExecute   MethodType = "execute"
	MethodUninstall MethodType = "uninstall"
	MethodUpdate    MethodType = "update"
	MethodValidate  MethodType = "validate"
) 