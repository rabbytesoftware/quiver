package v0_2

import (
	"fmt"

	"github.com/rabbytesoftware/quiver/internal/packages/manifest"
	"gopkg.in/yaml.v3"
)

// Init registers the v0.2 factory with the default registry
func Init() {
	manifest.RegisterFactory("0.2", NewFactory())
}

// Factory implements the ArrowFactory interface for v0.2 arrows
type Factory struct{}

// NewFactory creates a new v0.2 factory
func NewFactory() *Factory {
	return &Factory{}
}

// CreateArrow creates a v0.2 arrow from YAML data
func (f *Factory) CreateArrow(version string, data []byte) (manifest.ArrowInterface, error) {
	if version != "0.2" {
		return nil, fmt.Errorf("factory only supports version 0.2, got %s", version)
	}

	var arrow Arrow
	if err := yaml.Unmarshal(data, &arrow); err != nil {
		return nil, fmt.Errorf("failed to unmarshal v0.2 arrow: %w", err)
	}

	return &arrow, nil
}

// SupportedVersions returns the versions this factory supports
func (f *Factory) SupportedVersions() []string {
	return []string{"0.2"}
} 