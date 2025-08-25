package metadata

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed metadata.json
var embeddedMetadata []byte

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

// Load reads the embedded metadata and parses it into the Project variable
func Load() error {
	// Parse embedded JSON
	var metadata Metadata
	if err := json.Unmarshal(embeddedMetadata, &metadata); err != nil {
		return fmt.Errorf("failed to parse embedded metadata JSON: %w", err)
	}

	// Set the global Project variable
	Project = &metadata
	return nil
}

// LoadFromPath reads a metadata.json file from the specified path (for testing)
func LoadFromPath(path string) error {
	// This function is kept for backward compatibility and testing
	// but the main Load() function now uses embedded data
	return fmt.Errorf("LoadFromPath is deprecated, use Load() which uses embedded metadata")
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
