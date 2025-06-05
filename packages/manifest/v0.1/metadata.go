package v1_0

type metadata struct {
	name          string
	description   string
	maintainers   []string
	credits       []string
	license       string
	repository    string
	documentation string
	version       string
}

func (m *metadata) Name() *string {
	return &m.name
}

func (m *metadata) Description() *string {
	return &m.description
}

func (m *metadata) Maintainers() *[]string {
	return &m.maintainers
}

func (m *metadata) Credits() *[]string {
	return &m.credits
}

func (m *metadata) License() *string {
	return &m.license
}

func (m *metadata) Repository() *string {
	return &m.repository
}

func (m *metadata) Documentation() *string {
	return &m.documentation
}

func (m *metadata) Version() *string {
	return &m.version
}