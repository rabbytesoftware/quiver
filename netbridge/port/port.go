package port

type Port struct {
	Name 		string
	Port     	uint16
	Host     	string
	Protocol 	string
}

func NewPort(
	name 		string,
	port 		uint16, 
	host 		string, 
	protocol 	string,
) Port {
	return Port{
		Name:     name,
		Port:     port,
		Host:     "", // TODO Host is intentionally left empty while developing the arrow system 
		Protocol: protocol,
	}
}
