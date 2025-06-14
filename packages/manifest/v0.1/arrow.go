package v0_1

type Arrow struct {
	Version      string       `yaml:"version"`
	Metadata     Metadata     `yaml:"metadata"`
	Requirements Requirements `yaml:"requirements"`
	Dependencies []string     `yaml:"dependencies"`
	Netbridge    []Port 	  `yaml:"netbridge"`
	Variables    []Variable   `yaml:"variables"`
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

func (a *Arrow) Mainteiners() []string {
	return a.Metadata.Maintainers
}

func (a *Arrow) Credits() string {
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

func (a *Arrow) GetMetadata() *Metadata {
	return &a.Metadata
}

func (a *Arrow) GetRequirements() Requirements {
	return a.Requirements
}

func (a *Arrow) GetDependencies() []string {
	return a.Dependencies
}

func (a *Arrow) GetNetbridge() []Port {
	return a.Netbridge
}

func (a *Arrow) GetVariables() []Variable {
	return a.Variables
}

func (a *Arrow) GetMethods() Methods {
	return a.Methods
}
