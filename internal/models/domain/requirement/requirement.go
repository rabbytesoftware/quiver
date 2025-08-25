package requirement

import (
	system "github.com/rabbytesoftware/quiver/internal/models/domain/system"
)

type Requirement struct {
	CpuCores int `json:"cpu_cores"`
	Memory   int `json:"memory"`
	Disk     int `json:"disk"`
	OS       system.OS `json:"os"`
}

func (r* Requirement) IsValid() bool {
	return r.CpuCores > 0 && r.Memory > 0 && r.Disk > 0 && r.OS.IsValid()
}
