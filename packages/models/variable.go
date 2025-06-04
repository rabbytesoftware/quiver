package models

type variable struct {
	Name      string      `yaml:"name" json:"name"`
	Default   interface{} `yaml:"default" json:"default"`
	Values    []string    `yaml:"values,omitempty" json:"values,omitempty"`
	Min       *int        `yaml:"min,omitempty" json:"min,omitempty"`
	Max       *int        `yaml:"max,omitempty" json:"max,omitempty"`
	Sensitive bool        `yaml:"sensitive,omitempty" json:"sensitive,omitempty"`
}
