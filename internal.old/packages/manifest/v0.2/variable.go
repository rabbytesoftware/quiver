package v0_2

type Variable struct {
	Name      string      `yaml:"name"`
	Default   interface{} `yaml:"default"`
	Values    []string    `yaml:"values,omitempty"`
	Min       *int        `yaml:"min,omitempty"`
	Max       *int        `yaml:"max,omitempty"`
	Sensitive bool        `yaml:"sensitive,omitempty"`
}

func (v *Variable) GetName() string {
	return v.Name
}

func (v *Variable) GetDefault() interface{} {
	return v.Default
}

func (v *Variable) GetValues() []string {
	return v.Values
}

func (v *Variable) GetMin() *int {
	return v.Min
}

func (v *Variable) GetMax() *int {
	return v.Max
}

func (v *Variable) GetSensitive() bool {
	return v.Sensitive
} 