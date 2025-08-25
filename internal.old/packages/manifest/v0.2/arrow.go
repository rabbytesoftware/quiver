package v0_2

import "github.com/rabbytesoftware/quiver/internal/packages/manifest"

type Arrow struct {
	Version      string       `yaml:"version"`
	Metadata     Metadata     `yaml:"metadata"`
	Requirements Requirements `yaml:"requirements"`
	Dependencies []string     `yaml:"dependencies,omitempty"`
	Netbridge    []Port       `yaml:"netbridge,omitempty"`
	Variables    []Variable   `yaml:"variables,omitempty"`
	Methods      Methods      `yaml:"methods"`
}

func (a *Arrow) Manifest() string {
	return a.Version
}

func (a *Arrow) Name() string {
	return a.Metadata.Name
}

func (a *Arrow) Description() string {
	return a.Metadata.Description
}

func (a *Arrow) Maintainers() []string {
	return a.Metadata.Maintainers
}

func (a *Arrow) Credits() string {
	if len(a.Metadata.Credits) > 0 {
		return a.Metadata.Credits[0]
	}
	return ""
}

func (a *Arrow) License() string {
	return a.Metadata.License
}

func (a *Arrow) Repository() string {
	return a.Metadata.Repository
}

func (a *Arrow) Documentation() string {
	return a.Metadata.Documentation
}

func (a *Arrow) ArrowVersion() string {
	return a.Metadata.Version
}

func (a *Arrow) GetMetadata() manifest.Metadata {
	return &a.Metadata
}

func (a *Arrow) GetRequirements() manifest.Requirements {
	return &a.Requirements
}

func (a *Arrow) GetDependencies() []string {
	return a.Dependencies
}

func (a *Arrow) GetNetbridge() []manifest.Netbridge {
	bridges := make([]manifest.Netbridge, len(a.Netbridge))
	for i, port := range a.Netbridge {
		bridges[i] = &port
	}
	return bridges
}

func (a *Arrow) GetVariables() []manifest.Variable {
	vars := make([]manifest.Variable, len(a.Variables))
	for i, variable := range a.Variables {
		vars[i] = &variable
	}
	return vars
}

func (a *Arrow) GetMethods() manifest.Methods {
	return &a.Methods
}

func (a *Arrow) GetSupportedArchs(os string) []string {
	return a.Requirements.Compatible[os]
} 