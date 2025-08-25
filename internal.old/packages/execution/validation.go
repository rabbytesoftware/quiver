package execution

import (
	"fmt"
	"os"
	"runtime"

	"github.com/rabbytesoftware/quiver/internal/packages/manifest"
	"github.com/rabbytesoftware/quiver/internal/packages/types"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

// ValidateRequirements checks hardware requirements against system capabilities with verbose logging
func (e *Engine) ValidateRequirements(arrow manifest.ArrowInterface) error {
	e.logger.Info("Starting hardware requirements validation for arrow: %s", arrow.Name())
	
	requirements := arrow.GetRequirements()
	if requirements == nil {
		e.logger.Info("No requirements specified for arrow %s, skipping validation", arrow.Name())
		return nil
	}
	
	minimum := requirements.GetMinimum()
	if minimum == nil {
		e.logger.Info("No minimum requirements specified for arrow %s, skipping validation", arrow.Name())
		return nil
	}
	
	e.logger.Info("Validating hardware requirements against system capabilities...")
	
	// Check CPU cores
	if reqCores := minimum.GetCpuCores(); reqCores > 0 {
		e.logger.Info("Checking CPU cores requirement: %d cores needed", reqCores)
		
		cpuInfo, err := cpu.Info()
		if err != nil {
			e.logger.Error("Failed to get CPU info: %v", err)
			e.logger.Warn("Skipping CPU cores validation due to error")
		} else {
			totalCores := 0
			for i, cpuData := range cpuInfo {
				cores := int(cpuData.Cores)
				totalCores += cores
				e.logger.Debug("CPU %d: %s (%d cores)", i, cpuData.ModelName, cores)
			}
			
			e.logger.Info("System has %d total CPU cores", totalCores)
			
			if totalCores < reqCores {
				e.logger.Error("Insufficient CPU cores: required %d, available %d", reqCores, totalCores)
				return fmt.Errorf("insufficient CPU cores: required %d, available %d", reqCores, totalCores)
			}
			e.logger.Info("CPU cores requirement satisfied: %d/%d ✓", totalCores, reqCores)
		}
	} else {
		e.logger.Debug("No CPU cores requirement specified")
	}
	
	// Check RAM
	if reqRAM := minimum.GetRamGB(); reqRAM > 0 {
		e.logger.Info("Checking RAM requirement: %dGB needed", reqRAM)
		
		memInfo, err := mem.VirtualMemory()
		if err != nil {
			e.logger.Error("Failed to get memory info: %v", err)
			e.logger.Warn("Skipping RAM validation due to error")
		} else {
			totalGB := int(memInfo.Total / (1024 * 1024 * 1024))
			usedGB := int(memInfo.Used / (1024 * 1024 * 1024))
			availableGB := int(memInfo.Available / (1024 * 1024 * 1024))
			
			e.logger.Info("System RAM: %dGB total, %dGB used, %dGB available", totalGB, usedGB, availableGB)
			
			if totalGB < reqRAM {
				e.logger.Error("Insufficient RAM: required %dGB, available %dGB", reqRAM, totalGB)
				return fmt.Errorf("insufficient RAM: required %dGB, available %dGB", reqRAM, totalGB)
			}
			e.logger.Info("RAM requirement satisfied: %dGB/%dGB ✓", totalGB, reqRAM)
		}
	} else {
		e.logger.Debug("No RAM requirement specified")
	}
	
	// Check disk space
	if reqDisk := minimum.GetDiskGB(); reqDisk > 0 {
		e.logger.Info("Checking disk space requirement: %dGB needed", reqDisk)
		
		diskInfo, err := disk.Usage("/")
		if err != nil {
			e.logger.Error("Failed to get disk info: %v", err)
			e.logger.Warn("Skipping disk space validation due to error")
		} else {
			totalGB := int(diskInfo.Total / (1024 * 1024 * 1024))
			usedGB := int(diskInfo.Used / (1024 * 1024 * 1024))
			freeGB := int(diskInfo.Free / (1024 * 1024 * 1024))
			
			e.logger.Info("System disk: %dGB total, %dGB used, %dGB free", totalGB, usedGB, freeGB)
			
			if freeGB < reqDisk {
				e.logger.Error("Insufficient disk space: required %dGB, available %dGB", reqDisk, freeGB)
				return fmt.Errorf("insufficient disk space: required %dGB, available %dGB", reqDisk, freeGB)
			}
			e.logger.Info("Disk space requirement satisfied: %dGB/%dGB ✓", freeGB, reqDisk)
		}
	} else {
		e.logger.Debug("No disk space requirement specified")
	}
	
	e.logger.Info("All hardware requirements satisfied for arrow: %s", arrow.Name())
	return nil
}

// ValidateEnvironment checks if the environment is suitable for execution with verbose logging
func (e *Engine) ValidateEnvironment(arrow manifest.ArrowInterface, installPath string) error {
	e.logger.Info("Starting environment validation for arrow: %s", arrow.Name())
	e.logger.Info("Install path: %s", installPath)
	
	// Check if install path exists
	if _, err := os.Stat(installPath); os.IsNotExist(err) {
		e.logger.Error("Install path does not exist: %s", installPath)
		return fmt.Errorf("install path does not exist: %s", installPath)
	}
	e.logger.Info("Install path exists and is accessible")

	// Check platform support
	platform := runtime.GOOS
	arch := runtime.GOARCH
	
	e.logger.Info("System platform: %s/%s", platform, arch)
	
	methods := arrow.GetMethods()
	
	// Check if any method supports this platform and architecture
	supportedMethods := []types.MethodType{
		types.MethodInstall,
		types.MethodExecute,
		types.MethodUninstall,
	}
	
	e.logger.Debug("Checking platform support for methods: %v", supportedMethods)
	
	hasSupport := false
	supportedPlatforms := make(map[string][]string)
	
	for _, methodType := range supportedMethods {
		methodMap := methods.GetMethod(string(methodType))
		if methodMap != nil {
			e.logger.Debug("Checking method %s for platform support", methodType)
			
			for osName, osMap := range methodMap {
				if _, exists := supportedPlatforms[osName]; !exists {
					supportedPlatforms[osName] = make([]string, 0)
				}
				
				for archName := range osMap {
					// Add unique architectures
					found := false
					for _, existingArch := range supportedPlatforms[osName] {
						if existingArch == archName {
							found = true
							break
						}
					}
					if !found {
						supportedPlatforms[osName] = append(supportedPlatforms[osName], archName)
					}
				}
				
				if osName == platform {
					if _, exists := osMap[arch]; exists {
						e.logger.Debug("Method %s supports current platform %s/%s", methodType, platform, arch)
						hasSupport = true
					} else if len(osMap) > 0 {
						e.logger.Debug("Method %s supports current OS %s but not architecture %s", methodType, platform, arch)
						hasSupport = true // At least the OS is supported, fallback might work
					}
				}
			}
		}
	}
	
	// Log supported platforms for debugging
	e.logger.Debug("Arrow supports the following platforms:")
	for osName, arches := range supportedPlatforms {
		e.logger.Debug("  %s: %v", osName, arches)
	}
	
	if !hasSupport {
		e.logger.Error("Arrow does not support platform: %s/%s", platform, arch)
		return fmt.Errorf("arrow does not support platform: %s/%s", platform, arch)
	}
	
	e.logger.Info("Platform support validated successfully")
	
	// Validate hardware requirements
	e.logger.Debug("Proceeding to hardware requirements validation")
	return e.ValidateRequirements(arrow)
} 