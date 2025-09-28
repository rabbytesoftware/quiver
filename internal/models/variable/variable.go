package variable

type Variable struct {
	Name      string       `json:"name"`
	Default   string       `json:"default"`
	Values    []string     `json:"values"`
	Min       int          `json:"min"`
	Max       int          `json:"max"`
	Sensitive bool         `json:"sensitive"`
	Type      VariableType `json:"type"`
}
