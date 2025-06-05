package v1_0

type variable struct {
	name      	string      `yaml:"name"`
	_default   	interface{} `yaml:"default"`
	values    	[]string    `yaml:"values,omitempty"`
	min       	*int        `yaml:"min,omitempty"`
	max       	*int        `yaml:"max,omitempty"`
	sensitive 	bool        `yaml:"sensitive,omitempty"`
}

func (v *variable) Name() *string {
	return &v.name
}

func (v *variable) Default() *interface{} {
	return &v._default
}

func (v *variable) Values() *[]string {
	return &v.values
}

func (v *variable) Min() *int {
	return v.min
}

func (v *variable) Max() *int {
	return v.max
}

func (v *variable) Sensitive() *bool {
	return &v.sensitive
}
