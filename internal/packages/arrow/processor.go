package arrow

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rabbytesoftware/quiver/internal/logger"
	"github.com/rabbytesoftware/quiver/internal/packages/manifest"
	v0_1 "github.com/rabbytesoftware/quiver/internal/packages/manifest/v0.1"
	"gopkg.in/yaml.v3"
)

// Processor handles arrow loading and processing operations
type Processor struct {
	logger *logger.Logger
}

// NewProcessor creates a new arrow processor
func NewProcessor(logger *logger.Logger) *Processor {
	return &Processor{
		logger: logger.WithService("arrow-processor"),
	}
}

// LoadFromFile loads an arrow manifest from a file
func (p *Processor) LoadFromFile(path string) (manifest.ArrowInterface, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read arrow file: %w", err)
	}

	return p.LoadFromData(data)
}

// LoadFromData loads an arrow manifest from data
func (p *Processor) LoadFromData(data []byte) (manifest.ArrowInterface, error) {
	// Parse version first
	var versionInfo struct {
		Version string `yaml:"version"`
	}

	if err := yaml.Unmarshal(data, &versionInfo); err != nil {
		return nil, fmt.Errorf("failed to parse version: %w", err)
	}

	if versionInfo.Version == "" {
		return nil, fmt.Errorf("version field is required in arrow manifest")
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
		return nil, fmt.Errorf("unsupported arrow version: %s", versionInfo.Version)
	}
}

// LoadFromInstallation loads an arrow from its installation directory
func (p *Processor) LoadFromInstallation(installPath string) (manifest.ArrowInterface, error) {
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
func (p *Processor) ValidateArrow(arrow manifest.ArrowInterface) error {
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

	// Check if at least install and execute methods exist
	if methods.GetInstall() == nil {
		return fmt.Errorf("install method is required")
	}

	if methods.GetExecute() == nil {
		return fmt.Errorf("execute method is required")
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

	// Parse basic info only
	var basicInfo struct {
		Version  string `yaml:"version"`
		Metadata struct {
			Name        string `yaml:"name"`
			Description string `yaml:"description"`
			Version     string `yaml:"version"`
		} `yaml:"metadata"`
	}

	if err := yaml.Unmarshal(data, &basicInfo); err != nil {
		return nil, fmt.Errorf("failed to parse arrow basic info: %w", err)
	}

	return &ArrowBasicInfo{
		Name:        basicInfo.Metadata.Name,
		Description: basicInfo.Metadata.Description,
		Version:     basicInfo.Metadata.Version,
		ManifestVersion: basicInfo.Version,
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

// SupportedVersions returns a list of supported arrow manifest versions
func (p *Processor) SupportedVersions() []string {
	return []string{"0.1"}
} 