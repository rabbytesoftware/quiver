package requirements

import (
	"runtime"

	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/disk"
)

func GetSystemRules() map[string]interface{} {
	// OS & Arch
	osRules := map[string]interface{}{
		"allowed": []string{"windows", "linux"},
		"arch":    []string{"amd64", "arm64"},
	}

	// Memory
	vm, _ := mem.VirtualMemory()
	memMB := int(vm.Total / 1024 / 1024)

	// Disk
	path := "C:\\"
	if runtime.GOOS != "windows" {
		path = "/"
	}
	usage, _ := disk.Usage(path)
	totalDiskGB := int(usage.Total / 1024 / 1024 / 1024)

	// Require only 80% of total disk
	requiredDiskGB := int(float64(totalDiskGB) * 0.8)

	return map[string]interface{}{
		"os": osRules,
		"resources": map[string]interface{}{
			"memory_mb": memMB,
			"disk_gb":   requiredDiskGB,
		},
		"dependencies": map[string]interface{}{
			"required_bins": []string{"git", "docker"},
		},
	}
}
