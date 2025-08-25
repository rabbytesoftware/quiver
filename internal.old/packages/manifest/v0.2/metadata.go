package v0_2

type Metadata struct {
	Name          string   `yaml:"name"`
	Description   string   `yaml:"description"`
	Maintainers   []string `yaml:"mainteiners"`
	Credits       []string `yaml:"credits"`
	License       string   `yaml:"license"`
	Repository    string   `yaml:"repository"`
	Documentation string   `yaml:"documentation"`
	Version       string   `yaml:"version"`
}

func (m *Metadata) GetName() string {
	return m.Name
}

func (m *Metadata) GetDescription() string {
	return m.Description
}

func (m *Metadata) GetMaintainers() []string {
	return m.Maintainers
}

func (m *Metadata) GetCredits() []string {
	return m.Credits
}

func (m *Metadata) GetLicense() string {
	return m.License
}

func (m *Metadata) GetRepository() string {
	return m.Repository
}

func (m *Metadata) GetDocumentation() string {
	return m.Documentation
}

func (m *Metadata) GetVersion() string {
	return m.Version
} 