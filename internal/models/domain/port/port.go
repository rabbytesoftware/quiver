package port

type PortRule struct {
	ID        			string 				`json:"id"`
	StartPort 			int 				`json:"start_port"`
	EndPort   			int 				`json:"end_port"`
	Protocol  			Protocol 			`json:"protocol"`
	ForwardingStatus 	ForwardingStatus 	`json:"forwarding_status"`
}

func (p* PortRule) IsStartPortValid() bool {
	return p.StartPort > 0 && p.StartPort <= 65535
}

func (p* PortRule) IsEndPortValid() bool {
	return p.EndPort > 0 && p.EndPort <= 65535
}
