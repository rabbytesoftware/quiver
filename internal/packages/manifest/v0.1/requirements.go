package v0_1

import (
	"github.com/rabbytesoftware/quiver/internal/packages/manifest"
)

type Requirements struct {
	Minimum     Requirement `yaml:"minimum"`
	Recommended Requirement `yaml:"recommended"`
	Compatible  map[string][]string `yaml:"compatible"`
}

func (r *Requirements) GetMinimum() manifest.Requirement {
	return &r.Minimum
}

func (r *Requirements) GetRecommended() manifest.Requirement {
	return &r.Recommended
}

func (r *Requirements) GetCompatible() map[string][]string {
	return r.Compatible
} 