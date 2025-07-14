package manifest

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/rabbytesoftware/quiver/internal/logger"
)

// Processor handles arrow loading and processing operations
type Processor struct {
	logger *logger.Logger
	registry *VersionRegistry
}

// NewProcessor creates a new arrow processor
func NewProcessor(logger *logger.Logger) *Processor {
	return &Processor{
		logger: logger.WithService("manifest-processor"),
		registry: DefaultRegistry,
	}
}

// LoadFromFile loads an arrow manifest from a file
func (p *Processor) LoadFromFile(path string) (ArrowInterface, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read arrow file: %w", err)
	}

	return p.LoadFromData(data)
}

// LoadFromData loads an arrow manifest from data
func (p *Processor) LoadFromData(data []byte) (ArrowInterface, error) {
	return p.registry.LoadFromData(data)
}

// LoadFromInstallation loads an arrow from its installation directory
func (p *Processor) LoadFromInstallation(installPath string) (ArrowInterface, error) {
	arrowFile := filepath.Join(installPath, "arrow.yaml")
	
	// Check if arrow.yaml exists
	if _, err := os.Stat(arrowFile); os.IsNotExist(err) {
		// Try arrow.yml as fallback
		arrowFile = filepath.Join(installPath, "arrow.yml")
		if _, err := os.Stat(arrowFile); os.IsNotExist(err) {
			return nil, fmt.Errorf("arrow manifest not found in installation directory: %s", installPath)
		}
	}

	return p.LoadFromFile(arrowFile)
}

// ValidateArrow validates an arrow manifest
func (p *Processor) ValidateArrow(arrow ArrowInterface) error {
	// Basic validation
	if arrow.Name() == "" {
		return fmt.Errorf("arrow name is required")
	}

	if arrow.ArrowVersion() == "" {
		return fmt.Errorf("arrow version is required")
	}

	// Validate methods exist
	methods := arrow.GetMethods()
	if methods == nil {
		return fmt.Errorf("arrow methods are required")
	}

	osStr := runtime.GOOS
	archStr := runtime.GOARCH

	// Check install method for current platform
	install := methods.GetInstall()
	if install == nil || install[osStr] == nil || install[osStr][archStr] == nil || len(install[osStr][archStr]) == 0 {
		return fmt.Errorf("no install method defined for %s/%s", osStr, archStr)
	}

	// Check execute method for current platform
	execute := methods.GetExecute()
	if execute == nil || execute[osStr] == nil || execute[osStr][archStr] == nil || len(execute[osStr][archStr]) == 0 {
		return fmt.Errorf("no execute method defined for %s/%s", osStr, archStr)
	}

	// Check compatible requirements
	requirements := arrow.GetRequirements()
	if requirements != nil {
		compatible := requirements.GetCompatible()
		if len(compatible) > 0 {
			archs, ok := compatible[osStr]
			if !ok {
				return fmt.Errorf("os %s not in compatible list", osStr)
			}
			found := false
			for _, a := range archs {
				if a == archStr {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("arch %s not in compatible list for %s", archStr, osStr)
			}
		}
	}

	// Validate variables if present
	for _, variable := range arrow.GetVariables() {
		if variable.GetName() == "" {
			return fmt.Errorf("variable name is required")
		}
	}

	p.logger.Debug("Arrow validation passed for: %s", arrow.Name())
	return nil
}

// ValidateArrowFile validates an arrow file without fully loading it
func (p *Processor) ValidateArrowFile(path string) error {
	arrow, err := p.LoadFromFile(path)
	if err != nil {
		return fmt.Errorf("failed to load arrow for validation: %w", err)
	}

	return p.ValidateArrow(arrow)
}

// GetArrowInfo extracts basic information from an arrow without full loading
func (p *Processor) GetArrowInfo(path string) (*ArrowBasicInfo, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read arrow file: %w", err)
	}

	// Load the arrow to get info
	arrow, err := p.LoadFromData(data)
	if err != nil {
		return nil, fmt.Errorf("failed to load arrow: %w", err)
	}

	return &ArrowBasicInfo{
		Name:        arrow.Name(),
		Description: arrow.Description(),
		Version:     arrow.ArrowVersion(),
		ManifestVersion: arrow.Manifest(),
		Path:        path,
	}, nil
}

// ArrowBasicInfo contains basic arrow information
type ArrowBasicInfo struct {
	Name            string
	Description     string
	Version         string
	ManifestVersion string
	Path            string
}

// IsValidArrowFile checks if a file is a valid arrow manifest
func (p *Processor) IsValidArrowFile(path string) bool {
	// Check file extension
	ext := filepath.Ext(path)
	if ext != ".yaml" && ext != ".yml" {
		return false
	}

	// Try to parse basic structure
	err := p.ValidateArrowFile(path)
	return err == nil
}

// GetSupportedVersions returns all supported manifest versions
func (p *Processor) GetSupportedVersions() []string {
	return p.registry.GetSupportedVersions()
}

// HasVersion checks if a version is supported
func (p *Processor) HasVersion(version string) bool {
	return p.registry.HasVersion(version)
} 