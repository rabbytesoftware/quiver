package models

type Query struct {
	Syntax      string            `yaml:"syntax"`
	Description string            `yaml:"description"`
	REST        *REST             `yaml:"REST,omitempty"`
	Children    []Query           `yaml:"children,omitempty"`
	Variables   map[string]string `yaml:"variables,omitempty"`
}

type REST struct {
	URL    string `yaml:"url"`
	Method string `yaml:"method"`
	Body   string `yaml:"body,omitempty"`
}

type QueriesConfig struct {
	Queries []Query `yaml:"queries"`
}

type QueryRequest struct {
	URL    string
	Method string
	Args   map[string]string
}
