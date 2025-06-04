package models

type requirement struct {
	CpuCores   		int 	`yaml:"cpu_cores"`
	RamGb   		int 	`yaml:"ram_gb"`
	DiskGb   		int 	`yaml:"disk_gb"`
	NetworkMbps 	int 	`yaml:"network_mbps"`
}
