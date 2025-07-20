package v0_2

import (
	"strings"

	"github.com/rabbytesoftware/quiver/internal/packages/manifest"
)

type Requirements struct {
	System      []string            `yaml:"system,omitempty"`
	Minimum     Requirement         `yaml:"minimum"`
	Recommended Requirement         `yaml:"recommended"`
	Compatible  map[string][]string `yaml:"compatible,omitempty"`
}

func (r *Requirements) GetMinimum() manifest.Requirement {
	return &r.Minimum
}

func (r *Requirements) GetRecommended() manifest.Requirement {
	return &r.Recommended
}

func (r *Requirements) GetCompatible() map[string][]string {
	// If Compatible is already set, use it
	if r.Compatible != nil && len(r.Compatible) > 0 {
		return r.Compatible
	}
	
	// Convert System format to Compatible format
	if r.System != nil && len(r.System) > 0 {
		compatible := make(map[string][]string)
		for _, system := range r.System {
			parts := strings.Split(system, "/")
			if len(parts) == 2 {
				os := parts[0]
				arch := parts[1]
				if _, exists := compatible[os]; !exists {
					compatible[os] = []string{}
				}
				// Check if arch is not already in the slice
				found := false
				for _, existingArch := range compatible[os] {
					if existingArch == arch {
						found = true
						break
					}
				}
				if !found {
					compatible[os] = append(compatible[os], arch)
				}
			}
		}
		return compatible
	}
	
	return nil
} 