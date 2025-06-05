package v1_0

import (
	port "github.com/rabbytesoftware/quiver/netbridge/port"
)

type Arrow struct {
	manifest 		string 			`yaml:"version"`
	metadata 		metadata 		`yaml:"metadata"`
	requirements 	struct{
		minimum 	requirement 	`yaml:"minimum"`
		recommended	requirement 	`yaml:"recommended"`
	}
	dependencies 	[]string 		`yaml:"dependencies"`
	netbridge 		[]port.Port 	`yaml:"netbridge"`
	variables 		[]variable		`yaml:"variables"`
	methods 		methods 		`yaml:"methods"`
}

func (a *Arrow) Manifest() string {
	return a.manifest
}

func (a *Arrow) Metadata() *metadata {
	return &a.metadata
}

func (a *Arrow) Requirements() *struct{
	minimum 	requirement 	`yaml:"minimum"`
	recommended	requirement 	`yaml:"recommended"`
} {
	return &a.requirements
}

func (a *Arrow) Dependencies() *[]string {
	return &a.dependencies
}

func (a *Arrow) Netbridge() *[]port.Port {
	return &a.netbridge
}

func (a *Arrow) Variables() *[]variable {
	return &a.variables
}

func (a *Arrow) Methods() *methods {
	return &a.methods
}
