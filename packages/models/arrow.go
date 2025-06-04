package models

import (
	port "github.com/rabbytesoftware/quiver/netbridge/port"
)

type Arrow struct {
	Metadata 		metadata 		`yaml:"metadata"`
	Requirements 	struct{
		Minimum 	requirement 	`yaml:"minimum"`
		Mecommended	requirement 	`yaml:"recommended"`
	}
	Dependencies 	[]string 		`yaml:"dependencies"`
	Netbridge 		[]port.Port 	`yaml:"netbridge"`
	Variables 		[]variable		`yaml:"variables"`
	Methods 		methods 		`yaml:"methods"`
}

