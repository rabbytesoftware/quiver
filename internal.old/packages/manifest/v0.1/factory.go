package v0_1

import (
	"fmt"

	"github.com/rabbytesoftware/quiver/internal/packages/manifest"
	"gopkg.in/yaml.v3"
)

// init registers the v0.1 factory with the default registry
func Init() {
	manifest.RegisterFactory("0.1", NewFactory())
}

// Factory implements the ArrowFactory interface for v0.1 arrows
type Factory struct{}

// NewFactory creates a new v0.1 factory
func NewFactory() *Factory {
	return &Factory{}
}

// CreateArrow creates a v0.1 arrow from YAML data
func (f *Factory) CreateArrow(version string, data []byte) (manifest.ArrowInterface, error) {
	if version != "0.1" {
		return nil, fmt.Errorf("factory only supports version 0.1, got %s", version)
	}

	var arrow Arrow
	if err := yaml.Unmarshal(data, &arrow); err != nil {
		return nil, fmt.Errorf("failed to unmarshal v0.1 arrow: %w", err)
	}

	return &arrow, nil
}

// SupportedVersions returns the versions this factory supports
func (f *Factory) SupportedVersions() []string {
	return []string{"0.1"}
} 