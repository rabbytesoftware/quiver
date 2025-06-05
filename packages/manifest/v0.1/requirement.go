package v1_0

type requirement struct {
	cpuCores   		int 	`yaml:"cpu_cores"`
	ramGb   		int 	`yaml:"ram_gb"`
	diskGb   		int 	`yaml:"disk_gb"`
	networkMbps 	int 	`yaml:"network_mbps"`
}

func (r *requirement) GetCpuCores() *int {
	return &r.cpuCores
}

func (r *requirement) GetRamGb() *int {
	return &r.ramGb
}

func (r *requirement) GetDiskGb() *int {
	return &r.diskGb
}

func (r *requirement) GetNetworkMbps() *int {
	return &r.networkMbps
}