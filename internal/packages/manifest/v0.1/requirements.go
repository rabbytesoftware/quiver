package v0_1

type Requirements struct {
	Minimum     Requirement `yaml:"minimum"`
	Recommended Requirement `yaml:"recommended"`
}

func (r *Requirements) GetMinimum() Requirement {
	return r.Minimum
}

func (r *Requirements) GetRecommended() Requirement {
	return r.Recommended
} 