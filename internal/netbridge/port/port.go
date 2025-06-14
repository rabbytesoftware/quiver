package port

type Port struct {
	name 		string
	port     	uint16
	host     	string
	protocol 	string
}

func NewPort(
	name 		string,
	port 		uint16, 
	host 		string, 
	protocol 	string,
) Port {
	return Port{
		name:     name,
		port:     port,
		host:     "", // TODO Host is intentionally left empty while developing the arrow system 
		protocol: protocol,
	}
}

func (p Port) Name() string {
	return p.name
}

func (p Port) Port() uint16 {
	return p.port
}

func (p Port) Host() string {
	return p.host
}

func (p Port) Protocol() string {
	return p.protocol
}
