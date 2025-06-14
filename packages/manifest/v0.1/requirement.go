package v0_1

type Requirement struct {
	CpuCores    int `yaml:"cpu_cores"`
	RamGB       int `yaml:"ram_gb"`
	DiskGB      int `yaml:"disk_gb"`
	NetworkMbps int `yaml:"network_mbps"`
}

func (r *Requirement) GetCpuCores() int {
	return r.CpuCores
}

func (r *Requirement) GetRamGB() int {
	return r.RamGB
}

func (r *Requirement) GetDiskGB() int {
	return r.DiskGB
}

func (r *Requirement) GetNetworkMbps() int {
	return r.NetworkMbps
}