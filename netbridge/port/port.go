package port

type Port struct {
	Port     uint16
	Host     string
	Protocol string
}

func NewPort(port uint16, host string, protocol string) Port {
	return Port{
		Port:     port,
		Host:     host,
		Protocol: protocol,
	}
}
