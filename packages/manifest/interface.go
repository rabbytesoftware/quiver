package manifest

import (
	Port "github.com/rabbytesoftware/quiver/netbridge/port"
)

type ArrowInterface interface {
	Manifest() 		string

	Name() 			string
	Description() 	string
	Mainteiners() 	[]string
	Credits() 		string
	License() 		string
	Repository() 	string
	Documentation() string
	ArrowVersion() 	string

	Requirements() 	[]Requirements

	Dependencies() 	[]string
	Netbridge() 	[]Port.Port

	Variables() 	[]Variables
	Methods() 		[]map[string]map[string][]string
}

type Requirements interface {
	CpuCores() 		int
	RamGb() 		int
	DiskGb() 		int
	NetworkMbps() 	int
}

type Variables interface {
	Name() 			string
	Default() 		string
	Values() 		[]string
	Min() 			int
	Max() 			int
	Sensitive() 	bool
}
