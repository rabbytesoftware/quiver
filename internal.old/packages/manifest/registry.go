package manifest

import (
	"fmt"
	"sync"

	"gopkg.in/yaml.v3"
)

// VersionRegistry manages different manifest versions
type VersionRegistry struct {
	factories map[string]ArrowFactory
	mutex     sync.RWMutex
}

// NewVersionRegistry creates a new version registry
func NewVersionRegistry() *VersionRegistry {
	return &VersionRegistry{
		factories: make(map[string]ArrowFactory),
	}
}

// RegisterFactory registers a factory for a specific version
func (r *VersionRegistry) RegisterFactory(version string, factory ArrowFactory) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.factories[version] = factory
}

// LoadFromData loads an arrow manifest from data using the appropriate version factory
func (r *VersionRegistry) LoadFromData(data []byte) (ArrowInterface, error) {
	// Parse version first
	version, err := r.parseVersion(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse version: %w", err)
	}

	// Find factory for this version
	r.mutex.RLock()
	factory, exists := r.factories[version]
	r.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("unsupported arrow version: %s", version)
	}

	// Create arrow using the factory
	return factory.CreateArrow(version, data)
}

// parseVersion extracts version information from YAML data
func (r *VersionRegistry) parseVersion(data []byte) (string, error) {
	var versionInfo struct {
		Version string `yaml:"version"`
	}

	if err := yaml.Unmarshal(data, &versionInfo); err != nil {
		return "", fmt.Errorf("failed to parse version: %w", err)
	}

	if versionInfo.Version == "" {
		return "", fmt.Errorf("version field is required in arrow manifest")
	}

	return versionInfo.Version, nil
}

// GetSupportedVersions returns all supported versions
func (r *VersionRegistry) GetSupportedVersions() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	versions := make([]string, 0, len(r.factories))
	for version := range r.factories {
		versions = append(versions, version)
	}
	return versions
}

// HasVersion checks if a version is supported
func (r *VersionRegistry) HasVersion(version string) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	_, exists := r.factories[version]
	return exists
}

// DefaultRegistry is the global instance
var DefaultRegistry = NewVersionRegistry()

// RegisterFactory registers a factory in the default registry
func RegisterFactory(version string, factory ArrowFactory) {
	DefaultRegistry.RegisterFactory(version, factory)
}

// LoadFromData loads an arrow using the default registry
func LoadFromData(data []byte) (ArrowInterface, error) {
	return DefaultRegistry.LoadFromData(data)
} 