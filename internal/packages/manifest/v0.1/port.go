package v0_1

type Port struct {
	Name     string `yaml:"name"`
	Protocol string `yaml:"protocol"`
}

func (n *Port) GetName() string {
	return n.Name
}

func (n *Port) GetProtocol() string {
	return n.Protocol
} 