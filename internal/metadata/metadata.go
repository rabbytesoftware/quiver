package metadata

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Maintainer represents a project maintainer
type Maintainer struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	URL   string `json:"url"`
}

// Metadata represents the project metadata structure
type Metadata struct {
	Version     string       `json:"version"`
	Description string       `json:"description"`
	Author      string       `json:"author"`
	URL         string       `json:"url"`
	License     string       `json:"license"`
	Copyright   string       `json:"copyright"`
	Maintainers []Maintainer `json:"maintainers"`
}

var (
	// Project holds the loaded project metadata
	Project *Metadata
)

// Load reads the metadata.json file and parses it into the Project variable
func Load() error {
	return LoadFromPath("metadata.json")
}

// LoadFromPath reads a metadata.json file from the specified path
func LoadFromPath(path string) error {
	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for %s: %w", path, err)
	}

	// Read the file
	data, err := os.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("failed to read metadata file %s: %w", absPath, err)
	}

	// Parse JSON
	var metadata Metadata
	if err := json.Unmarshal(data, &metadata); err != nil {
		return fmt.Errorf("failed to parse metadata JSON: %w", err)
	}

	// Set the global Project variable
	Project = &metadata
	return nil
}

// GetVersion returns the project version
func GetVersion() string {
	if Project == nil {
		return "unknown"
	}
	return Project.Version
}

// GetDescription returns the project description
func GetDescription() string {
	if Project == nil {
		return ""
	}
	return Project.Description
}

// GetAuthor returns the project author
func GetAuthor() string {
	if Project == nil {
		return ""
	}
	return Project.Author
}

// GetURL returns the project URL
func GetURL() string {
	if Project == nil {
		return ""
	}
	return Project.URL
}

// GetLicense returns the project license
func GetLicense() string {
	if Project == nil {
		return ""
	}
	return Project.License
}

// GetCopyright returns the project copyright
func GetCopyright() string {
	if Project == nil {
		return ""
	}
	return Project.Copyright
}

// GetMaintainers returns the list of project maintainers
func GetMaintainers() []Maintainer {
	if Project == nil {
		return nil
	}
	return Project.Maintainers
}

// String returns a formatted string representation of the metadata
func (m *Metadata) String() string {
	return fmt.Sprintf("%s v%s by %s", m.Description, m.Version, m.Author)
}
